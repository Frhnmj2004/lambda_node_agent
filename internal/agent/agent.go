package agent

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"lamda_node_agent/internal/blockchain"
	"lamda_node_agent/internal/docker"
	"lamda_node_agent/internal/nats"
	"lamda_node_agent/internal/storage"

	"github.com/ethereum/go-ethereum/crypto"
)

// JobMessage represents a job assignment message from NATS
type JobMessage struct {
	JobID        string `json:"job_id"`
	ImageName    string `json:"image_name"`
	InputFileCID string `json:"input_file_cid"`
	OutputPath   string `json:"output_path"`
}

// StatusUpdate represents a status update message to NATS
type StatusUpdate struct {
	AgentAddress string    `json:"agent_address"`
	JobID        string    `json:"job_id"`
	Status       string    `json:"status"`
	OutputCID    string    `json:"output_cid,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// Agent is the main orchestrator for the lamda_node_agent
type Agent struct {
	blockchainClient blockchain.BlockchainClient
	dockerManager    docker.Manager
	storageManager   storage.Manager
	natsClient       nats.Client
	privateKey       *ecdsa.PrivateKey
	address          string
	heartbeatTicker  *time.Ticker
}

// NewAgent creates a new agent instance
func NewAgent(
	blockchainClient blockchain.BlockchainClient,
	dockerManager docker.Manager,
	storageManager storage.Manager,
	natsClient nats.Client,
	privateKey *ecdsa.PrivateKey,
) *Agent {
	// Derive address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Agent{
		blockchainClient: blockchainClient,
		dockerManager:    dockerManager,
		storageManager:   storageManager,
		natsClient:       natsClient,
		privateKey:       privateKey,
		address:          address.Hex(),
	}
}

// Run starts the agent and orchestrates all operations
func (a *Agent) Run(ctx context.Context, gpuModel string, vram uint64) error {
	log.Printf("Starting lamda_node_agent with address: %s", a.address)
	log.Printf("GPU Model: %s, VRAM: %d MiB", gpuModel, vram)

	// Register node with the blockchain
	log.Printf("Registering node with blockchain...")
	if err := a.blockchainClient.RegisterNode(ctx, gpuModel, vram); err != nil {
		return fmt.Errorf("failed to register node: %w", err)
	}
	log.Printf("Node registered successfully")

	// Start heartbeat goroutine
	a.startHeartbeat(ctx)

	// Subscribe to job assignments
	subject := fmt.Sprintf("jobs.dispatch.%s", a.address)
	log.Printf("Subscribing to job assignments on subject: %s", subject)

	if err := a.natsClient.SubscribeToJobs(ctx, subject, a.handleJobMessage); err != nil {
		return fmt.Errorf("failed to subscribe to jobs: %w", err)
	}

	// Keep the agent running
	<-ctx.Done()
	log.Printf("Agent shutting down...")

	// Cleanup
	if a.heartbeatTicker != nil {
		a.heartbeatTicker.Stop()
	}
	a.natsClient.Close()

	return nil
}

// startHeartbeat starts a goroutine to send periodic heartbeats
func (a *Agent) startHeartbeat(ctx context.Context) {
	a.heartbeatTicker = time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <-a.heartbeatTicker.C:
				if err := a.blockchainClient.SendHeartbeat(ctx); err != nil {
					log.Printf("Failed to send heartbeat: %v", err)
				} else {
					log.Printf("Heartbeat sent successfully")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// handleJobMessage processes incoming job messages
func (a *Agent) handleJobMessage(msg []byte) {
	var jobMsg JobMessage
	if err := json.Unmarshal(msg, &jobMsg); err != nil {
		log.Printf("Failed to unmarshal job message: %v", err)
		return
	}

	log.Printf("Received job assignment: %s", jobMsg.JobID)

	// Create local directories for the job
	jobDir := filepath.Join(os.TempDir(), "lamda_jobs", jobMsg.JobID)
	inputDir := filepath.Join(jobDir, "input")
	outputDir := filepath.Join(jobDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		log.Printf("Failed to create input directory: %v", err)
		return
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("Failed to create output directory: %v", err)
		return
	}

	// Update status to "processing"
	a.publishStatus(jobMsg.JobID, "processing")

	// Download input data from IPFS
	if err := a.storageManager.DownloadInput(context.Background(), jobMsg.InputFileCID, inputDir); err != nil {
		log.Printf("Failed to download input data: %v", err)
		a.publishStatus(jobMsg.JobID, "failed")
		return
	}

	// Run the job container
	if err := a.dockerManager.RunJobContainer(context.Background(), jobMsg.ImageName, inputDir, outputDir); err != nil {
		log.Printf("Failed to run job container: %v", err)
		a.publishStatus(jobMsg.JobID, "failed")
		return
	}

	// Upload output data to IPFS
	outputCID, err := a.storageManager.UploadOutput(context.Background(), outputDir)
	if err != nil {
		log.Printf("Failed to upload output data: %v", err)
		a.publishStatus(jobMsg.JobID, "failed")
		return
	}

	// Update status to "completed" with output CID
	a.publishStatus(jobMsg.JobID, "completed", outputCID)

	// Cleanup
	os.RemoveAll(jobDir)

	log.Printf("Job %s completed successfully", jobMsg.JobID)
}

// publishStatus publishes a status update to NATS
func (a *Agent) publishStatus(jobID, status string, outputCID ...string) {
	statusUpdate := StatusUpdate{
		AgentAddress: a.address,
		JobID:        jobID,
		Status:       status,
		Timestamp:    time.Now(),
	}

	// Add output CID if provided
	if len(outputCID) > 0 && outputCID[0] != "" {
		statusUpdate.OutputCID = outputCID[0]
	}

	statusBytes, err := json.Marshal(statusUpdate)
	if err != nil {
		log.Printf("Failed to marshal status update: %v", err)
		return
	}

	if err := a.natsClient.PublishStatusUpdate(context.Background(), statusBytes); err != nil {
		log.Printf("Failed to publish status update: %v", err)
	}
}

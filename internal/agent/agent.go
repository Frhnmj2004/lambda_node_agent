package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"lamda_node_agent/internal/blockchain"
	"lamda_node_agent/internal/config"
	"lamda_node_agent/internal/docker"
	"lamda_node_agent/internal/nats"
	"lamda_node_agent/internal/storage"
)

// JobMessage represents the structure of a job assignment from NATS.
type JobMessage struct {
	JobID                  string `json:"jobId"`
	DockerImage            string `json:"dockerImage"`
	GreenfieldInputUrl     string `json:"greenfieldInputUrl"`
	GreenfieldOutputBucket string `json:"greenfieldOutputBucket"`
	GreenfieldOutputName   string `json:"greenfieldOutputName"`
}

// StatusUpdate represents the structure of a job status update.
type StatusUpdate struct {
	JobID     string `json:"jobId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	OutputUrl string `json:"outputUrl,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// Agent is the main struct that manages the node agent lifecycle and job execution.
type Agent struct {
	Blockchain blockchain.BlockchainClient
	Docker     docker.Manager
	Nats       nats.Client
	Storage    storage.Manager
	Config     *config.Config
	gpuModel   string
	vram       uint64
}

// NewAgent constructs a new Agent instance.
func NewAgent(bc blockchain.BlockchainClient, dm docker.Manager, nc nats.Client, sm storage.Manager, cfg *config.Config, gpuModel string, vram uint64) *Agent {
	return &Agent{
		Blockchain: bc,
		Docker:     dm,
		Nats:       nc,
		Storage:    sm,
		Config:     cfg,
		gpuModel:   gpuModel,
		vram:       vram,
	}
}

// Run starts the agent's main loop, including registration, heartbeat, and job processing.
func (a *Agent) Run(ctx context.Context) error {
	// Log agent startup
	slog.Info("Agent starting up", "gpuModel", a.gpuModel, "vram", a.vram)
	// Register node on blockchain
	_, err := a.Blockchain.RegisterNode(ctx, a.gpuModel, a.vram)
	if err != nil {
		return err
	}

	// Start heartbeat goroutine
	heartbeatCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		ticker := time.NewTicker(time.Duration(a.Config.HeartbeatIntervalSeconds) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_, _ = a.Blockchain.SendHeartbeat(heartbeatCtx)
			case <-heartbeatCtx.Done():
				return
			}
		}
	}()

	// Subscribe to jobs
	jobSubject := "jobs.dispatch." + "TODO_AGENT_WALLET_ADDRESS"
	err = a.Nats.SubscribeToJobs(ctx, jobSubject, a.handleJob)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

// handleJob is the handler function for incoming job messages.
func (a *Agent) handleJob(msg []byte) {
	var job JobMessage
	if err := json.Unmarshal(msg, &job); err != nil {
		slog.Error("Failed to unmarshal job message", "error", err)
		return
	}
	slog.Info("Received job", "jobId", job.JobID)

	// Create temp dirs
	inputDir, err := ioutil.TempDir("", "lamda_input_")
	if err != nil {
		a.publishStatus(job.JobID, "failed", "Failed to create input temp dir: "+err.Error(), "")
		return
	}
	defer os.RemoveAll(inputDir)
	outputDir, err := ioutil.TempDir("", "lamda_output_")
	if err != nil {
		a.publishStatus(job.JobID, "failed", "Failed to create output temp dir: "+err.Error(), "")
		return
	}
	defer os.RemoveAll(outputDir)

	inputPath := filepath.Join(inputDir, "input.zip")
	outputPath := filepath.Join(outputDir, "output.zip")

	// Download input
	if err := a.Storage.DownloadInput(context.Background(), job.GreenfieldInputUrl, inputPath); err != nil {
		a.publishStatus(job.JobID, "failed", "Failed to download input: "+err.Error(), "")
		return
	}

	// Run Docker job
	if err := a.Docker.RunJobContainer(context.Background(), job.DockerImage, inputDir, outputDir); err != nil {
		a.publishStatus(job.JobID, "failed", "Docker job failed: "+err.Error(), "")
		return
	}

	// Upload output
	if err := a.Storage.UploadOutput(context.Background(), outputPath, job.GreenfieldOutputBucket, job.GreenfieldOutputName); err != nil {
		a.publishStatus(job.JobID, "failed", "Failed to upload output: "+err.Error(), "")
		return
	}

	outputUrl := fmt.Sprintf("gnfd://%s/%s", job.GreenfieldOutputBucket, job.GreenfieldOutputName)
	a.publishStatus(job.JobID, "success", "Job completed successfully", outputUrl)
}

func (a *Agent) publishStatus(jobID, status, message, outputUrl string) {
	update := StatusUpdate{
		JobID:     jobID,
		Status:    status,
		Message:   message,
		OutputUrl: outputUrl,
		Timestamp: time.Now().Unix(),
	}
	b, _ := json.Marshal(update)
	a.Nats.PublishStatusUpdate(context.Background(), b)
}

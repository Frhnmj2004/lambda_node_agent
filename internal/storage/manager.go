package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"lamda_node_agent/internal/config"

	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// Manager defines the interface for storage operations
type Manager interface {
	DownloadInput(ctx context.Context, jobID string, localPath string) error
	UploadOutput(ctx context.Context, jobID string, localPath string) error
}

// GreenfieldManager implements Manager using BNB Greenfield
type GreenfieldManager struct {
	client       client.IClient
	agentAddress common.Address
	bucketName   string
}

// NewGreenfieldManager creates a new Greenfield storage manager
func NewGreenfieldManager(cfg *config.Config, agentAddress common.Address) (*GreenfieldManager, error) {
	// Create account from private key
	account, err := types.NewAccountFromPrivateKey("agent", cfg.AgentPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create account from private key: %w", err)
	}

	// Initialize Greenfield client with the correct chain ID for testnet
	greenfieldClient, err := client.New(
		cfg.GreenfieldEndpoint,
		"greenfield_5600-1", // Testnet chain ID
		client.Option{
			DefaultAccount: account,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Greenfield client: %w", err)
	}

	return &GreenfieldManager{
		client:       greenfieldClient,
		agentAddress: agentAddress,
		bucketName:   cfg.GreenfieldBucketName,
	}, nil
}

// DownloadInput downloads input data for a job from Greenfield
func (g *GreenfieldManager) DownloadInput(ctx context.Context, jobID string, localPath string) error {
	// Create the local directory if it doesn't exist
	if err := os.MkdirAll(localPath, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// Define the object key for input data
	objectKey := fmt.Sprintf("jobs/%s/input", jobID)

	// Download the object from Greenfield
	resp, _, err := g.client.GetObject(ctx, g.bucketName, objectKey, types.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get object from Greenfield: %w", err)
	}
	defer resp.Close()

	// Create the local file
	localFile := filepath.Join(localPath, "input")
	file, err := os.Create(localFile)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	// Copy the data from Greenfield to local file
	_, err = io.Copy(file, resp)
	if err != nil {
		return fmt.Errorf("failed to copy data to local file: %w", err)
	}

	return nil
}

// UploadOutput uploads output data for a job to Greenfield
func (g *GreenfieldManager) UploadOutput(ctx context.Context, jobID string, localPath string) error {
	// Define the object key for output data
	objectKey := fmt.Sprintf("jobs/%s/output", jobID)

	// Open the local output file
	outputFile := filepath.Join(localPath, "output")
	file, err := os.Open(outputFile)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer file.Close()

	// Get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Upload the file to Greenfield
	err = g.client.PutObject(ctx, g.bucketName, objectKey, fileInfo.Size(), file, types.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload object to Greenfield: %w", err)
	}

	return nil
}

package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"lamda_node_agent/internal/config"
)

// Manager defines the interface for storage operations
type Manager interface {
	DownloadInput(ctx context.Context, ipfsCID string, localPath string) error
	UploadOutput(ctx context.Context, localPath string) (string, error)
}

// IPFSManager implements Manager using IPFS with Pinata
type IPFSManager struct {
	client    *http.Client
	pinataJWT string
}

// PinataResponse represents the response from Pinata API
type PinataResponse struct {
	IpfsHash  string `json:"IpfsHash"`
	PinSize   int    `json:"PinSize"`
	Timestamp string `json:"Timestamp"`
}

// NewIPFSManager creates a new IPFS storage manager
func NewIPFSManager(cfg *config.Config) (*IPFSManager, error) {
	return &IPFSManager{
		client:    &http.Client{},
		pinataJWT: cfg.PinataJWT,
	}, nil
}

// DownloadInput downloads input data for a job from IPFS
func (i *IPFSManager) DownloadInput(ctx context.Context, ipfsCID string, localPath string) error {
	// Create the local directory if it doesn't exist
	if err := os.MkdirAll(localPath, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// Construct the IPFS gateway URL
	gatewayURL := fmt.Sprintf("https://gateway.pinata.cloud/ipfs/%s", ipfsCID)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", gatewayURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Make the request
	resp, err := i.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download from IPFS: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download from IPFS: HTTP %d", resp.StatusCode)
	}

	// Create the local file
	localFile := filepath.Join(localPath, "input")
	file, err := os.Create(localFile)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	// Copy the data from IPFS to local file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy data to local file: %w", err)
	}

	return nil
}

// UploadOutput uploads output data for a job to IPFS via Pinata
func (i *IPFSManager) UploadOutput(ctx context.Context, localPath string) (string, error) {
	// Open the local output file
	outputFile := filepath.Join(localPath, "output")
	file, err := os.Open(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to open output file: %w", err)
	}
	defer file.Close()

	// Create a buffer to store the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field
	part, err := writer.CreateFormFile("file", "output")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy the file content to the form field
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file to form: %w", err)
	}

	// Close the multipart writer
	writer.Close()

	// Create HTTP request to Pinata API
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.pinata.cloud/pinning/pinFileToIPFS", &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+i.pinataJWT)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request
	resp, err := i.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload to Pinata: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload to Pinata: HTTP %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var pinataResp PinataResponse
	if err := json.NewDecoder(resp.Body).Decode(&pinataResp); err != nil {
		return "", fmt.Errorf("failed to parse Pinata response: %w", err)
	}

	return pinataResp.IpfsHash, nil
}

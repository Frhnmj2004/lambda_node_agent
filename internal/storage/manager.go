package storage

import (
	"context"
	"errors"

	greenfield "github.com/bnb-chain/greenfield-go-sdk/client"
)

// Manager defines the interface for decentralized storage operations.
type Manager interface {
	// DownloadInput downloads the input data from Greenfield to the local path.
	DownloadInput(ctx context.Context, greenfieldUrl, localPath string) error
	// UploadOutput uploads the output data from the local path to Greenfield.
	UploadOutput(ctx context.Context, localPath, bucket, objectName string) error
}

// GreenfieldManager implements Manager using the greenfield-go-sdk.
type GreenfieldManager struct {
	client *greenfield.Client
}

// NewGreenfieldManager creates a new GreenfieldManager instance.
func NewGreenfieldManager(endpoint, accessKey, secretKey string) (*GreenfieldManager, error) {
	cfg := &greenfield.Config{
		Endpoint:  endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
	cli, err := greenfield.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &GreenfieldManager{client: cli}, nil
}

// DownloadInput downloads the input data from Greenfield to the local path.
func (g *GreenfieldManager) DownloadInput(ctx context.Context, greenfieldUrl, localPath string) error {
	// TODO: Parse greenfieldUrl and use SDK to download file
	return errors.New("DownloadInput not implemented: SDK call required")
}

// UploadOutput uploads the output data from the local path to Greenfield.
func (g *GreenfieldManager) UploadOutput(ctx context.Context, localPath, bucket, objectName string) error {
	// TODO: Use SDK to upload file
	return errors.New("UploadOutput not implemented: SDK call required")
}

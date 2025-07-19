package docker

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// Manager defines the interface for Docker operations
type Manager interface {
	RunJobContainer(ctx context.Context, imageName string, inputPath string, outputPath string) error
}

// dockerManager implements Manager using Docker
type dockerManager struct {
	client *client.Client
}

// NewDockerManager creates a new Docker manager
func NewDockerManager() (Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &dockerManager{
		client: cli,
	}, nil
}

// RunJobContainer runs a Docker container for job execution with GPU access
func (d *dockerManager) RunJobContainer(ctx context.Context, imageName string, inputPath string, outputPath string) error {
	// Pull the Docker image
	log.Printf("Pulling Docker image: %s", imageName)
	reader, err := d.client.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", imageName, err)
	}
	defer reader.Close()
	io.Copy(io.Discard, reader)

	// Create container configuration
	config := &container.Config{
		Image: imageName,
		Cmd:   []string{"/bin/bash", "-c", "ls /input && ls /output"},
	}

	// Create host configuration with volume mounts
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: inputPath,
				Target: "/input",
			},
			{
				Type:   mount.TypeBind,
				Source: outputPath,
				Target: "/output",
			},
		},
		// Note: GPU access requires nvidia-docker runtime
		// This can be configured via Docker daemon settings
	}

	// Create the container
	log.Printf("Creating container for image: %s", imageName)
	resp, err := d.client.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	containerID := resp.ID
	log.Printf("Created container with ID: %s", containerID)

	// Start the container
	log.Printf("Starting container: %s", containerID)
	if err := d.client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	// Stream container logs
	logs, err := d.client.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		log.Printf("Warning: failed to get container logs: %v", err)
	} else {
		defer logs.Close()
		go func() {
			io.Copy(log.Writer(), logs)
		}()
	}

	// Wait for container to complete
	log.Printf("Waiting for container to complete: %s", containerID)
	statusCh, errCh := d.client.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error waiting for container: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("container exited with status code: %d", status.StatusCode)
		}
	}

	// Remove the container
	log.Printf("Removing container: %s", containerID)
	if err := d.client.ContainerRemove(ctx, containerID, container.RemoveOptions{}); err != nil {
		log.Printf("Warning: failed to remove container: %v", err)
	}

	log.Printf("Container execution completed successfully")
	return nil
}

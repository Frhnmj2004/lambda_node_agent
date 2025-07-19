package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// Manager defines the interface for running job containers with GPU access.
type Manager interface {
	// RunJobContainer pulls the image, creates, starts, and waits for a container with GPU access.
	// inputPath and outputPath are mounted as volumes.
	RunJobContainer(ctx context.Context, imageName string, inputPath string, outputPath string) error
}

// DockerManager implements Manager using the Docker Go SDK.
type DockerManager struct {
	client *dockerclient.Client
}

// NewDockerManager creates a new DockerManager instance.
func NewDockerManager() (*DockerManager, error) {
	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerManager{client: cli}, nil
}

// RunJobContainer runs a container with GPU access and mounts input/output directories.
func (d *DockerManager) RunJobContainer(ctx context.Context, imageName string, inputPath string, outputPath string) error {
	// Pull the image
	out, err := d.client.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer out.Close()
	io.Copy(io.Discard, out)

	// Prepare mounts
	mounts := []mount.Mount{
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
	}

	// Create container with GPU access
	resp, err := d.client.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
			Tty:   false,
		},
		&container.HostConfig{
			Mounts: mounts,
			DeviceRequests: []container.DeviceRequest{
				{
					Driver:       "nvidia",
					Count:        -1,
					Capabilities: [][]string{{"gpu"}},
				},
			},
		},
		nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := d.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	// Stream logs
	logReader, err := d.client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err == nil {
		defer logReader.Close()
		stdcopy.StdCopy(os.Stdout, os.Stderr, logReader)
	}

	// Wait for completion
	statusCh, errCh := d.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("container wait error: %w", err)
		}
	case <-statusCh:
	}

	// Remove container
	if err := d.client.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}

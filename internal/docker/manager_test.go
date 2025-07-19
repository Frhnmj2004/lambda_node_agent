package docker

import (
	"context"
	"testing"
)

type mockDockerClient struct {
	containerCreated bool
	gpuRequested     bool
	inputMounted     bool
	outputMounted    bool
}

func (m *mockDockerClient) ContainerCreateWithGPU(imageName, inputPath, outputPath string) {
	m.containerCreated = true
	m.gpuRequested = true // Simulate GPU config
	m.inputMounted = inputPath != ""
	m.outputMounted = outputPath != ""
}

func TestDockerManager_RunJobContainer(t *testing.T) {
	mock := &mockDockerClient{}
	dm := &DockerManager{
		// TODO: inject mock client
	}

	err := dm.RunJobContainer(context.Background(), "pytorch/pytorch:latest-cuda", "/tmp/input", "/tmp/output")
	// TODO: call mock.ContainerCreateWithGPU inside RunJobContainer for test

	if err != nil {
		t.Errorf("RunJobContainer returned error: %v", err)
	}
	if !mock.containerCreated {
		t.Error("Container was not created")
	}
	if !mock.gpuRequested {
		t.Error("GPU was not requested in HostConfig")
	}
	if !mock.inputMounted {
		t.Error("Input path was not mounted")
	}
	if !mock.outputMounted {
		t.Error("Output path was not mounted")
	}
}

package hwinfo

import (
	"os/exec"
	"testing"
)

func TestGetNvidiaGPUInfo(t *testing.T) {
	// Test successful detection
	gpuModel, vram, err := GetNvidiaGPUInfo()
	if err != nil {
		// This is expected if nvidia-smi is not available on the test system
		t.Logf("nvidia-smi not available or failed: %v", err)
		return
	}

	if gpuModel == "" {
		t.Error("GPU model should not be empty")
	}

	if vram == 0 {
		t.Error("VRAM should be greater than 0")
	}

	t.Logf("Detected GPU: %s with %d MiB VRAM", gpuModel, vram)
}

func TestGetNvidiaGPUInfo_CommandNotFound(t *testing.T) {
	// Test with a non-existent command
	originalCommand := "nvidia-smi"
	defer func() {
		// Restore original behavior after test
		exec.Command = exec.Command
	}()

	// Mock exec.Command to return a non-existent command
	exec.Command = func(name string, args ...string) *exec.Cmd {
		return exec.Command("non-existent-command")
	}

	_, _, err := GetNvidiaGPUInfo()
	if err == nil {
		t.Error("Expected error when nvidia-smi command is not found")
	}
}

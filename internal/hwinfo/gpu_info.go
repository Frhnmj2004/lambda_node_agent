package hwinfo

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GetNvidiaGPUInfo executes nvidia-smi to get GPU model and VRAM
func GetNvidiaGPUInfo() (gpuModel string, vramMiB uint64, err error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=gpu_name,memory.total", "--format=csv,noheader")
	output, err := cmd.Output()
	if err != nil {
		return "", 0, fmt.Errorf("failed to execute nvidia-smi: %w", err)
	}

	// Parse the output
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return "", 0, fmt.Errorf("no GPU information found")
	}

	// Take the first GPU (assuming single GPU setup)
	line := strings.TrimSpace(lines[0])
	parts := strings.Split(line, ",")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("unexpected nvidia-smi output format: %s", line)
	}

	gpuModel = strings.TrimSpace(parts[0])
	memoryStr := strings.TrimSpace(parts[1])

	// Parse memory (format: "8192 MiB")
	memoryParts := strings.Fields(memoryStr)
	if len(memoryParts) != 2 || memoryParts[1] != "MiB" {
		return "", 0, fmt.Errorf("unexpected memory format: %s", memoryStr)
	}

	vramMiB, err = strconv.ParseUint(memoryParts[0], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse VRAM value: %w", err)
	}

	return gpuModel, vramMiB, nil
}

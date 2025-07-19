package hwinfo

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GetNvidiaGPUInfo executes the nvidia-smi command to fetch the GPU model
// and total VRAM. It returns the GPU name, VRAM in MiB, and an error
// if the command fails or parsing is unsuccessful.
func GetNvidiaGPUInfo() (gpuModel string, vramMiB uint64, err error) {
	// Execute the nvidia-smi command
	cmd := exec.Command("nvidia-smi", "--query-gpu=gpu_name,memory.total", "--format=csv,noheader")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", 0, fmt.Errorf("nvidia-smi command failed: %v - output: %s", err, string(output))
	}

	// Parse the output
	rawOutput := strings.TrimSpace(string(output))
	parts := strings.Split(rawOutput, ",")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("unexpected nvidia-smi output format: %s", rawOutput)
	}

	gpuModel = strings.TrimSpace(parts[0])
	vramStr := strings.TrimSpace(parts[1])
	vramStr = strings.TrimSuffix(vramStr, " MiB")

	vramMiB, err = strconv.ParseUint(vramStr, 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse VRAM value '%s': %v", vramStr, err)
	}

	return gpuModel, vramMiB, nil
}

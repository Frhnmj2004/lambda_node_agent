package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"lamda_node_agent/internal/agent"
	"lamda_node_agent/internal/blockchain"
	"lamda_node_agent/internal/config"
	"lamda_node_agent/internal/docker"
	"lamda_node_agent/internal/hwinfo"
	"lamda_node_agent/internal/nats"
	"lamda_node_agent/internal/storage"
)

// main is the entry point for the lamda_node_agent application.
func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Detect GPU information
	gpuModel, vram, err := hwinfo.GetNvidiaGPUInfo()
	if err != nil {
		log.Fatalf("FATAL: Could not detect GPU information. Ensure NVIDIA drivers and nvidia-smi are installed and accessible. Error: %v", err)
	}
	log.Printf("Detected GPU: %s with %d MiB VRAM", gpuModel, vram)

	bc, err := blockchain.NewEthClient(cfg.OpBnbRpcUrl, cfg.AgentPrivateKey, cfg.ReputationContractAddress)
	if err != nil {
		slog.Error("Failed to initialize blockchain client", "error", err)
		os.Exit(1)
	}
	dm, err := docker.NewDockerManager()
	if err != nil {
		slog.Error("Failed to initialize Docker manager", "error", err)
		os.Exit(1)
	}
	nc, err := nats.NewNatsClient(cfg.NatsUrl)
	if err != nil {
		slog.Error("Failed to initialize NATS client", "error", err)
		os.Exit(1)
	}
	sm, err := storage.NewGreenfieldManager(cfg.GreenfieldEndpoint, cfg.GreenfieldAccessKey, cfg.GreenfieldSecretKey)
	if err != nil {
		slog.Error("Failed to initialize Greenfield manager", "error", err)
		os.Exit(1)
	}

	a := agent.NewAgent(bc, dm, nc, sm, cfg, gpuModel, vram)
	if err := a.Run(ctx); err != nil {
		slog.Error("Agent exited with error", "error", err)
		os.Exit(1)
	}
}

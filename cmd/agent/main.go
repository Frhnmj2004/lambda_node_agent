package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"lamda_node_agent/internal/agent"
	"lamda_node_agent/internal/blockchain"
	"lamda_node_agent/internal/config"
	"lamda_node_agent/internal/docker"
	"lamda_node_agent/internal/hwinfo"
	"lamda_node_agent/internal/nats"
	"lamda_node_agent/internal/storage"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Get GPU information
	gpuModel, vramMiB, err := hwinfo.GetNvidiaGPUInfo()
	if err != nil {
		log.Fatalf("Failed to get GPU information: %v", err)
	}
	log.Printf("Detected GPU: %s with %d MiB VRAM", gpuModel, vramMiB)

	// Parse private key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.AgentPrivateKey, "0x"))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Initialize blockchain client
	blockchainClient, err := blockchain.NewEthClient(
		cfg.OpBNBRPCURL,
		cfg.AgentPrivateKey,
		cfg.ReputationContractAddress,
	)
	if err != nil {
		log.Fatalf("Failed to create blockchain client: %v", err)
	}

	// Initialize Docker manager
	dockerManager, err := docker.NewDockerManager()
	if err != nil {
		log.Fatalf("Failed to create Docker manager: %v", err)
	}

	// Initialize storage manager
	storageManager, err := storage.NewIPFSManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create storage manager: %v", err)
	}

	// Initialize NATS client
	natsClient, err := nats.NewNatsClient(cfg.NatsURL)
	if err != nil {
		log.Fatalf("Failed to create NATS client: %v", err)
	}

	// Create agent
	agent := agent.NewAgent(
		blockchainClient,
		dockerManager,
		storageManager,
		natsClient,
		privateKey,
	)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v, shutting down...", sig)
		cancel()
	}()

	// Run the agent
	log.Printf("Starting lamda_node_agent...")
	if err := agent.Run(ctx, gpuModel, vramMiB); err != nil {
		log.Fatalf("Agent failed: %v", err)
	}

	log.Printf("lamda_node_agent stopped")
}

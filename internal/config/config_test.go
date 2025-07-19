package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary .env file
	dotenvContent := `
OPBNB_RPC_URL=https://opbnb-testnet-rpc.com
AGENT_PRIVATE_KEY=0xabc123
NATS_URL=nats://localhost:4222
REPUTATION_CONTRACT_ADDRESS=0xdef456
GREENFIELD_ENDPOINT=https://greenfield-endpoint.com
GREENFIELD_ACCESS_KEY=accesskey
GREENFIELD_SECRET_KEY=secretkey
GREENFIELD_BUCKET=test-bucket
HEARTBEAT_INTERVAL_SECONDS=120
LOG_LEVEL=debug
`
	tmpfile, err := ioutil.TempFile("", ".env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write([]byte(dotenvContent)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	// Set environment variable to point to the temp .env file
	os.Setenv("ENV", tmpfile.Name())
	os.Setenv("OPBNB_RPC_URL", "https://opbnb-testnet-rpc.com")
	os.Setenv("AGENT_PRIVATE_KEY", "0xabc123")
	os.Setenv("NATS_URL", "nats://localhost:4222")
	os.Setenv("REPUTATION_CONTRACT_ADDRESS", "0xdef456")
	os.Setenv("GREENFIELD_ENDPOINT", "https://greenfield-endpoint.com")
	os.Setenv("GREENFIELD_ACCESS_KEY", "accesskey")
	os.Setenv("GREENFIELD_SECRET_KEY", "secretkey")
	os.Setenv("GREENFIELD_BUCKET", "test-bucket")
	os.Setenv("HEARTBEAT_INTERVAL_SECONDS", "120")
	os.Setenv("LOG_LEVEL", "debug")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.OpBnbRpcUrl != "https://opbnb-testnet-rpc.com" {
		t.Errorf("OpBnbRpcUrl mismatch: got %s", cfg.OpBnbRpcUrl)
	}
	if cfg.AgentPrivateKey != "0xabc123" {
		t.Errorf("AgentPrivateKey mismatch: got %s", cfg.AgentPrivateKey)
	}
	if cfg.NatsUrl != "nats://localhost:4222" {
		t.Errorf("NatsUrl mismatch: got %s", cfg.NatsUrl)
	}
	if cfg.ReputationContractAddress != "0xdef456" {
		t.Errorf("ReputationContractAddress mismatch: got %s", cfg.ReputationContractAddress)
	}
	if cfg.GreenfieldEndpoint != "https://greenfield-endpoint.com" {
		t.Errorf("GreenfieldEndpoint mismatch: got %s", cfg.GreenfieldEndpoint)
	}
	if cfg.GreenfieldAccessKey != "accesskey" {
		t.Errorf("GreenfieldAccessKey mismatch: got %s", cfg.GreenfieldAccessKey)
	}
	if cfg.GreenfieldSecretKey != "secretkey" {
		t.Errorf("GreenfieldSecretKey mismatch: got %s", cfg.GreenfieldSecretKey)
	}
	if cfg.GreenfieldBucket != "test-bucket" {
		t.Errorf("GreenfieldBucket mismatch: got %s", cfg.GreenfieldBucket)
	}
	if cfg.HeartbeatIntervalSeconds != 120 {
		t.Errorf("HeartbeatIntervalSeconds mismatch: got %d", cfg.HeartbeatIntervalSeconds)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel mismatch: got %s", cfg.LogLevel)
	}
}

package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the lamda_node_agent
type Config struct {
	// Blockchain Configuration
	OpBNBRPCURL               string `env:"OPBNB_RPC_URL" envDefault:"https://opbnb-testnet-rpc.bnbchain.org"`
	AgentPrivateKey           string `env:"AGENT_PRIVATE_KEY,required"`
	ReputationContractAddress string `env:"REPUTATION_CONTRACT_ADDRESS" envDefault:"0x108f2c400C9828d8044a5F6985f0C9589B90758D"`

	// NATS Configuration
	NatsURL string `env:"NATS_URL" envDefault:"nats://localhost:4222"`

	// Docker Configuration
	DockerHost string `env:"DOCKER_HOST" envDefault:"unix:///var/run/docker.sock"`

	// Agent Configuration
	HeartbeatInterval string `env:"HEARTBEAT_INTERVAL" envDefault:"5m"`
	LogLevel          string `env:"LOG_LEVEL" envDefault:"info"`

	// Greenfield Configuration
	GreenfieldEndpoint   string `env:"GREENFIELD_ENDPOINT" envDefault:"https://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443"`
	GreenfieldBucketName string `env:"GREENFIELD_BUCKET_NAME" envDefault:"lamda-jobs"`
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

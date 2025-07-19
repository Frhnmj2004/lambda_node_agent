package config

import (
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config holds all configuration variables for the agent.
type Config struct {
	OpBnbRpcUrl               string `env:"OPBNB_RPC_URL,required"`
	AgentPrivateKey           string `env:"AGENT_PRIVATE_KEY,required"`
	NatsUrl                   string `env:"NATS_URL,required"`
	ReputationContractAddress string `env:"REPUTATION_CONTRACT_ADDRESS,required"`
	GreenfieldEndpoint        string `env:"GREENFIELD_ENDPOINT,required"`
	GreenfieldAccessKey       string `env:"GREENFIELD_ACCESS_KEY,required"`
	GreenfieldSecretKey       string `env:"GREENFIELD_SECRET_KEY,required"`
	GreenfieldBucket          string `env:"GREENFIELD_BUCKET,required"`
	HeartbeatIntervalSeconds  int    `env:"HEARTBEAT_INTERVAL_SECONDS" envDefault:"300"`
	LogLevel                  string `env:"LOG_LEVEL" envDefault:"info"`
}

var (
	config     *Config
	configOnce sync.Once
)

// LoadConfig loads configuration from a .env file and environment variables.
// It returns a pointer to the Config struct and an error if loading fails.
func LoadConfig() (*Config, error) {
	var err error
	configOnce.Do(func() {
		_ = godotenv.Load()
		cfg := &Config{}
		err = env.Parse(cfg)
		if err == nil {
			config = cfg
		}
	})
	return config, err
}

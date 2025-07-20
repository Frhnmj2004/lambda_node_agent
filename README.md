# Lamda Node Agent

The lamda_node_agent is the "workforce" of the Lamda Network, responsible for executing compute jobs with GPU access and maintaining node reputation on the blockchain.

## Architecture

The agent follows a clean architecture pattern with the following components:

```
/lamda_node_agent
|-- cmd/
|   |-- agent/
|       |-- main.go              # Application entry point
|-- internal/
|   |-- agent/
|   |   |-- agent.go             # Core orchestrator
|   |-- blockchain/
|   |   |-- client.go            # Ethereum client implementation
|   |   |-- nodereputation.go    # Smart contract bindings
|   |-- config/
|   |   |-- config.go            # Configuration management
|   |-- docker/
|   |   |-- manager.go           # Docker container management
|   |-- hwinfo/
|   |   |-- gpu_info.go          # GPU hardware detection
|   |-- nats/
|   |   |-- client.go            # NATS messaging client
|   |-- storage/
|       |-- manager.go           # Storage operations (Greenfield placeholder)
|-- go.mod                       # Go module dependencies
|-- .env                         # Environment configuration
```

## Features

- **Blockchain Integration**: Registers nodes and sends heartbeats to the NodeReputation smart contract on opBNB Testnet
- **GPU Detection**: Automatically detects NVIDIA GPU model and VRAM capacity
- **Docker Job Execution**: Runs compute jobs in Docker containers with GPU access
- **NATS Messaging**: Receives job assignments and publishes status updates
- **Graceful Shutdown**: Handles SIGINT/SIGTERM signals for clean shutdown

## Prerequisites

- Go 1.21 or later
- NVIDIA GPU with nvidia-smi available
- Docker with GPU support (nvidia-docker)
- NATS server running
- opBNB Testnet RPC access
- Ethereum private key for node registration

## Configuration

Create a `.env` file with the following variables:

```env
# Blockchain Configuration
OPBNB_RPC_URL=https://opbnb-testnet-rpc.bnbchain.org
AGENT_PRIVATE_KEY=your_private_key_here
REPUTATION_CONTRACT_ADDRESS=0x108f2c400C9828d8044a5F6985f0C9589B90758D

# NATS Configuration
NATS_URL=nats://localhost:4222

# Docker Configuration
DOCKER_HOST=unix:///var/run/docker.sock

# IPFS Configuration
PINATA_JWT=your_pinata_jwt_token_here

# Agent Configuration
HEARTBEAT_INTERVAL=5m
LOG_LEVEL=info
```

## Building

1. Download dependencies:
```bash
go mod tidy
```

2. Build the agent:
```bash
go build -o lamda_node_agent ./cmd/agent
```

3. Run the agent:
```bash
./lamda_node_agent
```

## How It Works

1. **Startup**: The agent loads configuration, detects GPU hardware, and initializes all clients
2. **Registration**: Registers the node with the NodeReputation smart contract using GPU specifications
3. **Heartbeat**: Sends periodic heartbeats every 5 minutes to maintain node status
4. **Job Processing**: Subscribes to `jobs.dispatch.<agent_address>` for job assignments
5. **Job Execution**: Downloads input data, runs Docker container with GPU access, uploads results
6. **Status Updates**: Publishes job status updates to NATS for monitoring

## Job Message Format

Jobs are received via NATS with the following JSON format:

```json
{
  "job_id": "unique-job-identifier",
  "image_name": "docker-image:tag",
  "input_file_cid": "QmX...",
  "output_path": "/path/to/output/data"
}
```

## Status Update Format

Status updates are published to NATS with the following JSON format:

```json
{
  "agent_address": "0x...",
  "job_id": "unique-job-identifier",
  "status": "processing|completed|failed",
  "output_cid": "QmX...",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Docker Integration

The agent runs Docker containers with:
- GPU access via nvidia-docker
- Input/output volume mounts
- Automatic cleanup after completion
- Log streaming to stdout/stderr

## Smart Contract Integration

The agent interacts with the NodeReputation contract at `0x108f2c400C9828d8044a5F6985f0C9589B90758D` on opBNB Testnet:

- `registerNode(gpuModel, vram)`: Registers node with hardware specifications
- `sendHeartbeat()`: Sends periodic heartbeat to maintain active status

## IPFS Storage Integration

The agent uses IPFS for decentralized storage with Pinata as the pinning service:
- Input data download from IPFS using Pinata gateway
- Output data upload to IPFS via Pinata API
- Automatic content addressing with IPFS CIDs

## Error Handling

The agent includes comprehensive error handling for:
- Network connectivity issues
- Docker container failures
- Blockchain transaction failures
- Hardware detection failures
- Configuration errors

## Logging

The agent provides detailed logging for:
- Startup and initialization
- Job processing steps
- Blockchain transactions
- Docker operations
- Error conditions

## Development

To run in development mode:

1. Set up a local NATS server
2. Configure Docker with GPU support
3. Set up opBNB Testnet RPC access
4. Use a test private key
5. Run with debug logging

## Production Deployment

For production deployment:

1. Use a production NATS cluster
2. Ensure Docker GPU support is properly configured
3. Use a secure private key management system
4. Set up monitoring and alerting
5. Configure proper logging and metrics collection

## Troubleshooting

Common issues and solutions:

- **GPU Detection Failed**: Ensure nvidia-smi is available and working
- **Docker Connection Failed**: Check Docker daemon is running and accessible
- **NATS Connection Failed**: Verify NATS server is running and accessible
- **Blockchain Registration Failed**: Check RPC URL and private key configuration
- **Container GPU Access Failed**: Ensure nvidia-docker is properly configured

## Contributing

1. Follow the existing code structure and patterns
2. Add comprehensive error handling
3. Include logging for debugging
4. Update documentation for new features
5. Test with the opBNB Testnet environment

## License

This project is part of the Lamda Network ecosystem. 
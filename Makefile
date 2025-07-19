# Lamda Node Agent Makefile

.PHONY: build clean test deps run help

# Build variables
BINARY_NAME=lamda_node_agent
BUILD_DIR=build
MAIN_PATH=./cmd/agent

# Default target
all: deps build

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod tidy
	go mod download

# Build the agent
build: deps
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for production (optimized)
build-prod: deps
	@echo "Building $(BINARY_NAME) for production..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Production build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the agent
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run in development mode
dev: build
	@echo "Running $(BINARY_NAME) in development mode..."
	LOG_LEVEL=debug ./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Generate contract bindings (placeholder)
generate-bindings:
	@echo "Generating smart contract bindings..."
	@echo "Note: This requires abigen and the contract ABI"
	# abigen --abi=NodeReputation.abi --pkg=blockchain --out=internal/blockchain/nodereputation.go

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t lamda-node-agent .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run --rm --gpus all -v /var/run/docker.sock:/var/run/docker.sock lamda-node-agent

# Show help
help:
	@echo "Available targets:"
	@echo "  deps              - Download dependencies"
	@echo "  build             - Build the agent"
	@echo "  build-prod        - Build optimized for production"
	@echo "  run               - Build and run the agent"
	@echo "  dev               - Run in development mode"
	@echo "  clean             - Clean build artifacts"
	@echo "  test              - Run tests"
	@echo "  test-coverage     - Run tests with coverage"
	@echo "  fmt               - Format code"
	@echo "  lint              - Lint code"
	@echo "  install-tools     - Install development tools"
	@echo "  generate-bindings - Generate smart contract bindings"
	@echo "  docker-build      - Build Docker image"
	@echo "  docker-run        - Run Docker container"
	@echo "  help              - Show this help" 
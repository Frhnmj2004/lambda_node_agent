# Multi-stage build for lamda_node_agent
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o lamda_node_agent ./cmd/agent

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S lamda && \
    adduser -u 1001 -S lamda -G lamda

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/lamda_node_agent .

# Change ownership to non-root user
RUN chown -R lamda:lamda /app

# Switch to non-root user
USER lamda

# Expose port (if needed for health checks)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ps aux | grep lamda_node_agent || exit 1

# Run the agent
ENTRYPOINT ["./lamda_node_agent"] 
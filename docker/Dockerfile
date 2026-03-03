# Multi-stage build for NeuronAgent
# Build arguments
ARG GO_VERSION=1.23
ARG VERSION=latest
ARG BUILD_DATE
ARG VCS_REF

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /build

# Copy go mod files first for better layer caching
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy source code
COPY src/ .

# Build the application with version information
ARG VERSION
ARG BUILD_DATE
ARG VCS_REF
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.vcsRef=${VCS_REF}" \
    -o neuronagent ./cmd/agent-server

# Final stage - use Debian slim for better compatibility
FROM debian:bookworm-slim

# Build arguments for labels
ARG VERSION=latest
ARG BUILD_DATE
ARG VCS_REF

# Add metadata labels
LABEL org.opencontainers.image.title="NeuronAgent" \
      org.opencontainers.image.description="AI agent runtime system providing REST API and WebSocket endpoints" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${VCS_REF}" \
      org.opencontainers.image.vendor="neurondb" \
      org.opencontainers.image.licenses="Proprietary" \
      org.opencontainers.image.source="https://github.com/neurondb/NeurondB" \
      maintainer="neurondb <admin@neurondb.com>"

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    netcat-openbsd && \
    rm -rf /var/lib/apt/lists/*

# Create non-root user for security with home directory
RUN groupadd -r neuronagent && useradd -r -g neuronagent -u 1000 -m -d /home/neuronagent neuronagent && \
    chmod 755 /home/neuronagent && \
    chown -R neuronagent:neuronagent /home/neuronagent

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/neuronagent .

# Copy configs directory if it exists
COPY --from=builder /build/configs ./configs

# Copy optional entrypoint script (from build context root)
COPY docker/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh && \
    chown root:root /usr/local/bin/docker-entrypoint.sh

# Change ownership of app directory to non-root user
RUN chown -R neuronagent:neuronagent /app

# Switch to non-root user
USER neuronagent

EXPOSE 8080

# Add health check (using netcat for TCP check)
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD nc -z localhost 8080 || exit 1

# Set entrypoint (can be overridden in docker-compose or runtime)
# The entrypoint script performs validation and then executes the binary
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["./neuronagent"]


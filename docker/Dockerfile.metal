# Placeholder for Metal (macOS) deployment.
# Metal is for Apple Silicon / macOS; production Linux images use Dockerfile or Dockerfile.cuda/rocm.
# For local macOS build: go build -o neuron-agent ./cmd/agent-server && ./neuron-agent
# This file keeps the manifest consistent; use the main Dockerfile for Linux production.
FROM golang:1.23-bookworm AS builder
WORKDIR /build
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o neuron-agent ./cmd/agent-server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata netcat-openbsd && rm -rf /var/lib/apt/lists/*
RUN groupadd -r neuronagent && useradd -r -g neuronagent -u 1000 -m neuronagent
WORKDIR /app
COPY --from=builder /build/neuron-agent .
RUN chown -R neuronagent:neuronagent /app
USER neuronagent
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 CMD nc -z localhost 8080 || exit 1
CMD ["./neuron-agent"]

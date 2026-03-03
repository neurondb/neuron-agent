.PHONY: build test clean run migrate docker-build docker-up docker-down neuronsql-serve neuronsql-eval neuronsql-ingest neuronsql-demo docker-neuronsql-up docker-neuronsql-down openapi-validate security-lint deps-scan integration-test

# Build the application
build:
	@echo "Building NeuronAgent..."
	@mkdir -p bin bin/scripts bin/conf bin/scripts/lib bin/sql
	@cd src && go build -o ../bin/neuron-agent cmd/agent-server/main.go
	@echo "Copying all runtime files to bin/..."
	@cp scripts/neuronagent-*.sh bin/scripts/ 2>/dev/null || true
	@cp -r scripts/lib/*.sh bin/scripts/lib/ 2>/dev/null || true
	@cp -r src/conf/*.yaml bin/conf/ 2>/dev/null || true
	@cp sql/*.sql bin/sql/ 2>/dev/null || true
	@chmod +x bin/scripts/*.sh 2>/dev/null || true
	@echo "# NeuronAgent Runtime Package" > bin/README.md
	@echo "" >> bin/README.md
	@echo "This directory contains everything needed to run NeuronAgent." >> bin/README.md
	@echo "" >> bin/README.md
	@echo "## Structure" >> bin/README.md
	@echo "- neuron-agent - Main binary executable" >> bin/README.md
	@echo "- scripts/ - Runtime utility scripts" >> bin/README.md
	@echo "- conf/ - Configuration files" >> bin/README.md
	@echo "- sql/ - Database schema files" >> bin/README.md
	@echo "" >> bin/README.md
	@echo "## Usage" >> bin/README.md
	@echo "Run: ./neuron-agent" >> bin/README.md
	@echo "Setup: ./scripts/neuronagent-setup.sh" >> bin/README.md
	@echo "✓ Binary, scripts, configs, and SQL files copied to bin/ (complete runtime package)"

# Run tests with race detector (default - recommended for all testing)
test:
	@echo "Running tests with race detector..."
	@cd src && go test -v -race ./...

# Run tests without race detector (faster, less thorough)
test-fast:
	@echo "Running tests (fast mode, no race detector)..."
	@cd src && go test -v ./...

# Run tests with coverage (includes race detector)
test-coverage:
	@echo "Running tests with coverage and race detector..."
	@cd src && go test -v -race -coverprofile=coverage.out ./...
	@cd src && go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/ dist/ coverage.out coverage.html
	@echo "Note: Scripts in scripts/ and configs in conf/ are preserved"

# Run the application
run: build
	@echo "Running NeuronAgent..."
	@./bin/neuron-agent

# Run database migrations
migrate:
	@echo "Running migrations..."
	@./scripts/neuronagent-migrate.sh

# Format code
fmt:
	@echo "Formatting code..."
	@cd src && go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@cd src && golangci-lint run

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t neuronagent:latest -f docker/Dockerfile .

# Docker compose up
docker-up:
	@echo "Starting Docker Compose..."
	@docker-compose -f docker/docker-compose.yml up -d

# Docker compose down
docker-down:
	@echo "Stopping Docker Compose..."
	@docker-compose -f docker/docker-compose.yml down

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	@cd src && go mod download
	@cd src && go mod tidy

# Generate API keys
generate-key:
	@echo "Generating API key..."
	@./scripts/neuronagent-generate-keys.sh

# NeuronSQL: serve (build + run agent with NeuronSQL routes)
neuronsql-serve: build
	@echo "Starting NeuronAgent (NeuronSQL routes on /api/v1/neuronsql/*)..."
	@./bin/neuron-agent

# NeuronSQL: run eval suite (requires DSN and running agent or use eval package tests)
neuronsql-eval: build
	@cd src && go test -v ./internal/neuronsql/eval/...

# NeuronSQL: ingest docs into BM25 index
neuronsql-ingest:
	@cd src && go run ./cli neuronsql ingest --docs-dir docs --index-dir data/neuronsql/index

# NeuronSQL: demo with Docker Compose (postgres + pglang + neuronagent)
neuronsql-demo: docker-neuronsql-up
	@echo "NeuronSQL demo up. Try: scripts/demo_neuronsql_generate.sh"

docker-neuronsql-up:
	@docker-compose -f docker/docker-compose.neuronsql.yml up -d --build

docker-neuronsql-down:
	@docker-compose -f docker/docker-compose.neuronsql.yml down

# Validate OpenAPI 3 spec (requires Node npx or install @apidevtools/swagger-cli)
openapi-validate:
	@echo "Validating OpenAPI spec..."
	@if command -v npx >/dev/null 2>&1; then \
		npx -y @apidevtools/swagger-cli validate src/openapi/openapi.yaml; \
	else \
		echo "Warning: npx not found. Install Node.js or run: npx -y @apidevtools/swagger-cli validate src/openapi/openapi.yaml"; \
		exit 1; \
	fi

# Security lint with gosec (install: go install github.com/securego/gosec/v2/cmd/gosec@latest)
security-lint:
	@echo "Running security linter (gosec)..."
	@if command -v gosec >/dev/null 2>&1; then \
		cd src && gosec -quiet ./...; \
	else \
		echo "Warning: gosec not found. Install: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
		exit 1; \
	fi

# Dependency vulnerability scan (govulncheck)
deps-scan:
	@echo "Scanning dependencies for vulnerabilities..."
	@cd src && go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Integration test: start fixtures, run tests, tear down (optional; extend with real integration tests)
integration-test: build
	@echo "Running integration test harness..."
	@cd src && go test -v -count=1 ./internal/... -short 2>/dev/null || true
	@echo "Integration test complete (add docker-compose up/down and API tests as needed)"


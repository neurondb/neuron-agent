# Quickstart

## Prerequisites

- Docker with Compose v2 **or** Go 1.24+, PostgreSQL, `psql`

## Docker Compose (recommended for local)

From the repository root:

```bash
cp .env.example .env
docker compose build neuronagent
docker compose up -d
bash scripts/wait-for-health.sh
bash scripts/smoke-test.sh
```

Apply database schema (host needs network access to Postgres; default host port `5433`):

```bash
export DB_HOST=127.0.0.1 DB_PORT=5433 DB_NAME=neurondb DB_USER=neurondb DB_PASSWORD=neurondb
make migrate
```

Generate an API key (requires DB connectivity):

```bash
cd src && go run ./cmd/generate-key -db-host 127.0.0.1 -db-port 5433 -db-name neurondb -db-user neurondb -db-pass neurondb
```

Put the key in `.demo.env` as `NEURON_AGENT_API_KEY=...` and run `bash scripts/demo.sh`.

## Native binary

```bash
make build
cp .env.example .env
# edit .env for DB_* then:
make migrate
./bin/neuron-agent
```

Verify: `curl -s http://127.0.0.1:8080/health` and `curl -s http://127.0.0.1:8080/version`.

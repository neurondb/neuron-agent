# Docker Compose deployment

## Root stack (`docker-compose.yml`)

Services:

- **postgres** — PostgreSQL 17; data in named volume `neuronagent_pgdata`
- **neuronagent** — NeuronAgent image built from `docker/Dockerfile`

Variables are driven by `.env` (copy from `.env.example`). Host ports:

- Postgres: `POSTGRES_HOST_PORT` (default `5433`)
- HTTP: `SERVER_HOST_PORT` (default `8080`)

## NeuronDB

Full vector/RAG features require the **NeuronDB** extension in PostgreSQL. The stock `postgres:17-alpine` image does not include it.

### Option A — Build Postgres from a local NeuronDB checkout (recommended for development)

If you have the [NeuronDB](https://github.com/neurondb/neurondb) repository cloned next to this repo (or anywhere on disk), merge **`docker-compose.neurondb.yml`** so `postgres` is built from `docker/neurondb/Dockerfile`:

```bash
# Default NEURONDB_REPO=../neurondb (sibling directory)
cp .env.example .env
# Or set explicitly: NEURONDB_REPO=/path/to/neurondb

make up-neurondb
# equivalent:
# docker compose -f docker-compose.yml -f docker-compose.neurondb.yml up -d --build
```

Build only the database image:

```bash
make build-neurondb-pg
```

The first build compiles the extension and can take several minutes. If you previously ran vanilla Postgres against the same compose volume, **`docker compose ... down -v`** before switching images.

### Option B — Pre-built image

Use a NeuronDB-published image and point `DB_*` at it, or install the extension in your own image before applying `sql/neuron-agent.sql`.

## Commands

```bash
docker compose up -d
docker compose logs -f neuronagent
docker compose down -v   # removes volumes
```

## Health

- Liveness: `GET /health`, `GET /healthz`
- Readiness (DB ping): `GET /readyz`

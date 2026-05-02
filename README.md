# NeuronAgent

**Database-native agent runtime** on PostgreSQL and NeuronDB: durable memory, RAG, SQL-native tooling, workflows, and audit-friendly APIs. Not a chat toy — it is the backend for reliable agent applications.

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-blue.svg)](https://www.postgresql.org/)
[![Version](https://img.shields.io/badge/version-3.0.0--devel-blue.svg)](https://github.com/neurondb/neuron-agent)
[![License](https://img.shields.io/badge/License-See%20LICENSE-lightgrey.svg)](LICENSE)

**Hosted docs:** [neurondb.ai — NeuronAgent](https://www.neurondb.ai/docs/neuronagent)  
**Roadmap:** [docs/roadmap/6-month-world-class-plan.md](docs/roadmap/6-month-world-class-plan.md)

---

## Install (Docker Compose)

```bash
git clone https://github.com/neurondb/neuron-agent.git
cd neuron-agent
cp .env.example .env
docker compose build neuronagent && docker compose up -d
bash scripts/wait-for-health.sh
bash scripts/smoke-test.sh
```

Apply SQL schema with `make migrate` (see [Quickstart](docs/quickstart.md)). **NeuronDB** (vectors/RAG) needs Postgres with the extension built in: clone [neurondb/neurondb](https://github.com/neurondb/neurondb), set `NEURONDB_REPO` in `.env`, then run **`make up-neurondb`** (see [Docker Compose](docs/deploy/docker-compose.md)). Plain `postgres:17-alpine` is fine for smoke tests without vectors/RAG.

---

## First demo

```bash
bash scripts/bootstrap-demo.sh
# Create DB + API key (see docs/quickstart.md), then:
bash scripts/demo.sh
```

---

## What you get

- REST API (`/api/v1/*`) for agents, sessions, messages, memory, tools, workflows, RAG ingest, NeuronSQL routes
- WebSocket `/ws` for streaming (authenticated)
- **OpenClaw-compatible** bridge: `GET /claw/v1/health`, `POST /claw/v1/tools/list`, `POST /claw/v1/tools/run`
- Metrics `GET /metrics`, ops probes `GET /health`, `/healthz`, `/readyz`, build info `GET /version`
- API docs **`GET /docs`** (Redoc) and **`GET /docs/openapi.yaml`**

---

## Why NeuronAgent

| Layer | Role |
|--------|------|
| **OpenClaw** | User-facing channels (Slack, web, …) |
| **NeuronAgent** | Transactional agent brain: memory, tools, workflows, audit |
| **NeuronDB + PostgreSQL** | Database-native embeddings, RAG, SQL |

---

## Use cases

- **Customer support** — grounded replies with policy or ticket RAG, HTTP/SQL tools, audit trail
- **Internal research & ops** — sessions with memory tiers, workflows, and human approval steps
- **SQL & data copilots** — NeuronSQL routes and safe database tooling against PostgreSQL
- **Channel-connected agents** — OpenClaw (or similar) for UX; NeuronAgent for durable state and tools
- **Scheduled automation** — workflow DAGs with cron schedules, retries, and execution history
- **Platform / multi-tenant** — API keys, roles, metrics, and standard HTTP probes for production

---

## Examples

- [examples/quickstart-chat/](examples/quickstart-chat/)
- [examples/rag-with-docs/](examples/rag-with-docs/)
- [examples/sql-agent/](examples/sql-agent/)

Templates: [templates/](templates/)

---

## Repository layout

- **`src/`** — Go module (`github.com/neurondb/NeuronAgent`), HTTP server in `cmd/agent-server`
- **`sql/`** — Schema files for `make migrate`
- **`docker/`** — `Dockerfile`, `docker-compose.neuronsql.yml` (NeuronSQL demo stack)
- **`docs/`** — Local docs; start at [docs/index.md](docs/index.md)

Build locally: `make build` → `./bin/neuron-agent`

---

## OpenClaw bridge

Channel gateways can integrate via **`/claw/v1`** (health, tool list/run for NeuronSQL-style tools). Scope, roles, and wiring are documented in [docs/integrations/openclaw.md](docs/integrations/openclaw.md).

---

## Production deployment

- **Docker Compose** — [docs/deploy/docker-compose.md](docs/deploy/docker-compose.md) (vanilla Postgres, or **`make up-neurondb`** with a local NeuronDB checkout)
- **Kubernetes** — [docs/deploy/kubernetes.md](docs/deploy/kubernetes.md)

---

## Documentation

- **Local docs hub:** [docs/index.md](docs/index.md) — quickstart, configuration, architecture, security, workflows, RAG, observability, roadmap [implementation status](docs/roadmap/implementation-status.md)
- **Hosted:** [neurondb.ai — NeuronAgent](https://www.neurondb.ai/docs/neuronagent)

---

## Contributing & security

- [CONTRIBUTING.md](CONTRIBUTING.md)
- [SECURITY.md](SECURITY.md)
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

---

## License

See [LICENSE](LICENSE).

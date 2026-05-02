# NeuronAgent documentation

NeuronAgent is a **database-native agent runtime** for PostgreSQL and NeuronDB: durable memory, tools, workflows, RAG, and audit-friendly operations.

## Start here

- [Quickstart](quickstart.md) — install, migrate, first request
- [Configuration](configuration.md) — environment variables
- [Deploy with Docker Compose](deploy/docker-compose.md) — vanilla Postgres or **NeuronDB-backed Postgres** (`make up-neurondb`, `NEURONDB_REPO`)
- [Six-month execution plan (roadmap)](roadmap/6-month-world-class-plan.md)

## Deep dives

- [Architecture](architecture.md)
- [Concepts](concepts.md)
- [Memory](memory.md)
- [RAG](rag.md)
- [Tools](tools.md)
- [Workflows](workflows.md)
- [Security](security.md)
- [Observability](observability.md)
- [Backup & restore](backup-restore.md)
- [Reliability](reliability.md)
- [Development](development.md)
- [FAQ](faq.md)
- [Examples index](examples.md)
- [HTTP API](api.md) — OpenAPI at **`GET /docs`** on a running server

### Deploy & integrations

- [Docker Compose](deploy/docker-compose.md) — includes merge file `docker-compose.neurondb.yml` for building Postgres from a local [NeuronDB](https://github.com/neurondb/neurondb) repo
- [Kubernetes](deploy/kubernetes.md)
- [OpenClaw bridge](integrations/openclaw.md)

### Benchmarks & roadmap

- [Benchmarks](benchmarks.md)
- [Six-month execution plan](roadmap/6-month-world-class-plan.md)
- [Implementation status](roadmap/implementation-status.md)

Hosted mirror: [neurondb.ai — NeuronAgent](https://www.neurondb.ai/docs/neuronagent)

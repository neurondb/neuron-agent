# NeuronAgent Overview

<div align="center">

**AI agent runtime with REST API and WebSocket for autonomous agents, persistent memory, and tool execution**

[![Status](https://img.shields.io/badge/status-stable-brightgreen)](.)
[![API](https://img.shields.io/badge/API-REST%20%7C%20WebSocket-blue)](.)
[![Tools](https://img.shields.io/badge/tools-18+-green)](.)

</div>

---

## What NeuronAgent Is

NeuronAgent is an **AI agent runtime service** that provides:

- **REST API and WebSocket** — Create and manage agents, sessions, and messages over HTTP (default port 8080) and stream responses in real time over WebSocket.
- **Persistent memory** — Long-term memory backed by vector search (HNSW in NeuronDB) so agents can store and retrieve context across conversations.
- **Tool execution** — 18+ built-in tools (SQL, HTTP, code, shell, browser, memory, RAG, etc.) plus custom tool registration; agents choose tools via LLM decisions.
- **Workflow engine** — DAG-based workflows with agent, tool, HTTP, SQL, approval, and conditional steps; human-in-the-loop (HITL) with notifications.
- **Multi-agent collaboration** — Workspaces, agent-to-agent messaging, task delegation, and hierarchical agent structures.

NeuronAgent **depends on NeuronDB** (PostgreSQL with the NeuronDB extension). It uses the same database for agent state, sessions, messages, memory embeddings, and workflow state. It does not replace NeuronDB; it adds an agent runtime layer on top of it.

---

## Key Capabilities

| Area | Description | Status |
|------|-------------|--------|
| **Agent runtime** | State machine for autonomous task execution with persistent memory | ✅ Stable |
| **REST API** | Full CRUD for agents, sessions, messages, tools, memory, workflows, plans, budgets, collaborations | ✅ Stable |
| **WebSocket** | Real-time streaming of agent responses with event streaming | ✅ Stable |
| **Tool system** | 18+ built-in tools; extensible via custom tool registration | ✅ Stable |
| **Multi-agent** | Workspaces, agent-to-agent communication, task delegation, hierarchies | ✅ Stable |
| **Workflow engine** | DAG execution with HITL, retries, idempotency, scheduling | ✅ Stable |
| **Memory** | Hierarchical memory (working / episodic / semantic), HNSW search, promotion, feedback | ✅ Stable |
| **Integration** | NeuronDB for embeddings, LLM, vector search, ML, RAG | ✅ Stable |

---

## Quick Start

1. **Prerequisites:** PostgreSQL 16+ with NeuronDB extension, port 8080 available, API key for authentication.

2. **Start the service** (e.g. with Docker from the neurondb repo root):
   ```bash
   docker compose up -d neurondb
   docker compose up -d neuronagent
   ```

3. **Health check:**
   ```bash
   curl -sS http://localhost:8080/health
   # {"status":"ok"}
   ```

4. **Create an agent and send a message** (see [API Reference](api-reference.md)):
   ```bash
   export API_KEY="your_api_key"
   # Create agent
   curl -X POST http://localhost:8080/api/v1/agents \
     -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
     -d '{"name":"my-agent","system_prompt":"You are helpful.","model_name":"gpt-4","enabled_tools":[]}'
   # Create session (use agent id from response)
   curl -X POST http://localhost:8080/api/v1/sessions \
     -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
     -d '{"agent_id":"<agent_id>"}'
   # Send message (use session id from response)
   curl -X POST "http://localhost:8080/api/v1/sessions/<session_id>/messages" \
     -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
     -d '{"content":"Hello, what can you do?"}'
   ```

Full setup (database, config, Docker) is in the [README](../README.md).

---

## Documentation Map

| Document | Description |
|----------|-------------|
| [Features](features.md) | Full feature list: tools, memory, workflows, HITL, evaluation, security, etc. |
| [Architecture](architecture.md) | System design, components, database layer, integration with NeuronDB. |
| [API Reference](api-reference.md) | REST endpoints, request/response shapes, errors, WebSocket. |
| [API (detailed)](api.md) | Additional endpoints: workflows, plans, budgets, collaborations, evaluation, RAG, etc. |
| [Troubleshooting](troubleshooting.md) | Common issues: server, database, API, WebSocket, workflows, sandbox. |

**Machine-readable API:** OpenAPI 3.0 spec at `src/openapi/openapi.yaml` (or `/openapi.yaml` when the server is running).

---

## Docker Service Variants

| Service | Description |
|---------|-------------|
| **neuronagent** | Main service (CPU) |
| **neuronagent-cuda** | NVIDIA GPU variant |
| **neuronagent-rocm** | AMD GPU variant |
| **neuronagent-metal** | Apple Silicon GPU variant |

---

## Where to Go Next

- **First-time setup:** [README Quick Start](../README.md#quick-start) and [Configuration](../README.md#configuration).
- **Deep dive on capabilities:** [Features](features.md).
- **How the system is built:** [Architecture](architecture.md).
- **Calling the API:** [API Reference](api-reference.md) and [API (full)](api.md).
- **Problems:** [Troubleshooting](troubleshooting.md).

For the official NeuronDB docs site: [https://www.neurondb.ai/docs/neuronagent](https://www.neurondb.ai/docs/neuronagent).

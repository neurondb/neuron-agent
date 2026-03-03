# NeuronAgent Architecture

<div align="center">

**System architecture and design of the NeuronAgent runtime**

[![Architecture](https://img.shields.io/badge/architecture-documented-brightgreen)](.)
[![Status](https://img.shields.io/badge/status-stable-blue)](.)

</div>

---

## Overview

NeuronAgent is a Go service that provides an agent runtime on top of the NeuronDB PostgreSQL extension. It handles agent lifecycle, sessions, messages, tool execution, long-term memory (vector search), workflows, and multi-agent collaboration. All persistent state (agents, sessions, messages, memory chunks, workflows, API keys, etc.) is stored in PostgreSQL; vector operations and embeddings use NeuronDB extension types and functions.

---

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    NeuronAgent Service                        │
├─────────────────────────────────────────────────────────────┤
│  REST API Layer        │  WebSocket        │  Health/Metrics   │
├─────────────────────────────────────────────────────────────┤
│  Agent Runtime Engine  │  Session Manager  │  Tool Registry    │
├─────────────────────────────────────────────────────────────┤
│  Memory Manager        │  Planner          │  Reflector        │
├─────────────────────────────────────────────────────────────┤
│  Hierarchical Memory   │  Event Stream     │  VFS              │
├─────────────────────────────────────────────────────────────┤
│  Collaboration         │  Async Tasks      │  Sub-Agents        │
├─────────────────────────────────────────────────────────────┤
│  Verification Agent    │  Multimodal Proc  │  Browser Driver    │
├─────────────────────────────────────────────────────────────┤
│  Background Workers    │  Job Queue        │  Notifications     │
├─────────────────────────────────────────────────────────────┤
│              NeuronDB PostgreSQL Extension                    │
│  (Vector Search │ Embeddings │ LLM │ ML │ RAG)               │
└─────────────────────────────────────────────────────────────┘
```

---

## Core Components

### Database Layer (`internal/db/`)

All persistence goes through this layer.

- **Models** — Go structs for Agent, Session, Message, Tool, MemoryChunk, Workflow, and related entities.
- **Connection** — Pooled PostgreSQL connections with configurable limits and health checks.
- **Queries** — Prepared statements for security and performance.
- **Transactions** — Used for multi-step operations (e.g. store message + update session).
- **Migrations** — Schema evolution via migration files; NeuronAgent expects the `neurondb_agent` schema (or equivalent) and NeuronDB extension enabled.

### Agent Runtime (`internal/agent/`)

Orchestrates agent behavior and execution.

- **Runtime** — Main execution engine; drives the agent state machine (receive message → load context → call LLM → execute tools → store memory → respond).
- **MemoryManager** — Long-term memory using HNSW-based vector search (NeuronDB); store and retrieve by embedding similarity.
- **HierarchicalMemoryManager** — Multi-tier memory: working (current session), episodic (recent conversations), semantic (long-term). Promotion moves important content to higher tiers.
- **EventStreamManager** — Event logging with optional summarization for context windows.
- **LLMClient** — Calls NeuronDB LLM/embedding functions (or configured providers) for completions and embeddings.
- **ContextLoader** — Builds context from recent messages and retrieved memories.
- **PromptBuilder** — Assembles system prompt, tools, history, and current message for the LLM.
- **Planner** — Task planning and decomposition (LLM-based).
- **Reflector** — Self-reflection and response quality assessment.
- **VerificationAgent** — Optional output verification against configurable rules.
- **VirtualFileSystem (VFS)** — Per-agent file storage (DB and/or S3) with atomic operations.
- **AsyncTaskExecutor** — Background task execution with notifications.
- **SubAgentManager** — Routing and delegation to other agents.
- **TaskNotifier** — Alerts and approval notifications (email, webhook).
- **EnhancedMultimodalProcessor** — Image, audio, and code input processing.

### Tool System (`internal/tools/`)

Registry of built-in and custom tools.

- **Built-in tools** — SQL, HTTP, Code (Python sandbox), Shell (allowlist), Browser (Playwright), Memory, Filesystem, Visualization, NeuronDB tools (vector, RAG, ML, hybrid search, reranking, analytics), Multimodal, Web Search, Retrieval, Collaboration, etc. See [features.md](features.md).
- **Custom tools** — JSON Schema definitions; parameters validated and executed via the tool interface.
- **Analytics** — Tool usage and performance can be tracked per agent/tool.

### API Layer (`internal/api/`)

HTTP and WebSocket entrypoints.

- **REST** — Versioned under `/api/v1`: agents, sessions, messages, tools, memory, workflows, plans, budgets, collaborations, evaluation, and more. See [api-reference.md](api-reference.md) and [api.md](api.md).
- **WebSocket** — `/ws?session_id=...` for streaming agent responses; auth via query param or `Authorization` header.
- **Health / metrics** — `/health`, `/metrics` (Prometheus).

### Authentication (`internal/auth/`)

- **APIKeyManager** — API key storage and lookup (e.g. in `neurondb_agent.api_keys`).
- **Hasher / Validator** — Secure hashing (e.g. bcrypt) and validation.
- **RateLimiter** — Per-key rate limiting; returns 429 with `Retry-After` when exceeded.
- **PrincipalManager** — Identity and tenant context.
- **RBAC** — Role-based access control and fine-grained permissions.
- **AuditLogger** — Audit trail for sensitive actions.

### Background Jobs (`internal/jobs/`)

PostgreSQL-backed job queue.

- **MemoryPromoter** — Promotes memories across hierarchical tiers.
- **VerifierWorker** — Runs verification and quality checks.
- **AsyncTaskWorker** — Executes async tasks and sends notifications.
- **Cleanup workers** — Expiry and maintenance of old data.

---

## Agent Execution Flow

1. **Request** — Client sends a message via `POST /api/v1/sessions/{id}/messages` (or WebSocket).
2. **Context** — Runtime loads session history and retrieves relevant memories (vector search in NeuronDB).
3. **Prompt** — PromptBuilder builds the full prompt (system prompt, tools, history, retrieved context, user message).
4. **LLM** — LLMClient calls the configured model (via NeuronDB or external provider); may request tool calls.
5. **Tools** — Tool registry executes requested tools (SQL, HTTP, code, etc.); some tools call back into NeuronDB (vector, RAG, ML).
6. **Memory** — New information is embedded and stored in the memory manager (NeuronDB vector tables).
7. **Response** — Final answer (and optional tool results) is streamed to the client (WebSocket) or returned in the HTTP response.

Planning, reflection, and verification can be invoked as part of this flow or via dedicated API endpoints.

---

## Memory Hierarchy

- **Working memory** — Current session; recent turns and temporary context.
- **Episodic memory** — Recent conversations; short- to medium-term.
- **Semantic memory** — Long-term knowledge; vector embeddings, HNSW index in NeuronDB.

Promotion (e.g. via MemoryPromoter job) moves important content from working → episodic → semantic. Retrieval uses vector similarity (NeuronDB) to fetch relevant chunks for context.

---

## Workflow Engine

Workflows are **directed acyclic graphs (DAGs)** of steps. Step types include:

- **Agent** — Run an agent (LLM + tools).
- **Tool** — Run a specific tool.
- **HTTP** — External API call.
- **SQL** — Database query.
- **Approval** — Human-in-the-loop gate; can trigger email/webhook notifications.
- **Conditional** — Branch on previous step outputs.

Features: dependencies between steps, input/output mapping, retries, idempotency, optional compensation on failure, and cron-based scheduling. Workflow state and execution history are stored in PostgreSQL.

---

## Integration with NeuronDB

NeuronAgent uses the same PostgreSQL instance as NeuronDB and relies on:

- **Extension** — `CREATE EXTENSION neurondb;` must be enabled.
- **Schema** — Agent-specific tables (e.g. in `neurondb_agent`) for agents, sessions, messages, memory chunks, workflows, API keys, etc.
- **Vector search** — NeuronDB vector types and HNSW indexes for memory storage and similarity search.
- **Embeddings** — NeuronDB embedding/LLM functions for generating vectors when storing or querying memory.
- **LLM** — Completion and chat calls via NeuronDB LLM integration or configured external providers.
- **Tools** — SQL tool and NeuronDB-specific tools execute against this database (vector, RAG, ML, etc.).

No separate vector DB is required; NeuronDB is the single store for both extension features and agent state.

---

## Configuration

- **Environment variables** — Preferred in production (e.g. `DB_*`, `SERVER_*`, `LOG_*`, API key path).
- **Config file** — Optional YAML (e.g. `config.yaml`) for database, server, logging, auth, LLM, tools.
- **Defaults** — Sensible defaults for port (8080), pool sizes, timeouts; env overrides file.

---

## Monitoring and Observability

- **Metrics** — Prometheus endpoint (`/metrics`); scrape for dashboards and alerts.
- **Logging** — Structured (e.g. JSON) with request IDs and log levels.
- **Tracing** — Optional distributed tracing for request flow.
- **Health** — `/health` for liveness; can include DB and component checks.

---

## See Also

- [Overview](overview.md) — What NeuronAgent is and quick start.
- [Features](features.md) — Full list of tools, memory, workflows, and capabilities.
- [API Reference](api-reference.md) — REST and WebSocket API.
- [API (full)](api.md) — All endpoint groups and options.

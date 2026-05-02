# NeuronAgent: Six-Month World-Class Development Plan

> Last updated: 2026-04-30  
> Status: Active execution blueprint  
> Covers: May – October 2026

---

## Table of Contents

1. [Product Understanding](#1-product-understanding)
2. [Where NeuronAgent Fits](#2-where-neuronagent-fits)
3. [Six-Month Goal](#3-six-month-goal)
4. [Development Principles](#4-development-principles)
5. [Month 1: Foundation, Install, and Product Clarity](#5-month-1-foundation-install-and-product-clarity)
6. [Month 2: Core Runtime Hardening](#6-month-2-core-runtime-hardening)
7. [Month 3: Enterprise Readiness Baseline](#7-month-3-enterprise-readiness-baseline)
8. [Month 4: Workflows, RAG, and Integrations](#8-month-4-workflows-rag-and-integrations)
9. [Month 5: Developer Ecosystem and Adoption Assets](#9-month-5-developer-ecosystem-and-adoption-assets)
10. [Month 6: Production Readiness, Release, and Public Launch](#10-month-6-production-readiness-release-and-public-launch)
11. [Work Breakdown by Engineering Track](#11-work-breakdown-by-engineering-track)
12. [Measurable Success Metrics](#12-measurable-success-metrics)
13. [README Target Structure](#13-readme-target-structure)
14. [Documentation Target Structure](#14-documentation-target-structure)
15. [Final Summary](#15-final-summary)

---

## 1. Product Understanding

### What NeuronAgent is

NeuronAgent is a database-native agent runtime. It provides a structured backend for AI agents that need to remember things, retrieve information, reason over data, execute tools safely, run multi-step workflows, and maintain a complete audit trail. It runs as a server that your applications talk to over HTTP or WebSocket.

NeuronAgent is not a chatbot. It is not a wrapper around an LLM. It is not a prompt chain. It is the persistent, stateful, transactional layer that makes AI agents reliable in production.

When an AI agent needs to recall something from three sessions ago, NeuronAgent provides that memory. When an agent needs to answer a question using company documents, NeuronAgent provides the RAG pipeline. When an agent needs to inspect a database schema before writing a query, NeuronAgent provides the SQL tooling. When an enterprise needs to know exactly what an agent did and why, NeuronAgent provides the audit log.

### What NeuronDB provides

NeuronDB is a PostgreSQL extension that adds AI-native capabilities directly inside the database engine. It provides:

- **Vector embeddings**: generate, store, and search embeddings using `neurondb_embed()` and `<=>` cosine distance
- **RAG functions**: chunk documents, retrieve context, rerank results, and generate answers with `neurondb_chunk_text()`, `neurondb_rerank_results()`, `neurondb_generate_answer()`
- **Hybrid search**: combine semantic and keyword search with `neurondb_hybrid_search()` and reciprocal rank fusion
- **ML primitives**: `neurondb.predict()`, training APIs, and model management inside SQL
- **LLM calls from SQL**: call language models directly from queries

NeuronDB turns PostgreSQL into an AI-capable substrate. NeuronAgent uses NeuronDB as its intelligence layer.

### Why PostgreSQL matters

PostgreSQL is the most trusted open-source relational database. Organizations that already run PostgreSQL get NeuronAgent's full capabilities without adopting a new database technology. Their data stays where it is. Their backup, replication, access control, and operational tooling all apply. ACID guarantees mean agent state transitions are atomic. Foreign keys and schemas mean data is structured and trustworthy. Point-in-time recovery means agent memory is recoverable after failure.

PostgreSQL is not a trend. It is infrastructure that enterprises trust for critical data. Building agent memory on top of it is the right engineering decision.

### Why agents need durable memory

A stateless AI agent forgets everything between requests. It cannot improve over time, cannot recall previous decisions, cannot track user preferences, and cannot detect contradictions in what it has learned. This is fine for a toy demo and fatal for production applications.

Durable memory means:
- An agent remembers what happened in past sessions with a specific user
- Memory persists across container restarts, deployments, and failures
- Memory can be searched, audited, corrected, and deleted under GDPR
- Memory quality improves over time through feedback and consolidation
- Memory is versioned so you can understand what an agent knew at any point

NeuronAgent provides three memory tiers (short-term, medium-term, long-term) with automatic promotion, importance scoring, temporal decay, and episodic memory for event sequences. All memory is stored in PostgreSQL. None of it is lost when the process restarts.

### Why RAG should live close to the database

Retrieval-Augmented Generation requires embedding documents, storing vectors, searching by similarity, reranking results, and generating answers from retrieved context. Most architectures scatter these steps across multiple services: a document store here, a vector database there, an embedding service elsewhere, a generation model somewhere else.

When RAG lives inside the database (via NeuronDB), the benefits are:
- No network hops between chunking, embedding, storage, and retrieval
- ACID consistency between documents and their embeddings
- Unified backup covers both relational data and vector indexes
- Access control applies uniformly to all data
- No synchronization problems between the document and its vector representation

NeuronAgent's RAG pipeline uses NeuronDB functions to chunk documents, generate embeddings, and retrieve context entirely within PostgreSQL. The result is simpler architecture, fewer failure points, and operational consistency.

### Why enterprise users care about audit, permissions, reliability, backups, and observability

Enterprise deployments are evaluated on risk. An AI agent that cannot prove what it did is a liability. The questions enterprise evaluators ask are:

- Who authorized this agent to run?
- What did it access?
- What tools did it call?
- What did it write to memory?
- Can we reproduce what it did last Tuesday?
- What happens if the server crashes during a workflow?
- How do we back up agent memory?
- How do we know the agent is healthy right now?
- How do we revoke access if a key is compromised?
- Can we run this in our Kubernetes cluster?

NeuronAgent must answer every one of these questions with working code, not documentation promises.

---

## 2. Where NeuronAgent Fits

### The ecosystem

| Category | Examples | Role | NeuronAgent relationship |
|---|---|---|---|
| Personal assistant tools | ChatGPT, Claude apps | User-facing chat | Different layer entirely |
| Generic agent frameworks | LangChain, LlamaIndex | Orchestration libraries | NeuronAgent is a runtime server, not a library |
| Channel gateways | OpenClaw, Slack bots | User-facing message routing | OpenClaw calls NeuronAgent for intelligence |
| Vector databases | Pinecone, Weaviate, Qdrant | Vector storage only | NeuronDB replaces these inside PostgreSQL |
| SaaS chatbot builders | Intercom AI, Zendesk AI | Product features | Different audience |
| MLOps platforms | MLflow, Vertex AI | Model training and serving | Different concern |

### The positioning triangle

```
User / App
    │
    ▼
OpenClaw (channels, routing, user interfaces)
    │
    ▼
NeuronAgent (agent runtime: memory, RAG, tools, workflows, audit)
    │
    ▼
NeuronDB + PostgreSQL (database-native AI substrate)
```

**OpenClaw** handles user-facing channels: Slack, Teams, web chat, API gateways, message routing. It knows how to talk to users.

**NeuronAgent** handles the agent brain: durable memory, RAG, SQL-native reasoning, tool execution, multi-step workflows, permission enforcement, and complete audit logs. It knows how to think reliably.

**NeuronDB** handles the AI substrate: embeddings, vector search, hybrid retrieval, ML models, and LLM calls from inside PostgreSQL. It knows how to store and process intelligence.

NeuronAgent is the runtime layer between user applications and database-native intelligence. It is not a chat toy. It is the backend brain for reliable AI applications.

---

## 3. Six-Month Goal

By the end of October 2026, NeuronAgent must be the reference implementation of a database-native agent runtime. A developer who finds the project on GitHub must be able to install it, run a demo, and understand its value within five minutes.

### Six-month milestone table

| Month | Theme | Key deliverable |
|---|---|---|
| 1 | Foundation and install | One-command local stack, honest README, smoke test |
| 2 | Core runtime hardening | Config validation, agent lifecycle tests, tool audit logs |
| 3 | Enterprise baseline | Auth, RBAC, audit logs, observability, backup docs |
| 4 | Workflows, RAG, integrations | RAG demo, DBA agent, workflow demo, OpenClaw bridge |
| 5 | Developer ecosystem | SDK, templates, contributor docs, demo GIF |
| 6 | Release and launch | CI automation, Kubernetes, benchmarks, v1.0 release |

### End-state checklist

- [ ] One-command local install (`git clone` + `docker compose up -d` + `make demo`)
- [ ] Docker Hub and GitHub Container Registry images published
- [ ] Stable Docker Compose stack with NeuronDB-enabled PostgreSQL
- [ ] Clear README focused on first success in under two minutes
- [ ] Docs site structure with 20+ pages
- [ ] OpenAPI spec and Swagger UI at `/docs`
- [ ] Working demos for chat, RAG, SQL agent, and workflow
- [ ] At least five agent templates in `templates/`
- [ ] Stable API versioning under `/api/v1`
- [ ] Go unit test coverage above 70% for core packages
- [ ] Smoke test in CI on every pull request
- [ ] Enterprise security baseline: API keys, JWT, RBAC, audit logs
- [ ] Prometheus metrics, structured JSON logs, OpenTelemetry tracing hooks
- [ ] Backup and restore guide
- [ ] Production Docker Compose guide
- [ ] Kubernetes deployment path (Helm chart complete)
- [ ] OpenClaw bridge endpoints documented and tested
- [ ] At least five example agents with working scripts
- [ ] CONTRIBUTING.md, issue templates, PR template
- [ ] Automated release workflow (tag → build → publish → release notes)
- [ ] Benchmark script with local reproducible results
- [ ] Demo GIF generated by script
- [ ] GitHub issue templates for bugs, features, documentation
- [ ] No broken quickstart commands
- [ ] No feature claims without supporting code

---

## 4. Development Principles

These principles are not aspirational. They are constraints. Any contribution that violates them must be revised before merging.

1. **First success in under five minutes.** A new developer must be running NeuronAgent locally, with a working demo, in under five minutes from a clean machine with Docker installed.

2. **README sells and guides. Docs explain depth.** The README must answer "what is it, how do I install it, and what does it do." Deep architecture belongs in `docs/`.

3. **Every feature needs one working example.** A feature without an example is invisible. Every public capability must have at least one script or file that demonstrates it working.

4. **Every public claim needs code or a tracked issue.** If the README says NeuronAgent supports JWT auth, that means the code supports it. If the code does not exist yet, the claim must be removed or replaced with a link to the open issue.

5. **Install scripts must fail with clear messages.** If a prerequisite is missing, the install script tells the user what to install and where to get it. It does not silently fail or produce cryptic errors.

6. **Local mode should be easy.** Local development requires only Docker. No cloud accounts, no paid services, no manual configuration beyond copying `.env.example`.

7. **Production mode should be safe.** Production config validation must refuse weak secrets, default passwords, and unsafe defaults. The server must not start in production mode with insecure settings.

8. **Security defaults should be strict.** Authentication is required on all endpoints except `/health` and `/metrics`. CORS is locked down by default. Rate limiting is on by default.

9. **Logs should help users fix problems.** Error messages must explain what went wrong, why it happened, and what the user can do about it. Stack traces are for debugging, not for end users.

10. **Enterprise features should be real, not decorations.** Audit logs must be durable, queryable, and tamper-evident. RBAC must actually enforce permissions. Metrics must reflect real system state.

11. **Database state must be durable and recoverable.** Agent memory, session history, workflow state, and audit logs must survive crashes, restarts, and deployments. Backup and restore procedures must be documented and tested.

12. **Tests must protect the golden path.** The most important user journey (install → start → create agent → run demo) must be covered by an automated test that runs on every pull request.

13. **APIs need stable shapes and versioning.** Breaking changes to `/api/v1` endpoints require a deprecation notice, migration guide, and version bump. No silent breaking changes.

14. **Contributors need simple targets.** Every good-first-issue must have a clear description, clear acceptance criteria, and a clear location in the code. No vague tasks.

15. **Every month must produce demo-worthy progress.** Each month ends with something new that can be shown, not just refactored or documented.

---

## 5. Month 1: Foundation, Install, and Product Clarity

**Theme:** Make NeuronAgent installable, runnable, and understandable for the first time.

**Acceptance criteria:**
- A new user installs locally from README instructions, no tribal knowledge required
- Demo runs end to end on a clean machine with only Docker installed
- NeuronDB and NeuronAgent start together via `docker compose up -d`
- README explains the project value in under two minutes
- Health endpoints respond correctly
- Smoke test catches regressions in the golden path
- No manual environment setup required for local demo

---

### Week 1: Repository Audit and Truth Map

**Goal:** Establish ground truth. Know exactly what works, what is broken, and what is missing.

**Tasks:**
- Map current architecture by reading `src/cmd/agent-server/main.go` and `src/internal/`
- List all live API endpoints from `buildRouter`
- List all agent runtime features from `src/internal/agent/`
- List all memory features and tiers from `src/internal/agent/memory*.go`
- List all registered tool handlers in `NewRegistry` vs `NewRegistryWithNeuronDB`
- List workflow step types from `src/internal/workflow/engine.go`
- List all NeuronDB function calls from `src/pkg/neurondb/`
- Identify missing install steps by running a fresh clone
- Identify broken Docker path (entrypoint binary name mismatch)
- Identify broken compose stack (plain postgres vs NeuronDB-enabled postgres)
- List missing or outdated documentation
- List Go unit test coverage gaps in `src/internal/api/`
- List CI gaps (workflow_dispatch only, no push triggers)
- List security gaps from config defaults
- List observability gaps from metrics coverage

**Deliverables:**
- `docs/roadmap/current-state-audit.md`
- `docs/roadmap/gap-analysis.md`
- `docs/roadmap/feature-truth-table.md`

**Key files to inspect:**
- `src/cmd/agent-server/main.go`
- `src/internal/tools/registry.go`
- `docker/Dockerfile`
- `docker/docker-entrypoint.sh`
- `docker/docker-compose.neuronsql.yml`
- `Makefile`
- `scripts/`
- `.github/workflows/`

---

### Week 2: One-Command Local Stack

**Goal:** `git clone` + `docker compose up -d` + `make demo` must work on a clean machine.

**Tasks:**

*Docker fix (critical):*
- Fix `docker/Dockerfile`: binary name must match `docker/docker-entrypoint.sh`
- Create root-level `docker-compose.yml` with:
  - NeuronDB-enabled PostgreSQL service (custom image or init script that installs extension)
  - NeuronAgent service
  - Named volumes for postgres data
  - Health checks for both services
  - `depends_on` with condition `service_healthy`
  - `.env` file loading

*Scripts:*
- `scripts/install.sh` — checks Docker version, clones repo, copies `.env.example` to `.env`, runs `docker compose up -d`, waits for health, runs `make demo`
- `scripts/wait-for-health.sh` — polls `/health` endpoint until ready or timeout
- `scripts/bootstrap-demo.sh` — creates API key, workspace, agent, ingests a fixture document
- `scripts/reset.sh` — stops stack, removes volumes, restores clean state
- `scripts/logs.sh` — wraps `docker compose logs -f`

*Environment:*
- `.env.example` at repo root with all required variables, safe local defaults, no production secrets
- Auto-generate local API key in bootstrap script

*Makefile targets:*
- `make up` — `docker compose up -d`
- `make down` — `docker compose down`
- `make logs` — `docker compose logs -f`
- `make reset` — `docker compose down -v && docker compose up -d`
- `make demo` — runs `scripts/bootstrap-demo.sh` then `scripts/demo.sh`
- `make smoke` — runs golden path smoke test
- `make test` — runs Go unit tests in `src/`
- `make lint` — runs `golangci-lint` in `src/`
- `make docker-build` — builds the Docker image locally

*Migration fix:*
- Align migration naming: use `001_initial.sql` pattern
- Add migration runner call to docker-entrypoint or compose init

**Target install flow:**
```bash
# Option A: script install
curl -fsSL https://raw.githubusercontent.com/neurondb/neuron-agent/main/scripts/install.sh | bash

# Option B: manual
git clone https://github.com/neurondb/neuron-agent.git
cd neuron-agent
cp .env.example .env
docker compose up -d
make demo
```

---

### Week 3: README Rewrite

**Goal:** A visitor who has never heard of NeuronAgent reads the README and thinks "I understand it, I want to try it."

**README structure:**
1. Project statement (two sentences)
2. Badges (only for real, passing workflows)
3. One-command install
4. First demo (GIF placeholder or terminal output)
5. What you get (feature list, no architecture fog)
6. Why NeuronAgent (positioning vs alternatives)
7. Use cases (four to six concrete scenarios)
8. Examples section (links to `examples/`)
9. OpenClaw bridge (one paragraph, link to docs)
10. Production deployment (link to docs)
11. Docs
12. Contributing
13. Security
14. License

**Rules for the README rewrite:**
- No architecture diagrams in README (move to `docs/architecture.md`)
- No internal package names in README
- No badges for workflows that do not exist
- No feature claims that are not backed by working code
- Every command in README is tested before merging

**Move to docs:**
- `docs/architecture.md` — all architecture diagrams
- `docs/api.md` — full API reference
- `docs/configuration.md` — all environment variables

---

### Week 4: Golden Path Demo and Smoke Test

**Goal:** A working demo that covers every major feature and a CI smoke test that catches regressions.

**Demo script (`scripts/demo.sh`):**
1. Wait for health endpoint to respond
2. Create an API key using `/api/v1/admin/keys` or generate-key binary
3. Create a workspace (or use default)
4. Create an agent with a basic system prompt
5. Ingest a fixture document (`tests/fixtures/sample-doc.txt`)
6. Send a question about the document
7. Print the agent's answer
8. Run a safe read-only SQL query via the SQL tool
9. Print pass/fail for each step

**New examples:**
- `examples/quickstart-chat/` — minimal agent creation and chat
- `examples/rag-with-docs/` — document ingest and RAG question
- `examples/sql-agent/` — SQL schema inspection and read-only query

**Smoke test:**
- `make smoke` runs the golden path against a live stack
- Verifies HTTP status codes and response shape
- Fails with clear error messages if any step breaks
- Added to CI as a nightly job initially; promote to PR gate when stack startup is fast enough

---

## 6. Month 2: Core Runtime Hardening

**Theme:** Make the runtime predictable. Define clear behavior for every public API. Add tests that catch regressions.

**Acceptance criteria:**
- Core runtime behaves predictably under error conditions
- Agent create/read/update/delete APIs have Go unit tests
- Memory system has documented write, search, promote, forget, and conflict paths
- Tool execution enforces permissions and writes audit log entries
- Startup config errors produce human-readable messages

---

### Week 1: Configuration and Startup Reliability

**Goal:** The server never starts in a broken state without telling the user why.

**Tasks:**
- Create `.env.example` with every `DB_*`, `AUTH_*`, `SERVER_*`, `LOG_*` variable and inline comments
- Add strict config validation for production mode in `src/internal/config/config.go`:
  - Refuse `AUTH_API_KEY_SECRET=""` in production
  - Refuse default DB passwords in production
  - Require `SERVER_MODE` to be one of `local`, `development`, `test`, `production`
- Add clear startup error messages: "Missing required config: AUTH_API_KEY_SECRET. Set this in .env or environment."
- Print version summary at startup (without secret values):
  ```
  NeuronAgent v3.0.0-devel [production]
  Git commit: abc1234 | Build: 2026-04-30
  API: v1 | DB: neurondb@localhost:5432/neuronagent
  Migrations: v12 | NeuronDB: compatible
  ```
- Add `/version` endpoint returning JSON:
  ```json
  {
    "version": "3.0.0-devel",
    "git_commit": "abc1234",
    "build_date": "2026-04-30",
    "api_version": "v1",
    "neurondb_compatible": true,
    "migration_version": 12,
    "runtime_mode": "production"
  }
  ```
- Add `docs/configuration.md` with all environment variables, types, defaults, and examples

---

### Week 2: Agent Runtime Stability

**Goal:** Agent lifecycle is well-defined, validated, and tested.

**Tasks:**
- Document agent model: id, workspace, name, system_prompt, status, tool_permissions, config
- Define agent states: `active`, `paused`, `archived`
- Add validation on agent creation: system_prompt max length, name uniqueness per workspace
- Add default system prompt if none provided
- Add tool permission model: agents declare which handler types they are allowed to use
- Add workspace boundary enforcement: agents cannot access sessions or memory from other workspaces
- Add session cleanup: orphaned sessions older than configurable TTL are marked expired
- Add deterministic error responses: `{"error": "...", "code": "...", "request_id": "..."}`
- Add Go unit tests for agent CRUD lifecycle in `src/internal/api/`

---

### Week 3: Memory System Polish

**Goal:** Memory behavior is documented, predictable, and tested.

**Tasks:**
- Write `docs/memory.md` covering:
  - Three-tier model: STM (session-scoped), MTM (agent-scoped, days), LPM (persistent, weeks+)
  - Write path: when and how memories are created
  - Search path: embedding query, cosine distance, tier selection
  - Promotion path: STM → MTM → LPM based on access frequency and importance
  - Forget path: explicit delete API, TTL expiry, GDPR delete
  - Conflict behavior: what happens when contradictory memories exist
- Add memory quality fields to the API: importance score, access count, last accessed, source
- Add memory metadata fields: source session, source message, tags
- Add memory audit log entry on write and delete operations
- Add Go unit tests for `MemoryManager.Store`, `MemoryManager.Search`, `MemoryManager.Forget`
- Add integration test: store memory, search for it, verify similarity score > threshold

---

### Week 4: Tool Execution Reliability

**Goal:** Every tool call is permission-checked, timed out, audited, and returns a structured result.

**Tasks:**
- Document tool registry design in `docs/tools.md`
- Fix `NewRegistry` in `src/internal/tools/registry.go` to optionally include NeuronDB handlers when DB client is available — `agent-server` should use `NewRegistryWithNeuronDB` when NeuronDB client is initialized
- Fix `RetrievalTool` registration when agentic retrieval is enabled
- Define permission check flow: agent tool_permissions → registry handler → execute
- Define SQL tool as read-only by default; require explicit `allow_writes: true` in tool config
- Add HTTP tool domain allowlist and request size limit
- Define shell tool policy: disabled by default, requires explicit opt-in per agent
- Add timeout enforcement: every tool execution respects `timeout_ms` from tool config
- Add structured tool result schema: `{"success": bool, "result": ..., "error": ..., "duration_ms": int, "tool_id": "..."}`
- Add audit log entry for every tool execution: tool_id, agent_id, session_id, status, duration, input hash
- Add Go unit tests: success path, error path, timeout, permission denied, malformed input

---

## 7. Month 3: Enterprise Readiness Baseline

**Theme:** An enterprise evaluator must be able to tick every box on their security, observability, and operations checklist.

**Acceptance criteria:**
- Auth and RBAC are fully implemented and tested with permission-denied test cases
- Audit log schema is durable and every defined action type produces an entry
- Prometheus metrics cover all defined counters and histograms
- Graceful shutdown, backup docs, and restore docs are complete
- A production checklist exists

---

### Week 1: Auth and Permissions

**Goal:** Every API call is authorized against a role-based permission model.

**Tasks:**
- Review current auth: `src/internal/auth/`, `src/internal/api/middleware.go`
- Confirm API key auth covers all `/api/v1` routes
- Add JWT auth support if config `AUTH_JWT_SECRET` is set (for service-to-service)
- Add workspace-scoped permissions: a key can be scoped to one or more workspaces
- Implement role model in `auth.Principal`:
  - `owner` — full control including delete workspace
  - `admin` — full control, no workspace delete
  - `developer` — create and run agents, read logs
  - `viewer` — read-only access to agents, sessions, memory
  - `service` — tool execution and session creation only (for OpenClaw integration)
- Add permission checks for: agent create/update/delete, tool run, memory write, memory read, workflow run, admin actions
- Add `docs/security.md` covering: auth methods, role model, key management, rotation, workspace scoping
- Add Go unit tests for permission denied cases on every protected resource type

---

### Week 2: Audit Logging

**Goal:** Every security-relevant action is recorded durably and queryable by admins.

**Tasks:**
- Confirm or create durable audit log schema: `neurondb_agent.audit_log` with columns: id, timestamp, actor_id, actor_type, workspace_id, action, resource_type, resource_id, request_id, status, reason, metadata JSONB
- Ensure the following events produce audit entries:
  - API key usage (login equivalent)
  - Auth failure (rate limit hit, invalid key, wrong role)
  - Agent create / update / delete
  - Session create
  - Message send
  - Tool execute (success and failure)
  - Memory write / delete
  - Workflow start / approve / fail
  - Admin config change
- Audit entry metadata must never include secrets, API keys, or PII beyond workspace-scoped identifiers
- Add `GET /api/v1/admin/audit` — paginated audit log read for admin role
- Add audit log retention policy config: `AUDIT_LOG_RETENTION_DAYS`
- Add Go unit tests for audit entry creation on every defined action type

---

### Week 3: Observability

**Goal:** A production NeuronAgent deployment can be monitored with standard tooling.

**Tasks:**
- Confirm `/health` and `/metrics` endpoints work (already in code, verify)
- Add `/healthz` as alias for Kubernetes liveness probe
- Add `/readyz` for Kubernetes readiness probe (returns 503 until DB is connected and migrations are applied)
- Confirm all Prometheus metrics are wired via `prometheus/client_golang`:
  - `neuronagent_http_requests_total{method, path, status}` — HTTP request counter
  - `neuronagent_http_request_duration_seconds{method, path}` — HTTP latency histogram
  - `neuronagent_agent_requests_total{agent_id, status}` — agent execution counter
  - `neuronagent_agent_request_duration_seconds` — agent execution latency
  - `neuronagent_tool_executions_total{handler_type, status}` — tool counter
  - `neuronagent_tool_failures_total{handler_type, reason}` — tool failure counter
  - `neuronagent_memory_writes_total` — memory store counter
  - `neuronagent_memory_search_duration_seconds` — memory search latency
  - `neuronagent_rag_requests_total{status}` — RAG counter
  - `neuronagent_rag_duration_seconds` — RAG latency
  - `neuronagent_workflow_runs_total{status}` — workflow counter
  - `neuronagent_workflow_failures_total{reason}` — workflow failure counter
  - `neuronagent_neurondb_query_duration_seconds{operation}` — NeuronDB latency
- Add `docs/observability.md` with:
  - Endpoint reference
  - Sample `prometheus.yml` scrape config
  - Key metrics and alert thresholds
  - Grafana dashboard JSON (at minimum a 6-panel starter dashboard)

---

### Week 4: Reliability and Operations

**Goal:** Operators have the documentation and tooling to run NeuronAgent in production safely.

**Tasks:**
- Confirm graceful shutdown path in `main.go` (signal handling exists, verify correctness)
- Add worker restart behavior: `MemoryPromoter` and other background workers restart with backoff on panic
- Add queue visibility endpoint: `GET /api/v1/admin/queue-stats` showing pending jobs
- Add idempotency key support for workflow starts: `POST /api/v1/workflows/{id}/execute` accepts `X-Idempotency-Key` header
- Add retry policy documentation for transient DB failures
- Add `docs/backup-restore.md` covering:
  - `pg_dump` for full backup
  - `pg_restore` for restore
  - Point-in-time recovery notes
  - What to back up: schema, memory chunks, sessions, audit logs
  - Estimated backup frequency guidance
- Add `docs/deploy/docker-compose-production.md` covering: environment variable hardening, volume mounts, network security, resource limits
- Add `docs/reliability.md` covering: graceful shutdown, worker lifecycle, queue behavior, idempotency, retry policy

---

## 8. Month 4: Workflows, RAG, and Integrations

**Theme:** The three major differentiated capabilities (RAG, SQL-native reasoning, workflows) must each have a working demo and clear documentation.

**Acceptance criteria:**
- RAG demo: ingest a document, ask a question, receive an answer with source reference
- DBA agent demo: inspect a schema, profile a table, explain a query
- Workflow demo: create a multi-step workflow with an approval gate, execute it
- OpenClaw bridge: all three endpoints respond correctly with auth
- Docs explain every integration clearly

---

### Week 1: RAG Quality Path

**Goal:** Document ingest-to-answer works reliably and is fully documented.

**Tasks:**
- Stabilize `POST /api/v1/rag/ingest` handler: accept `{title, content, metadata, chunk_size, chunk_overlap}`
- Add `GET /api/v1/rag/query` endpoint: accept `{question, top_k, workspace_id}`, return answer with sources
- Support document metadata: `source_url`, `author`, `created_at`, `tags`
- Expose chunking controls: `chunk_size` (default 512), `chunk_overlap` (default 50)
- Use `RAGClient.ChunkDocument` → embed → store in knowledge table via NeuronDB
- Add source citation in answer: return `sources: [{title, chunk_text, score}]` alongside answer
- Add RAG tests using a small, fixed fixture file (`tests/fixtures/rag-sample.txt`) to avoid LLM dependency in tests — test chunking, retrieval shape, and response schema
- Add `examples/rag-with-docs/` with README, sample document, and working demo script
- Add `docs/rag.md` with: architecture diagram, API reference, chunking explanation, retrieval explanation, reranking explanation

---

### Week 2: SQL-Native Reasoning and DBA Agent

**Goal:** Agents can safely inspect and reason over PostgreSQL schemas.

**Tasks:**
- Confirm NeuronSQL module tools are wired correctly in `agent-server`
- Add or confirm read-only SQL tool: default SQL tool must refuse DDL and DML unless explicitly allowed per agent config
- Add schema summary tool: `neuronsql.schema_snapshot` — returns table names, columns, types, indexes as structured JSON
- Add table profile tool: `neuronsql.table_profile` — row count, column stats, null rates, sample values
- Add query explain tool: `neuronsql.explain_json` — runs `EXPLAIN (FORMAT JSON)` on a query
- Add index suggestion helper: analyze table profile + explain output, suggest missing indexes
- Ensure `RegisterNeuronSQLTools` in `src/internal/neuronsql/tools/register.go` is called or its tools are otherwise reachable — currently dead code
- Add `examples/postgres-dba-agent/` with README, demo script, and sample queries
- Add `templates/postgres-dba-agent.yaml` — agent template with NeuronSQL tools enabled, safe read-only SQL, and DBA-focused system prompt
- Add `docs/examples/postgres-dba-agent.md`

---

### Week 3: Workflow Engine Polish

**Goal:** Workflows are reliable, observable, and have a working approval-gate demo.

**Tasks:**
- Review workflow engine: `src/internal/workflow/engine.go`
- Add `conditional` step type to the main engine (it exists in `advanced_engine.go` and TypeScript SDK types but not in the Go engine switch)
- Implement workflow retry scheduler: when a step fails and `retry_count < max_retries`, schedule re-execution after `retry_delay_ms`
- Implement schedule runner: add a background worker that calls `ListWorkflowSchedulesByNextRun` and enqueues scheduled workflows
- Fix agent step session handling in `executeAgentStep`: use a real session per workflow execution, not `uuid.Nil`
- Add workflow history: `GET /api/v1/workflows/{id}/executions` returns paginated execution history with step-level status
- Add `examples/workflow-approval/` — a two-step workflow: agent analyzes a request, human approves, agent executes
- Add `docs/workflows.md` with: step type reference, execution model diagram, retry and timeout explanation, approval gate guide, idempotency guide

---

### Week 4: OpenClaw Bridge

**Goal:** OpenClaw can call NeuronAgent as a tool provider with clear auth and permissions.

**Tasks:**
- Confirm `/claw/v1/health`, `/claw/v1/tools/list`, `/claw/v1/tools/run` exist in `buildRouter`
- Add auth check on all `/claw/v1` routes using `service` role API key
- Expose the following tools via the claw bridge:
  - `agent.chat` — send a message to a named agent, return response
  - `rag.answer` — ask a question over the knowledge base, return answer with sources
  - `rag.ingest` — ingest a document into the knowledge base
  - `memory.search` — search agent memory by query string
  - `schema.inspect` — return schema summary for a workspace database
  - `sql.read` — run a read-only SQL query (requires explicit allow in bridge config)
  - `workflow.trigger` — start a named workflow with input parameters
- Add permission mapping: claw tool calls go through the same permission model as direct API calls
- Add request logging: every claw tool call produces a structured log entry and audit event
- Add `docs/integrations/openclaw.md` covering: auth setup, available tools, permission model, example OpenClaw config
- Add `examples/openclaw-bridge/` with a sample OpenClaw config and test script

---

## 9. Month 5: Developer Ecosystem and Adoption Assets

**Theme:** Make the project easy to find, easy to understand, easy to contribute to, and easy to demo.

**Acceptance criteria:**
- A first-time visitor understands NeuronAgent's value within two minutes of opening the README
- A developer finds a working example for their use case within five minutes
- A potential contributor finds an actionable good-first-issue within ten minutes
- API docs are accessible at `/docs` locally
- Demo GIF is generated by a script, not manually created

---

### Week 1: SDK and API Experience

**Goal:** Developers can use NeuronAgent from Python and TypeScript with minimal friction.

**Tasks:**
- Review and validate `src/openapi/openapi.yaml` against live router endpoints in `buildRouter`
- Expose Swagger UI or Redoc at `/docs` (use embedded static assets or redirect to hosted spec)
- Generate and commit `openapi.json` as part of `make generate`
- Add curl examples to `docs/api.md` for the ten most important API calls
- Review existing Python client in `src/examples/neurondb_client/` — package and publish as `neuronagent-client`
- Review TypeScript SDK in `src/sdks/typescript/neurondb-agent/` — validate against actual API
- Add example scripts using both clients
- Add `docs/api.md` with: endpoint reference, curl examples, Python examples, TypeScript examples, error codes

---

### Week 2: Templates and Examples

**Goal:** Five agent templates that work out of the box.

**Tasks:**
- Create `templates/postgres-dba-agent.yaml` — schema inspection, query explain, read-only SQL
- Create `templates/rag-knowledge-agent.yaml` — document ingest, RAG question answering, source citations
- Create `templates/customer-support-agent.yaml` — memory-enabled support with escalation workflow
- Create `templates/research-agent.yaml` — multi-source retrieval, synthesis, memory consolidation
- Create `templates/workflow-agent.yaml` — multi-step workflow with approval gate
- Add `make load-template TEMPLATE=postgres-dba-agent` command
- Add `docs/templates.md` with: template reference, how to load, how to customize
- Add sample prompts and expected outputs for each template
- Add cleanup scripts for each example

---

### Week 3: Contributor Experience

**Goal:** A new contributor can find a task, understand the codebase, and submit a PR without asking questions.

**Tasks:**
- Add `CONTRIBUTING.md`: setup guide, development workflow, test guide, PR process, coding standards
- Add `CODE_OF_CONDUCT.md`: Contributor Covenant
- Add `SECURITY.md`: vulnerability reporting process, response SLA, supported versions
- Add `SUPPORT.md`: community channels, bug reporting, feature requests
- Update `ROADMAP.md`: link to `docs/roadmap/6-month-world-class-plan.md`
- Update `CHANGELOG.md`: add first versioned release entry
- Add `.github/ISSUE_TEMPLATE/bug_report.yml`
- Add `.github/ISSUE_TEMPLATE/feature_request.yml`
- Add `.github/ISSUE_TEMPLATE/documentation.yml`
- Add `.github/PULL_REQUEST_TEMPLATE.md`
- Add `.github/labels.yml` with standard label set
- Add `docs/development.md`: architecture overview for contributors, package structure guide, how to add a new tool handler, how to add a new API endpoint

---

### Week 4: Demo and Marketing Assets

**Goal:** The project makes a strong first impression and every use case has a visual asset.

**Tasks:**
- Add `scripts/record-demo.sh` using VHS (if installed) or asciinema — records the full demo flow
- Generate `docs/assets/demo.gif` from the recording script
- Add demo GIF to README (replace placeholder)
- Add `docs/use-cases.md` with six concrete scenarios:
  1. AI knowledge assistant for internal documentation
  2. PostgreSQL DBA agent for database analysis
  3. Internal support agent with escalation
  4. Enterprise RAG over private documents
  5. Research assistant with memory and synthesis
  6. Workflow automation with human approval gates
- Add `docs/comparisons.md`: honest comparison with LangChain, LlamaIndex, OpenClaw (complementary, not competitive), vector databases, SaaS chatbots
- Add `docs/faq.md`: 20 most common questions with direct answers
- Add `docs/index.md`: product page style intro — two paragraphs, key features, quick start link

---

## 10. Month 6: Production Readiness, Release, and Public Launch

**Theme:** Ship it. Everything must work from a clean machine. Automation covers build, test, and release.

**Acceptance criteria:**
- README install instructions work from a clean machine with only Docker installed
- Docker images are published to Docker Hub and GHCR on every tagged release
- Docs are complete and all links work
- Security baseline is documented and enforced in production mode
- Observability is confirmed working end to end
- Demos and examples all run
- CI protects the golden path on every pull request

---

### Week 1: CI and Release Automation

**Goal:** Every pull request triggers lint and unit tests. Every tag triggers a full release pipeline.

**Tasks:**
- Create `.github/workflows/ci.yml` with triggers on `push` to `main` and `pull_request`:
  - `golangci-lint` on `src/`
  - `go test ./...` in `src/`
  - Docker build (smoke-build, no push)
  - Docker smoke test (start stack, call `/health`, stop stack)
- Create `.github/workflows/release.yml` with trigger on `push` tags `v*`:
  - Build multi-arch Docker image (`linux/amd64`, `linux/arm64`)
  - Push to Docker Hub: `neurondb/neuron-agent:VERSION` and `:latest`
  - Push to GHCR: `ghcr.io/neurondb/neuron-agent:VERSION` and `:latest`
  - Generate SBOM using Syft
  - Scan for vulnerabilities using Trivy
  - Sign image using cosign (if keys are available)
  - Create GitHub Release with auto-generated release notes
  - Attach binary artifacts for Linux amd64/arm64 and macOS arm64
- Add semantic versioning: `VERSION` file at root, `make release-tag VERSION=x.y.z`
- Fix `Makefile` `integration-test` target to remove `|| true` so failures surface

---

### Week 2: Kubernetes and Production Deployment

**Goal:** NeuronAgent can be deployed in a production Kubernetes cluster with minimal manual steps.

**Tasks:**
- Review existing `helm/` chart: `helm/Chart.yaml`, `helm/values.yaml`, `helm/templates/`
- Complete Helm chart with:
  - NeuronAgent Deployment with configurable replicas, resource requests/limits
  - Service (ClusterIP)
  - Ingress with TLS example
  - Secret for `AUTH_API_KEY_SECRET`, `DB_PASSWORD`
  - ConfigMap for non-secret config
  - Liveness probe (`/healthz`), Readiness probe (`/readyz`)
  - HorizontalPodAutoscaler example
  - PersistentVolume note for NeuronDB (external dependency)
  - PodDisruptionBudget
  - Prometheus ServiceMonitor (if prometheus-operator is available)
- Add `docs/deploy/kubernetes.md`:
  - Prerequisites
  - Install via Helm
  - NeuronDB as external dependency (instructions to configure connection)
  - Secrets management
  - Autoscaling configuration
  - Backup and restore in Kubernetes context
- Add `docs/deploy/docker-compose-production.md`:
  - Production `.env` checklist
  - Resource limits
  - Network isolation
  - Volume backup strategy
  - Reverse proxy example (nginx/Caddy)
  - TLS termination

---

### Week 3: Benchmarking and Scale Testing

**Goal:** Performance characteristics are measured, documented, and reproducible.

**Tasks:**
- Create `scripts/benchmark.sh` that:
  - Starts the local stack if not running
  - Measures: agent request latency (p50/p95/p99), memory search latency, RAG latency, tool execution latency, workflow start latency, concurrent request handling at 10/50/100 RPS
  - Uses `hey` or `wrk` for HTTP load testing, custom Go benchmark in `src/cmd/bench/` for NeuronDB-specific measurements
  - Outputs a structured JSON report to `docs/benchmarks-results/`
- Review and extend `src/cmd/bench/main.go` — it uses `NewRegistryWithNeuronDB`, which is the correct baseline
- Add `docs/benchmarks.md`:
  - Test environment specs (hardware, OS, Postgres version, NeuronDB version)
  - Methodology explanation
  - Results table (from actual local run — no invented numbers)
  - How to reproduce
  - Performance tuning guide: connection pool size, worker concurrency, memory batch size, embedding batch size

---

### Week 4: Final Launch Polish

**Goal:** Everything works. Every claim is true. Every command in the docs has been tested.

**Tasks:**
- Run complete install from scratch on a clean machine
- Execute every command in the README exactly as written
- Run `make demo`, verify output
- Run every example in `examples/`
- Run `make test`, verify all tests pass
- Run `make smoke`, verify golden path
- Run `make lint`, fix any issues
- Run release workflow dry run using `act` or manual steps
- Audit all docs links (use `markdown-link-check` or similar)
- Audit spelling (`codespell`)
- Audit Docker images: inspect sizes, confirm entrypoints are correct
- Audit security defaults: confirm production mode refuses weak config
- Create `docs/release-checklist.md` — operational checklist for each release
- Create `docs/v1-readiness-checklist.md` — gate criteria for calling the project v1.0
- Create `docs/launch-checklist.md` — public announcement checklist

---

## 11. Work Breakdown by Engineering Track

### Track A: Install and Packaging

**Goal:** Zero-friction local install. Production-grade Docker image.

**Current state:** Docker entrypoint is broken (binary name mismatch). Compose stack uses plain PostgreSQL (no NeuronDB). No root-level compose file. No install script. Makefile targets use deprecated `docker-compose` command.

**Target state:** `git clone` + `cp .env.example .env` + `docker compose up -d` + `make demo` works on any machine with Docker 24+.

**Monthly milestones:**
- M1: Fix Docker entrypoint. Create root docker-compose.yml with NeuronDB. Write install.sh. Write demo.sh.
- M2: Add production Docker Compose hardening. Document all env vars.
- M3: Add health probes to compose. Document backup strategy.
- M4: Verify all examples run from compose.
- M5: Record demo GIF from compose-based demo.
- M6: Publish multi-arch images to Docker Hub and GHCR.

**Key files:** `docker/Dockerfile`, `docker/docker-entrypoint.sh`, `docker/docker-compose.neuronsql.yml`, `Makefile`

**Tests needed:** Compose-based smoke test in CI

**Risks:** NeuronDB extension availability as a Docker image. Mitigation: use `init-scripts` approach to install extension in standard postgres image, or publish a `neurondb/postgres` image.

**Acceptance criteria:** Install time under five minutes. Demo succeeds on clean machine.

---

### Track B: Core Runtime

**Goal:** Agent lifecycle, session management, and configuration are predictable and tested.

**Current state:** Agent runtime exists and is functional. `SetUseStateMachine` is never called so state-machine run path is inactive. Config validation exists but startup errors can be opaque. `/version` endpoint is missing.

**Target state:** Every agent lifecycle event is validated, documented, and tested. Startup errors are human-readable. Version endpoint exists.

**Monthly milestones:**
- M1: Audit current runtime behavior. Document gaps.
- M2: Config validation hardening. Agent lifecycle tests. `/version` endpoint.
- M3: Permission enforcement in runtime. Workspace boundary tests.
- M4: Workflow agent step session fix.
- M5: State machine path evaluation — enable or remove.
- M6: Load test agent endpoint at 100 RPS.

**Key files:** `src/cmd/agent-server/main.go`, `src/internal/agent/runtime.go`, `src/internal/config/config.go`

**Tests needed:** Agent CRUD, session lifecycle, concurrent execution, error response shape

**Risks:** State machine path may have undiscovered bugs. Mitigation: start with legacy path, add state machine tests before enabling.

---

### Track C: Memory and RAG

**Goal:** Memory is durable, queryable, and well-documented. RAG pipeline is stable and has a working demo.

**Current state:** Hierarchical memory (STM/MTM/LPM) and episodic memory are implemented. RAG client wraps NeuronDB functions. `agent-server` does not register `rag` tool. `VectorClient.VectorSearch` has a SQL bug (query vector not bound).

**Target state:** Memory system has documented behavior for all paths. RAG has a stable ingest/query API and a working demo. Vector search SQL is fixed.

**Monthly milestones:**
- M1: Document current memory and RAG code paths.
- M2: Memory tests. Memory docs. Forget path API test.
- M3: Memory audit entries. Memory quality metadata API.
- M4: RAG ingest/query API stable. RAG demo working. Fix VectorClient SQL.
- M5: RAG example with fixture documents. Templates using RAG.
- M6: RAG latency benchmark.

**Key files:** `src/internal/agent/memory.go`, `src/internal/agent/hierarchical_memory.go`, `src/pkg/neurondb/rag_client.go`, `src/pkg/neurondb/vector_client.go`, `src/internal/tools/rag_tool.go`

**Tests needed:** Memory write/search/forget, RAG ingest/query with fixed fixtures, vector search correctness

**Risks:** NeuronDB extension must be running for vector and embedding tests. Mitigation: mark these tests with `requires_neurondb` and run in integration suite only.

---

### Track D: Tools and Workflows

**Goal:** Tool execution is safe, audited, and the full NeuronDB tool suite is available. Workflows are reliable with retry and scheduling.

**Current state:** Base registry (sql, http, code, shell, browser, visualization) is wired. NeuronDB tools (rag, vector, ml, analytics, hybrid_search, reranking) are only in `bench` binary. `RetrievalTool` is never registered. Workflow schedule runner does not exist. Workflow retries are incomplete. `conditional` step type is only in `advanced_engine.go` (not wired).

**Target state:** All tool handlers available via config. Workflow retry, scheduling, and conditional steps work. Approval gates are tested.

**Monthly milestones:**
- M1: Audit tool registry gaps. Document tool system.
- M2: Fix tool registry to include NeuronDB tools. Fix RetrievalTool registration. Add tool audit logs.
- M3: Tool permission enforcement tests.
- M4: Workflow retry scheduler. Schedule runner. Conditional step. Approval gate demo.
- M5: Tool templates. Tool examples.
- M6: Tool execution latency benchmarks.

**Key files:** `src/internal/tools/registry.go`, `src/cmd/agent-server/main.go`, `src/internal/workflow/engine.go`, `src/internal/workflow/advanced_engine.go`

**Tests needed:** Tool permission denied, tool timeout, tool audit log entry, workflow retry, workflow scheduling, conditional step branching

**Risks:** Workflow schedule runner requires a reliable ticker — must be restartable and idempotent. MCP sync is a stub (low priority for M1-M3).

---

### Track E: Security and Enterprise Controls

**Goal:** NeuronAgent passes an enterprise security evaluation.

**Current state:** Auth middleware is in place (API keys, org scoping, admin role check). Role model exists in code. Permission checks are partial. Audit log tables likely exist in schema but coverage is incomplete. JWT is mentioned in config struct but implementation depth is unclear.

**Target state:** All endpoints are protected. Role-based permission checks are enforced for every resource type. Audit log covers all defined event types. Production mode refuses weak defaults.

**Monthly milestones:**
- M1: Security gap audit. Document what is enforced vs what is missing.
- M2: Production config validation. Startup security checks.
- M3: Full RBAC implementation. Audit log coverage. Permission tests. `docs/security.md`.
- M4: OpenClaw bridge uses `service` role.
- M5: `SECURITY.md`. Vulnerability reporting process.
- M6: Trivy scan in release pipeline. Image signing with cosign.

**Key files:** `src/internal/auth/`, `src/internal/api/middleware.go`, `src/internal/api/org_middleware.go`, `sql/neuron-agent.sql`

**Tests needed:** Permission denied for every role × resource combination, audit log entry for every event type, production config refuses weak secrets

---

### Track F: Observability and Operations

**Goal:** NeuronAgent in production can be monitored, debugged, and operated without tribal knowledge.

**Current state:** Prometheus metrics endpoint exists. Request ID middleware exists. zerolog structured logging is in use. OpenTelemetry is a dependency. Graceful shutdown is in main. `MemoryPromoter` background worker exists but restart behavior on panic is not confirmed.

**Target state:** All defined Prometheus metrics are instrumented. `/readyz` is distinct from `/healthz`. Grafana starter dashboard exists. Backup and restore docs are complete.

**Monthly milestones:**
- M1: Audit which metrics are instrumented vs defined.
- M2: Confirm graceful shutdown. Add worker panic recovery.
- M3: Full metrics instrumentation. `/readyz` endpoint. `docs/observability.md`. Sample Grafana dashboard.
- M4: Queue depth metrics. Workflow execution metrics.
- M5: Tracing span examples using OpenTelemetry.
- M6: End-to-end observability demo in production compose stack.

**Key files:** `src/cmd/agent-server/main.go`, `src/internal/api/middleware.go`, `src/internal/worker/`

**Tests needed:** `/healthz` and `/readyz` response shape, metric counter increments on known actions

---

### Track G: Docs and Developer Experience

**Goal:** Docs are accurate, well-structured, and never claim features that don't work.

**Current state:** Docs exist (`docs/`) but mix `.md` and `.txt`. No `docs/README.md` index. Hosted docs at neurondb.ai are canonical narrative but local docs are more audit/reference oriented. README has version drift and claims that need verification.

**Target state:** 20+ doc pages in `docs/`. No `.txt` files in docs tree (convert or remove). Every doc is accurate and linked from an index.

**Monthly milestones:**
- M1: Audit docs. Create docs/roadmap/. Fix README version drift.
- M2: `docs/configuration.md`, `docs/memory.md`, `docs/tools.md`.
- M3: `docs/security.md`, `docs/observability.md`, `docs/reliability.md`.
- M4: `docs/rag.md`, `docs/workflows.md`, `docs/integrations/openclaw.md`.
- M5: `docs/api.md`, `docs/use-cases.md`, `docs/comparisons.md`, `docs/faq.md`, `docs/index.md`, `CONTRIBUTING.md`.
- M6: Full link audit. Spelling check. Docs completeness checklist.

**Key files:** `docs/`, `README.md`, `CHANGELOG.md`, `CONTRIBUTING.md`

**Tests needed:** Markdown link checker in CI. Spelling check.

---

### Track H: Examples and Adoption

**Goal:** Every use case has a working, runnable example.

**Current state:** `examples/` has Python and TypeScript NeuronSQL minimal examples. `src/examples/` has a Python client library and modular examples. CLI examples exist. No examples cover the agent quickstart, RAG demo, or DBA agent.

**Target state:** Five working examples (quickstart, RAG, DBA agent, workflow approval, OpenClaw) each with a README and a runnable script.

**Monthly milestones:**
- M1: `examples/quickstart-chat/`. `examples/rag-with-docs/`. `examples/sql-agent/`.
- M2: Add expected output to each example.
- M3: Add cleanup scripts to each example.
- M4: `examples/postgres-dba-agent/`. `examples/workflow-approval/`. `examples/openclaw-bridge/`.
- M5: Five agent templates in `templates/`. Demo GIF. `docs/use-cases.md`.
- M6: Verify all examples on clean machine.

---

### Track I: Integrations

**Goal:** OpenClaw bridge works. SDK clients work. OpenAPI spec is accurate.

**Current state:** Claw gateway routes exist in `buildRouter`. OpenAPI spec exists at `src/openapi/openapi.yaml`. Python client exists in `src/examples/`. TypeScript SDK exists in `src/sdks/typescript/`. MCP sync is a stub.

**Target state:** Claw bridge has auth, permissions, docs, and tests. OpenAPI spec matches live API. Python and TypeScript clients are packaged and tested.

**Monthly milestones:**
- M1: Document claw bridge current state.
- M2: Confirm claw bridge auth.
- M3: Claw bridge uses `service` role.
- M4: All claw bridge tools exposed. Auth confirmed. Docs written. Example added.
- M5: OpenAPI spec validated against live router. Clients packaged.
- M6: Claw bridge integration test in CI.

---

### Track J: Release and Community

**Goal:** Regular, automated releases. A healthy GitHub presence with clear contribution pathways.

**Current state:** CI is `workflow_dispatch` only. CHANGELOG has no versioned releases. Release checklist is manual. No Docker Hub publish automation. No image signing. No SBOM.

**Target state:** Every PR triggers CI. Every tag triggers a full automated release pipeline. Community health files are complete.

**Monthly milestones:**
- M1: Add push/PR trigger to CI workflow.
- M2: Fix `integration-test` Makefile target.
- M3: Add PR template and issue templates.
- M4: `CONTRIBUTING.md`. `CODE_OF_CONDUCT.md`. `SECURITY.md`.
- M5: Complete CHANGELOG. Update ROADMAP.md.
- M6: Full CI/CD pipeline. Docker Hub and GHCR publish. Release automation. Good-first-issues labeled.

---

## 12. Measurable Success Metrics

These are pass/fail gates for calling the project ready for public adoption:

| Metric | Target | How to measure |
|---|---|---|
| Install time | Under 5 minutes from `git clone` | Time from clean machine with Docker only |
| Quickstart commands | 4 or fewer | Count commands in README quickstart |
| Demo success rate | 100% on clean machines | Run on 3 different machines before release |
| Golden path smoke test | Passes on every PR | CI check |
| API docs accessible locally | Yes | `curl http://localhost:8080/docs` returns 200 |
| Working examples | 5 or more | Count in `examples/` that have passing demo scripts |
| Agent templates | 5 or more | Count in `templates/` |
| Docker images published | Yes | Check Docker Hub and GHCR |
| Security docs complete | Yes | `docs/security.md` exists and is accurate |
| Production deployment docs | Yes | `docs/deploy/` has docker-compose and kubernetes docs |
| OpenClaw bridge working | Yes | All 3 claw endpoints return correct responses |
| Benchmark script present | Yes | `scripts/benchmark.sh` produces a report |
| README structure | Per target structure | Manual review |
| No fake feature claims | Zero | Manual audit |
| Release checklist complete | Yes | `docs/release-checklist.md` and it is followed |
| CI automation | Push + PR triggers | Check `.github/workflows/` |
| Go test coverage (core) | Above 60% | `go test -cover ./internal/...` |

---

## 13. README Target Structure

The README must answer three questions in 90 seconds: what is it, how do I run it, and why do I want it.

```
# NeuronAgent

[One-sentence description]

[Badges: CI status, Go version, License — only for real workflows]

## Install

[4-5 commands. Docker only. No manual steps.]

## First demo

[GIF or terminal output of the demo running]

## What you get

[6-8 bullet points of concrete capabilities]

## Why NeuronAgent

[3-4 sentences on positioning vs alternatives]

## Use cases

[4-6 one-line scenarios with links to examples]

## Examples

[Links to examples/ directory]

## OpenClaw bridge

[One paragraph, link to docs/integrations/openclaw.md]

## Production deployment

[Link to docs/deploy/]

## Docs

[Link to docs/ index]

## Contributing

[Link to CONTRIBUTING.md]

## Security

[Link to SECURITY.md]

## License

[License statement]
```

**Rules:**
- No architecture diagrams in README
- No internal package names
- No fake badges
- No features not backed by working code
- Every command tested before merge

---

## 14. Documentation Target Structure

```
docs/
├── index.md                    # Product intro, quick navigation
├── quickstart.md               # Install + first demo, step by step
├── concepts.md                 # Core concepts: agents, memory, RAG, tools, workflows
├── architecture.md             # Architecture diagrams, component map
├── configuration.md            # All environment variables with types and defaults
├── api.md                      # API reference with curl examples
├── memory.md                   # Memory system: tiers, write, search, forget
├── rag.md                      # RAG: ingest, chunk, search, answer, citations
├── tools.md                    # Tool system: registry, permissions, handler types
├── workflows.md                # Workflow engine: steps, retries, approval gates
├── security.md                 # Auth, RBAC, key management, audit logs
├── observability.md            # Metrics, logs, tracing, health endpoints
├── backup-restore.md           # Backup strategy, restore procedure, PITR
├── reliability.md              # Graceful shutdown, worker lifecycle, idempotency
├── benchmarks.md               # Performance measurements and tuning guide
├── use-cases.md                # Six concrete use case descriptions
├── comparisons.md              # Honest comparison with alternatives
├── faq.md                      # 20 common questions
├── examples.md                 # Index of all examples
├── development.md              # Contributor guide, architecture, adding features
├── templates.md                # Agent template reference
├── deploy/
│   ├── docker-compose.md       # Local and production Compose setup
│   └── kubernetes.md           # Helm chart, Kubernetes deployment guide
└── integrations/
    └── openclaw.md             # OpenClaw bridge setup and tool reference
```

---

## 15. Final Summary

### What NeuronAgent is today

NeuronAgent is a sophisticated Go runtime (module `github.com/neurondb/NeuronAgent`, Go 1.24) with:
- A rich HTTP/WebSocket API (50+ endpoints across agents, sessions, memory, RAG, tools, workflows)
- A multi-tier memory system (STM/MTM/LPM + episodic) backed by PostgreSQL
- A NeuronDB client wrapper for embeddings, RAG, hybrid search, ML, and analytics
- A workflow engine supporting agent, tool, approval, HTTP, and SQL steps
- Auth middleware with API keys, org scoping, and admin roles
- Prometheus metrics, structured logging, and OpenTelemetry wiring
- An existing Helm chart and Terraform configuration
- 86 Go test files and 150+ Python test files

### What NeuronAgent should become

NeuronAgent should become the reference implementation of a database-native agent runtime: the standard answer to "how do I give my AI agents durable memory, safe tool execution, reliable workflows, and enterprise-grade audit trails, all backed by PostgreSQL." It should be:
- Installable in five minutes
- Demonstrable in ten minutes
- Trustworthy enough for enterprise evaluation
- Clear enough for contributors to get productive in one day
- Impressive enough for GitHub visitors to star and return to

### Top ten gaps

1. **Docker entrypoint is broken** — binary built as `neuron-agent`, entrypoint executes `neuronagent`. Container exits immediately.
2. **Compose stack has no NeuronDB** — uses plain `postgres:17-alpine`; schema migrations require `CREATE EXTENSION neurondb`.
3. **NeuronDB tool handlers not wired in agent-server** — rag, vector, ml, analytics, hybrid_search, reranking are absent from the live registry.
4. **RetrievalTool is never registered** — agentic retrieval is broken by default.
5. **VectorClient.VectorSearch has broken SQL** — query vector is not bound; vector tool path is non-functional.
6. **No CI automation** — CI is `workflow_dispatch` only; no push or PR triggers.
7. **No one-command install** — no `install.sh`, no root-level `docker-compose.yml`, no `demo.sh`.
8. **Workflow schedule runner missing** — cron tables exist, no runner.
9. **`integration-test` silently passes on failure** — `|| true` swallows test errors.
10. **Migration naming inconsistency** — bulk psql loop vs versioned Go loader with incompatible naming conventions.

### Top ten Month 1 tasks

1. Fix Docker entrypoint binary name mismatch
2. Create root-level `docker-compose.yml` with NeuronDB-enabled PostgreSQL
3. Create `.env.example` with all required variables
4. Create `scripts/install.sh` with prerequisite checks and clear error messages
5. Create `scripts/demo.sh` that covers the golden path
6. Add `make up`, `make down`, `make demo`, `make smoke` to Makefile
7. Rewrite README with correct structure and no false claims
8. Create `examples/quickstart-chat/` with working script
9. Add smoke test to CI with push/PR trigger
10. Create `docs/roadmap/current-state-audit.md`, `gap-analysis.md`, `feature-truth-table.md`

### Risks

| Risk | Severity | Mitigation |
|---|---|---|
| NeuronDB extension is not publicly available as a Docker image | High | Provide init-script approach or publish `neurondb/postgres` test image |
| LLM API key required for full demo | Medium | Use mock/stub responses for smoke tests; document real key requirement clearly |
| State machine agent path may have undiscovered bugs | Medium | Keep legacy path as default; enable state machine only with explicit config flag |
| Python test suite may not run against live Go server | Medium | Separate unit tests (no server) from integration tests (require server) |
| Workflow scheduling requires persistent worker — may conflict with stateless deployment | Medium | Document that schedule runner requires single-writer mode or distributed lock |
| `SyncFromMCP` stub will block MCP integration forever | Low | Either implement or remove the stub; do not leave dead code |

### Recommended first PR

**Fix Docker container startup.**
- Fix binary name in `docker/Dockerfile` to match `docker/docker-entrypoint.sh`
- Add basic health check to Dockerfile
- Add root-level `docker-compose.yml` with correct service names
- Add `.env.example`
- This is the prerequisite for everything else.

### Recommended second PR

**Create one-command local install.**
- Create `scripts/install.sh`
- Create `scripts/demo.sh`
- Add `make up`, `make down`, `make demo`, `make smoke` targets
- Create `examples/quickstart-chat/`
- Add smoke test GitHub Actions workflow with push trigger

### Recommended third PR

**Rewrite README for adoption.**
- Restructure README to target structure
- Remove false claims and unverified badges
- Add confirmed working quickstart commands
- Link to new examples
- Move architecture diagrams to `docs/architecture.md`
- Fix Go version badge from 1.23 to 1.24

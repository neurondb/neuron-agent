# NeuronAgent Current State Audit

> Last updated: 2026-04-30  
> Audit basis: Code inspection of repository at `/mnt/pge/ndb/neuron-agent`  
> Auditor: Automated code analysis  
> Status: Ground truth — do not update without re-verifying code

This document records what NeuronAgent actually does today, verified by reading code. It is the source of truth for the gap analysis and feature truth table.

---

## Repository Structure

```
neuron-agent/
├── agent-server/           # Compiled binary output (not source)
├── bin/                    # Build artifacts, scripts mirror, SQL mirror
├── docker/                 # Dockerfile, entrypoint, compose variants
├── docs/                   # Documentation (mix of .md and .txt)
├── examples/               # Python and TypeScript NeuronSQL examples only
├── helm/                   # Helm chart (Chart.yaml, values.yaml, templates/)
├── scripts/                # Operational shell scripts
├── sql/                    # SQL migration files
├── src/                    # Go module root (go.mod here, not at repo root)
│   ├── cmd/
│   │   ├── agent-server/   # Main HTTP server binary
│   │   ├── bench/          # Benchmark binary
│   │   └── generate-key/   # API key generation utility
│   ├── cli/                # Cobra CLI binary
│   ├── examples/           # Python client library and modular examples
│   ├── internal/
│   │   ├── agent/          # Agent runtime, memory, planning, reflection
│   │   ├── api/            # HTTP handlers, middleware
│   │   ├── auth/           # API keys, rate limiting, principal management
│   │   ├── config/         # Config structs, loading, validation, defaults
│   │   ├── core/           # App initialization, DB wiring
│   │   ├── db/             # Database queries, models, migrations
│   │   ├── modules/        # Module system (NeuronSQL module)
│   │   ├── neuronsql/      # NeuronSQL handler implementations
│   │   ├── tools/          # Tool handlers and registry
│   │   ├── worker/         # Background workers
│   │   └── workflow/       # Workflow engine
│   ├── openapi/            # OpenAPI specification
│   ├── pkg/
│   │   ├── llm_sql/        # LLM SQL sidecar client
│   │   ├── module/         # Module interface
│   │   ├── neuronsql/      # NeuronSQL package
│   │   └── neurondb/       # NeuronDB Go client wrappers
│   ├── sdks/
│   │   └── typescript/     # TypeScript SDK
│   ├── services/           # Additional service layer
│   └── tests/              # Integration and Python test suite
├── terraform/              # Terraform infrastructure
└── tests/                  # Root-level test fixtures
    └── fixtures/
```

**Go module:** `github.com/neurondb/NeuronAgent` (declared in `src/go.mod`)  
**Go version:** 1.24.0  
**Note:** There is no `go.mod` at the repository root. All Go commands must run from `src/`.

---

## HTTP Server

### Status: Working

The server starts via `src/cmd/agent-server/main.go`. It:
1. Loads config (env vars or YAML file if `CONFIG_PATH` is set)
2. Builds `core.NewApp(cfg)` — opens DB connection
3. Initializes NeuronDB client via `neurondb.NewClient(app.DB().DB)`
4. Registers NeuronSQL module
5. Constructs tool registry
6. Builds agent runtime
7. Initializes workflow engine
8. Calls `buildRouter` to wire all HTTP routes
9. Starts HTTP server with graceful shutdown on SIGINT/SIGTERM

**HTTP framework:** `gorilla/mux`

**Middleware stack (applied in order):**
1. `RequestIDMiddleware` — adds `X-Request-ID` to every request
2. `RequestTimeoutMiddleware(60s)` — aborts requests that take over 60 seconds
3. `AuthMiddleware` — validates `Authorization: Bearer <key>` or `Authorization: ApiKey <key>`
4. `OrgMiddleware` — resolves org ID from API key
5. `RejectUnknownFieldsMiddleware` — rejects JSON payloads with unknown fields (if enabled in config)
6. `RequestBodyLimitMiddleware(10MiB)` — limits request body size
7. `SecurityHeadersMiddleware` — adds `X-Content-Type-Options`, `X-Frame-Options`, etc.
8. `CORSMiddleware` — CORS with configurable allowed origins
9. `LoggingMiddleware` — structured request logging via zerolog

**Endpoints not behind auth:** `/health`, `/metrics` only.

**Known gap:** `/ws` WebSocket endpoint is behind auth middleware (requires `Authorization` header), but handler-level code also checks for `api_key` query parameter. The query parameter path is effectively dead unless `/ws` is exempted from `AuthMiddleware`.

---

## API Surface

### Status: Working — all routes are registered in `buildRouter`

#### Agent management
| Method | Path | Status |
|---|---|---|
| GET | `/api/v1/agents` | Working |
| POST | `/api/v1/agents` | Working |
| GET | `/api/v1/agents/{id}` | Working |
| PUT | `/api/v1/agents/{id}` | Working |
| DELETE | `/api/v1/agents/{id}` | Working |
| POST | `/api/v1/agents/{id}/clone` | Working |
| POST | `/api/v1/agents/{id}/plan` | Working |
| POST | `/api/v1/agents/{id}/reflect` | Working |
| POST | `/api/v1/agents/{id}/delegate` | Working |
| GET | `/api/v1/agents/{id}/metrics` | Working |
| GET | `/api/v1/agents/{id}/costs` | Working |
| GET, PUT | `/api/v1/agents/{id}/budget` | Working |

#### Session and message management
| Method | Path | Status |
|---|---|---|
| GET, POST | `/api/v1/sessions` | Working |
| GET, PUT, DELETE | `/api/v1/sessions/{id}` | Working |
| GET, POST | `/api/v1/sessions/{id}/messages` | Working (POST triggers agent execution) |
| GET, PUT, DELETE | `/api/v1/sessions/{id}/messages/{msgId}` | Working |
| POST | `/api/v1/sessions/{id}/reflect` | Working |
| POST | `/api/v1/sessions/{id}/feedback` | Working |

#### Memory management
| Method | Path | Status |
|---|---|---|
| GET | `/api/v1/agents/{id}/memory` | Working |
| POST | `/api/v1/agents/{id}/memory/search` | Working |
| POST | `/api/v1/agents/{id}/memory/summarize` | Working |
| POST | `/api/v1/agents/{id}/memory/consolidate` | Working |
| POST | `/api/v1/agents/{id}/memory/corruption` | Working |
| POST | `/api/v1/agents/{id}/memory/forget` | Working |
| POST | `/api/v1/agents/{id}/memory/conflicts` | Working |
| POST | `/api/v1/agents/{id}/memory/feedback` | Working |
| GET | `/api/v1/agents/{id}/memory/quality` | Working |
| GET, DELETE | `/api/v1/memory/{chunkId}` | Working |

#### Run management
| Method | Path | Status |
|---|---|---|
| POST | `/api/v1/agents/{id}/runs` | Working |
| GET | `/api/v1/runs/{id}` | Working |
| POST | `/api/v1/runs/{id}/cancel` | Working |
| GET | `/api/v1/runs/{id}/plan` | Working |
| GET | `/api/v1/runs/{id}/steps` | Working |
| GET | `/api/v1/runs/{id}/traces` | Working |
| GET | `/api/v1/runs/{id}/explain/tool/{inv_id}` | Working |
| GET | `/api/v1/runs/{id}/explain/memory` | Working |
| GET | `/api/v1/runs/{id}/explain/model` | Working |
| GET | `/api/v1/runs/{id}/explain/plan` | Working |

#### Tool management
| Method | Path | Status |
|---|---|---|
| GET, POST | `/api/v1/tools` | Working |
| GET, PUT, DELETE | `/api/v1/tools/{id}` | Working |
| GET | `/api/v1/tools/{id}/analytics` | Working |

#### Workflow management
| Method | Path | Status |
|---|---|---|
| GET, POST | `/api/v1/workflows` | Working |
| GET, PUT, DELETE | `/api/v1/workflows/{id}` | Working |
| POST, GET | `/api/v1/workflows/{workflow_id}/steps` | Working |
| POST | `/api/v1/workflows/{workflow_id}/execute` | Working |
| GET | `/api/v1/workflows/{workflow_id}/executions` | Working |
| GET | `/api/v1/workflows/executions/{execution_id}` | Working |
| POST, GET, PUT, DELETE | `/api/v1/workflows/{workflow_id}/schedule` | Working (CRUD only — no runner) |
| GET | `/api/v1/workflows/schedules` | Working (CRUD only — no runner) |

#### RAG
| Method | Path | Status |
|---|---|---|
| POST | `/api/v1/rag/ingest` | Conditional — only registered when `ragClient != nil` |

**Gap:** No `GET /api/v1/rag/query` endpoint. RAG query may go through agent execution, not a direct API call.

#### NeuronSQL (module-mounted at `/api/v1/neuronsql/`)
| Method | Path | Status |
|---|---|---|
| POST | `/api/v1/neuronsql/generate` | Working |
| POST | `/api/v1/neuronsql/optimize` | Working |
| POST | `/api/v1/neuronsql/validate` | Working |
| POST | `/api/v1/neuronsql/plpgsql` | Working |

#### Approvals and feedback
| Method | Path | Status |
|---|---|---|
| GET | `/api/v1/approvals` | Working |
| GET | `/api/v1/approvals/{id}` | Working |
| POST | `/api/v1/approvals/{id}/approve` | Working |
| POST | `/api/v1/approvals/{id}/reject` | Working |
| GET | `/api/v1/feedback` | Working |
| GET | `/api/v1/feedback/stats` | Working |

#### Analytics and governance
| Method | Path | Status |
|---|---|---|
| GET | `/api/v1/analytics/overview` | Working |
| GET | `/api/v1/analytics/retrieval-stats` | Working |
| GET | `/api/v1/governance/costs` | Working |
| GET | `/api/v1/governance/tool-risk` | Working |
| GET | `/api/v1/governance/policy-blocks` | Working |
| GET | `/api/v1/governance/memory-growth` | Working |
| GET | `/api/v1/governance/agent-performance` | Working |

#### Admin
| Method | Path | Status |
|---|---|---|
| GET | `/api/v1/admin/config` | Working (admin role required) |
| GET | `/api/v1/admin/diagnostics` | Working (admin role required) |

#### Batch operations
| Method | Path | Status |
|---|---|---|
| POST | `/api/v1/batch/agents` | Working |
| POST | `/api/v1/batch/agents/delete` | Working |
| POST | `/api/v1/batch/messages/delete` | Working |
| POST | `/api/v1/batch/tools/delete` | Working |

#### OpenClaw gateway
| Method | Path | Status |
|---|---|---|
| GET | `/claw/v1/health` | Working |
| POST | `/claw/v1/tools/list` | Working |
| POST | `/claw/v1/tools/run` | Working |

#### Blueprint
| Method | Path | Status |
|---|---|---|
| GET | `/api/v1/blueprints` | Working |
| POST | `/api/v1/agents/from-blueprint` | Working |

#### Missing from live router
- `/health` returns health status
- `/metrics` returns Prometheus metrics
- `/version` — **does not exist** (planned)
- `/healthz` — **does not exist** (Kubernetes liveness alias)
- `/readyz` — **does not exist** (Kubernetes readiness, needs DB check)
- `/docs` or `/swagger` — **does not exist** (OpenAPI UI not exposed at runtime)

---

## Authentication

### Status: Working for API keys; JWT implementation depth unclear

**API key auth:** `AuthMiddleware` in `src/internal/api/middleware.go` requires `Authorization: Bearer <key>` or `Authorization: ApiKey <key>` on all routes except `/health` and `/metrics`. Keys are validated via `auth.APIKeyManager`.

**Rate limiting:** `auth.RateLimiter` is applied per API key in `AuthMiddleware`.

**Org scoping:** `OrgMiddleware` attaches org ID from API key record.

**Admin role:** `RequireRole(..., RoleAdmin)` is checked in admin handlers.

**JWT:** `AUTH_JWT_SECRET` exists in config struct (`src/internal/config/config.go`). Implementation depth in `auth/` package is not confirmed — needs verification before claiming JWT support.

**Workspace scoping:** API keys can be scoped to workspaces — exists in the auth model.

**Role model in code:** `RoleAdmin`, `RoleDeveloper`, `RoleViewer` exist. `RoleService` and `RoleOwner` need verification.

---

## Agent Runtime

### Status: Working — legacy execution path is active; state machine path is inactive

**Runtime struct:** `agent.Runtime` in `src/internal/agent/runtime.go`

**Construction:** `NewRuntime(db, queries, tools, embedClient, ragClient, hybridClient)` wires memory managers, planner, reflector, LLM client, NeuronDB clients, and tool permission checker.

**Execution path:** `Execute(ctx, sessionID, userMessage)` runs through:
1. State machine check — **`SetUseStateMachine` is never called from `agent-server/main.go`**, so this path is inactive
2. Distributed coordinator check — only if `distributedCoordinator` is set (not set by default)
3. Async path — only if `asyncExecutor` says so
4. **Legacy pipeline** — the active path: event stream → planning → LLM → tools → memory/reflection

**Agent states:** Active, paused, archived — defined in DB schema.

**Known gaps:**
- `SetUseStateMachine` never called → state machine run persistence path is off
- Planning and reflection are in code but interaction with LLM depends on LLM client configuration

---

## Memory System

### Status: Implemented — multiple tiers, background worker, vector search

**Memory tiers:**

| Tier | Table | Scope | Expiry |
|---|---|---|---|
| Short-term (STM) | `memory_stm` | Session-scoped | TTL-based |
| Medium-term (MTM) | `memory_mtm` | Agent-scoped | Days |
| Long-term (LPM) | `memory_lpm` | Agent-scoped | Weeks+ |
| General chunks | `memory_chunks` | Agent-scoped | Manual/TTL |
| Episodic | Separate table | Event sequences | Configurable |

**Vector search:** Uses `embedding <=> $1::neurondb_vector` cosine distance. Requires NeuronDB extension.

**Memory promotion:** `MemoryPromoter` worker runs STM→MTM→LPM promotion based on access frequency and importance. Worker exists in `src/internal/worker/memory_promoter.go`.

**Summarization:** Current implementation uses **truncation**, not an LLM. Comment in code notes this is a placeholder for LLM-based summarization.

**Agentic retrieval:** Conditional path in `ContextLoader`. When `agentic_retrieval_enabled` is in agent config, uses `RetrievalAdapter` → hierarchical first, then flat `memory_chunks`.

**Known gaps:**
- `RetrievalTool` (`src/internal/tools/retrieval_tool.go`) is never registered on the tool registry — agentic retrieval prompts reference a tool that does not exist
- Summarization is truncation-based, not LLM-based
- Memory conflict detection exists as an API endpoint but behavior is not documented

---

## RAG System

### Status: Partially wired — client exists, HTTP ingest endpoint is conditional

**NeuronDB RAG client:** `src/pkg/neurondb/rag_client.go` — wraps `neurondb_chunk_text`, `neurondb_rerank_results`, `neurondb_generate_answer`, `rag_query`, `rag_ingest_document`.

**RAG HTTP endpoint:** `POST /api/v1/rag/ingest` is registered only when `ragClient != nil`. The `ragClient` is initialized as `neurondb.NewClient(app.DB().DB).RAG` only when the DB is connected. So in practice, when NeuronDB extension is present, this endpoint is live.

**RAG tool:** `src/internal/tools/rag_tool.go` exists but is only registered in `NewRegistryWithNeuronDB` — **not in the `agent-server` default registry** (`NewRegistry`).

**Advanced RAG:** `src/internal/agent/advanced_rag.go` and `src/internal/agent/modular_rag.go` exist as experimental paths but are not confirmed to be in the hot path.

**Known gaps:**
- No `GET /api/v1/rag/query` HTTP endpoint
- RAG tool not wired in default `agent-server` registry
- No citation/source reference in API response documented

---

## Tool System

### Status: Base tools working; NeuronDB tools missing from agent-server

**Registry:** `src/internal/tools/registry.go` — `ExecuteTool` validates args, runs optional RBAC, applies circuit breaker and timeout, writes audit log.

#### Tools registered in default `agent-server` (via `NewRegistry`)

| Handler type | File | Status |
|---|---|---|
| `sql` | `src/internal/tools/sql_tool.go` | Working — but defaults allow writes (needs read-only default) |
| `http` | `src/internal/tools/http_tool.go` | Working |
| `code` | `src/internal/tools/code_tool.go` | Working — reads/analyzes files in allowed dirs |
| `shell` | `src/internal/tools/shell_tool.go` | Working — exists with constraints |
| `browser` | `src/internal/tools/browser_tool.go` | Working — requires Chromedp |
| `visualization` | `src/internal/tools/visualization_tool.go` | Working |

#### NeuronSQL module tools (registered via module system)

| Handler type | Status |
|---|---|
| `neuronsql.schema_snapshot` | Working |
| `neuronsql.validate_sql` | Working |
| `neuronsql.explain_json` | Working |
| `neuronsql.generate_select` | Working |
| `neuronsql.optimize_select` | Working |
| `neuronsql.plpgsql_generate` | Working |
| `neuronsql_generate` (legacy) | Working |
| `neuronsql_optimize` (legacy) | Working |

#### Tools NOT registered in agent-server (only in bench or unused)

| Handler type | File | Status |
|---|---|---|
| `rag` | `src/internal/tools/rag_tool.go` | Missing from agent-server registry |
| `vector` | `src/internal/tools/vector_tool.go` | Missing from agent-server registry + VectorSearch SQL bug |
| `ml` | `src/internal/tools/ml_tool.go` | Missing from agent-server registry |
| `analytics` | `src/internal/tools/analytics_tool.go` | Missing from agent-server registry |
| `hybrid_search` | `src/internal/tools/hybrid_search_tool.go` | Missing from agent-server registry |
| `reranking` | `src/internal/tools/reranking_tool.go` | Missing from agent-server registry |
| `filesystem` | `src/internal/tools/filesystem_tool.go` | Only in `NewRegistryWithVFS` — not wired |
| `memory` | `src/internal/tools/memory_tool.go` | Only in `NewRegistryWithAllFeatures` — not wired |
| `collaboration` | `src/internal/tools/collaboration_tool.go` | Only in `NewRegistryWithAllFeatures` — not wired |
| `multimodal` | `src/internal/tools/multimodal_tool.go` | Only in `NewRegistryWithAllFeatures` — not wired |
| `web_search` | `src/internal/tools/web_search_tool.go` | Only in `NewRegistryWithAllFeatures` — not wired |
| `retrieval` | `src/internal/tools/retrieval_tool.go` | Never registered anywhere |

#### Dead code
- `RegisterNeuronSQLTools` in `src/internal/neuronsql/tools/register.go` — no callers; the table/index NeuronSQL tools are unreachable
- `SyncFromMCP` in registry — stub that returns nil
- `VectorClient.VectorSearch` SQL — query vector is not bound, makes the tool non-functional even if registered

---

## Workflow Engine

### Status: Basic execution working; scheduling, retries, and conditional steps are incomplete

**Engine:** `src/internal/workflow/engine.go`

**Supported step types (in main engine switch):** `agent`, `tool`, `approval`, `http`, `sql`

**Not supported in main engine:** `conditional` — exists in `advanced_engine.go` and TypeScript SDK types but not in the main engine `switch`.

**DAG execution:** Topological ordering of steps is implemented. Steps execute in dependency order.

**Idempotency:** `idempotency_key` field exists on steps. Short-circuit on duplicate keys is implemented.

**Retries:** On failure, `retry_count` is incremented. Comment in code explicitly states "retry scheduling is not implemented" — retries are tracked but not automatically re-executed.

**Compensation:** `CompensateStep` exists and can run a linked compensation step. Depth of integration is unclear.

**Schedule runner:** Schedule CRUD endpoints exist. `ListWorkflowSchedulesByNextRun` exists in DB layer. **No background worker calls this function.** Schedules are stored but never executed automatically.

**Advanced engine:** `src/internal/workflow/advanced_engine.go` is **not referenced** from HTTP handlers or the main engine. It is dead code.

**Agent step bug:** `executeAgentStep` uses `GetSession(uuid.Nil)` then creates a session — using a nil UUID for session lookup is likely to cause issues in production workflows.

---

## Configuration System

### Status: Working — env-first loading with YAML file support

**Config loading order:**
1. `DefaultConfig()` — safe defaults
2. `LoadFromEnv` — reads `SERVER_*`, `DB_*`, `AUTH_*`, `LOG_*`, etc.
3. `ApplyProfile` — applies environment profile adjustments
4. `ValidateConfig` — basic validation

**Config path:** If `CONFIG_PATH` env is set, loads YAML then applies env overrides.

**Key env vars:**
- `SERVER_HOST`, `SERVER_PORT` — default `0.0.0.0:8080`
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- `AUTH_API_KEY_SECRET`
- `CORS_ALLOWED_ORIGINS`
- `LOG_LEVEL`, `LOG_FORMAT`
- `LLM_SQL_BASE_URL`, `LLM_SQL_API_KEY` — for NeuronSQL sidecar

**Modules config:** `Modules["neuronsql"].Enabled: true` by default.

**Gap:** No `/version` endpoint. No startup banner with config summary. Production mode does not refuse weak secrets at startup.

---

## Database and Migrations

### Status: Schema exists; migration approach is inconsistent

**Schema:** `sql/neuron-agent.sql` — creates `neurondb_agent` schema, tables, indexes. Requires `neurondb`, `pgcrypto`, `uuid-ossp` extensions.

**Additional schema files:** `sql/neuronagent_*_schema.sql` files for partitioning, RLS, and advanced features.

**Shell migration script:** `scripts/neuronagent-migrate.sh` — runs all `*.sql` files in `sql/` via `psql`. Order is filesystem glob order, not explicit.

**Go migration system:** `src/internal/db/schema.go` expects files named `001_name.sql`, `002_name.sql`, etc. The main schema file `neuron-agent.sql` does not match this pattern.

**Gap:** Two migration approaches that do not align. Migration version tracking is inconsistent.

---

## Docker and Containers

### Status: BROKEN — container exits immediately due to binary name mismatch

**Dockerfile:** `docker/Dockerfile` — builds Go binary, names it `neuron-agent` (hyphen).

**Entrypoint:** `docker/docker-entrypoint.sh` — checks for and executes `/app/neuronagent` (no hyphen).

**Result:** Container starts, entrypoint cannot find binary, exits with "Binary not found." The Docker image does not work.

**Compose file:** `docker/docker-compose.neuronsql.yml` — uses `postgres:17-alpine` (plain PostgreSQL, no NeuronDB extension). Schema migrations require NeuronDB extension. First-run migrations will fail or produce errors.

**No root-level compose:** There is no `docker-compose.yml` at the repository root. The existing compose file is in `docker/` and is framed as a NeuronSQL demo, not a complete stack.

**Additional Dockerfiles:** `docker/Dockerfile.cuda`, `docker/Dockerfile.rocm`, `docker/Dockerfile.metal`, `docker/Dockerfile.pglang` — hardware-specific variants.

---

## CI/CD

### Status: Manual only — no automatic triggers

**CI file:** `.github/workflows/neuron-agent-build-matrix.yml`

**Trigger:** `workflow_dispatch` only. **No `push` or `pull_request` triggers.**

**Jobs:**
- Ubuntu + macOS matrix (Go 1.24): `make build`, `make test-fast`, `golangci-lint`
- Rocky Linux 9 container: same steps

**Gaps:**
- No automatic CI on push or PR
- No Docker build in CI
- No smoke test in CI
- No migration test in CI
- No Docker image publish
- No release automation
- `Makefile` `integration-test` target uses `|| true` — failures are silently swallowed

---

## Scripts

### Status: Exists — operational scripts, no install or demo script

**Scripts in `scripts/`:**
- `neuronagent-setup.sh` — server setup helper
- `neuronagent-run.sh` — run the agent server
- `neuronagent-run-server.sh` — run server variant
- `neuronagent-migrate.sh` — run SQL migrations
- `neuronagent-generate-keys.sh` — generate API keys
- `neuronagent-verify.sh` — verify installation
- `demo_neuronsql_generate.sh` — NeuronSQL demo
- `demo_neuronsql_optimize.sh` — NeuronSQL demo
- `lib/neuronagent-cli.sh` — CLI helper

**Missing scripts:**
- `scripts/install.sh` — one-command installer
- `scripts/demo.sh` — golden path demo
- `scripts/bootstrap-demo.sh` — create demo agent and workspace
- `scripts/wait-for-health.sh` — health polling
- `scripts/reset.sh` — clean state reset
- `scripts/benchmark.sh` — performance benchmark
- `scripts/record-demo.sh` — demo GIF recording

---

## Tests

### Status: Extensive test files; coverage depth unverified; CI is manual-only

**Go tests:** 86 `*_test.go` files across `src/internal/`, `src/cmd/`, `src/cli/`, `src/pkg/`, `src/tests/`.

**Go test coverage gaps:**
- `src/internal/api/` has only `errors_test.go`, `context_test.go`, `request_id_test.go` — most HTTP handlers have no Go unit tests
- No confirmed coverage percentage for core packages

**Python tests:** 150+ `test_*.py` files in `src/tests/`. Organized with pytest markers: `requires_db`, `requires_server`, `requires_neurondb`. Covers API, tools, security, workflow, NeuronDB integration, workers.

**Python test gap:** Coverage target `--cov=NeuronAgent` assumes a Python package layout that may not match the Go-first repository structure without the Python client installed.

**No smoke test:** There is no automated golden path test that runs on CI.

---

## Documentation

### Status: Exists but disorganized; mix of formats; no index

**Docs directory structure:**
```
docs/
├── api-reference.md
├── api.md
├── architecture.md
├── features.md
├── overview.md
├── troubleshooting.md
├── release_checklist.md
├── deployment_guide.txt
├── operations_runbook.txt
├── config_env_schema.txt
├── compliance_profiles.txt
├── security_model.txt
├── product_gap_report.txt
├── product_readiness_audit.txt
├── architecture_v2.txt
├── neuronsql_design.txt
├── neuronsql/
│   ├── quickstart.txt
│   ├── api.txt
│   ├── architecture.txt
│   ├── repo_map.txt
│   ├── security.txt
│   └── eval.txt
└── workflow_templates/
    ├── README.md
    └── *.yaml
```

**Issues:**
- Mix of `.md` and `.txt` files — `.txt` files are not rendered by GitHub or doc sites
- No `docs/README.md` or `docs/index.md` as navigation entry point
- No dedicated docs for: memory, RAG, tools, workflows, security, observability, deployment, contributing
- Hosted docs at neurondb.ai are canonical narrative; local docs are reference/audit-oriented

**README issues:**
- Go version badge says `1.23+`; `go.mod` declares `1.24.0`; prerequisites section says `1.24+`
- Troubleshooting references `docker compose logs agent-server`; compose service is named `neuronagent`
- `CHANGELOG.md` references `Docs/` (capital D); actual path is `docs/` (lowercase)

---

## Examples

### Status: Minimal — only NeuronSQL examples exist

**`examples/` directory:**
- `examples/python/` — `neuronsql_minimal.py`, README
- `examples/typescript/` — `neuronsql_minimal.ts`, `package.json`, README

**`src/examples/` directory:**
- Python client library (`neurondb_client/`)
- Modular examples
- `go_client.go`
- Shell/SQL helpers

**Missing examples:**
- `examples/quickstart-chat/` — basic agent chat
- `examples/rag-with-docs/` — document ingest and RAG
- `examples/sql-agent/` — SQL schema inspection
- `examples/postgres-dba-agent/` — full DBA agent
- `examples/workflow-approval/` — approval gate workflow
- `examples/openclaw-bridge/` — OpenClaw integration

---

## Release and Versioning

### Status: No automated release process; no published images

**Version in README badge:** `3.0.0-devel`

**CHANGELOG:** `CHANGELOG.md` uses Keep a Changelog format. `[Unreleased]` section has one bullet. No versioned release entries.

**No `VERSION` file** at repository root.

**Release checklist:** `docs/release_checklist.md` exists as a manual procedure document.

**No Docker Hub publish.** No GitHub Container Registry publish. No image signing. No SBOM.

---

## Observability

### Status: Partially wired — Prometheus and logging exist; completeness unverified

**Prometheus:** `/metrics` endpoint is registered and the `prometheus/client_golang` library is a dependency. Which specific metrics are registered is not confirmed without running the server.

**Structured logging:** `rs/zerolog` is the logging library. `LoggingMiddleware` logs requests.

**OpenTelemetry:** `go.opentelemetry.io/otel` is a dependency. Tracing spans are likely partial — dependency exists but instrumentation depth is unverified.

**Request IDs:** `RequestIDMiddleware` adds `X-Request-ID` to every request.

**Health endpoint:** `/health` returns health status. Confirmed in router.

**Missing:**
- `/healthz` — Kubernetes liveness probe alias
- `/readyz` — Kubernetes readiness probe (should check DB connectivity)
- `/version` — version and build information

---

## Security

### Status: Baseline exists; production hardening incomplete

**What is confirmed:**
- Auth middleware on all routes except `/health`, `/metrics`
- API key validation via `auth.APIKeyManager`
- Rate limiting via `auth.RateLimiter`
- Admin role check for admin endpoints
- Security headers middleware (`X-Content-Type-Options`, `X-Frame-Options`)
- CORS middleware with configurable origins

**What is unclear or incomplete:**
- JWT implementation depth — config key exists, handler implementation not confirmed
- Full RBAC enforcement — role model exists in code but per-resource permission check coverage is not verified
- Production mode refusing weak secrets — `ValidateConfig` exists but strictness is unclear
- Audit log completeness — tables likely exist in schema but not all event types are confirmed to produce entries

**Gap:** No `docs/security.md`. No `SECURITY.md` for vulnerability reporting.

---

## Helm Chart

### Status: Exists — completeness and correctness unverified

**Location:** `helm/Chart.yaml`, `helm/values.yaml`, `helm/templates/`

**Not verified:** Whether all templates are complete, whether values.yaml reflects all required config variables.

---

## Summary Scorecard

| Area | Status | Score |
|---|---|---|
| HTTP server | Working | 8/10 |
| API surface | Working (50+ endpoints) | 8/10 |
| Authentication | Working (API keys) | 7/10 |
| Agent runtime | Working (legacy path) | 7/10 |
| Memory system | Implemented | 7/10 |
| RAG system | Partial (conditional ingest, no query endpoint) | 5/10 |
| Tool system | Base tools working, NeuronDB tools missing | 5/10 |
| Workflow engine | Basic execution working, scheduling/retries missing | 5/10 |
| Configuration | Working, production hardening weak | 6/10 |
| Database/migrations | Works but two inconsistent approaches | 5/10 |
| Docker/containers | **Broken** (binary name mismatch) | 1/10 |
| CI/CD | Manual only, no automation | 2/10 |
| Install experience | No install script, no demo script | 2/10 |
| Documentation | Exists but disorganized | 4/10 |
| Examples | Minimal (NeuronSQL only) | 3/10 |
| Release process | No automation | 2/10 |
| Observability | Partial (metrics + logging, completeness unverified) | 6/10 |
| Security hardening | Baseline only | 5/10 |
| **Overall** | | **5/10** |

The core runtime is solid. The operational layer (install, Docker, CI, docs, examples) needs significant work before the project is adoption-ready.

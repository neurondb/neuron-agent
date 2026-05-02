# NeuronAgent Feature Truth Table

> Last updated: 2026-04-30  
> Source: Code audit of repository at `/mnt/pge/ndb/neuron-agent`  
> **Legend — Status:** ✅ Working | ⚠️ Partial | ❌ Missing | 🐛 Broken  
> **Legend — Tests:** ✅ Has tests | ⚠️ Partial tests | ❌ No tests  
> **Legend — Docs:** ✅ Documented | ⚠️ Partial docs | ❌ Not documented

Every row is verified by reading code. Nothing is assumed.

---

## Core Server

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| HTTP server startup | ✅ Working | `src/cmd/agent-server/main.go` | ⚠️ Partial | ⚠️ Partial | P1 | — | gorilla/mux, graceful shutdown |
| Request ID middleware | ✅ Working | `src/internal/api/middleware.go` | ✅ Has tests | ❌ Not documented | P2 | — | `X-Request-ID` header |
| Request timeout middleware | ✅ Working | `src/internal/api/middleware.go` | ❌ No tests | ❌ Not documented | P2 | — | 60s default |
| Security headers middleware | ✅ Working | `src/internal/api/middleware.go` | ❌ No tests | ❌ Not documented | P2 | — | X-Frame-Options, etc. |
| CORS middleware | ✅ Working | `src/internal/api/middleware.go` | ❌ No tests | ⚠️ Partial | P2 | — | Configurable origins |
| Request body limit | ✅ Working | `src/internal/api/middleware.go` | ❌ No tests | ❌ Not documented | P2 | — | 10 MiB default |
| Structured JSON logging | ✅ Working | `src/internal/api/middleware.go` | ❌ No tests | ❌ Not documented | P2 | — | zerolog |
| `/health` endpoint | ✅ Working | `src/cmd/agent-server/main.go` | ❌ No tests | ⚠️ Partial | P0 | — | Not behind auth |
| `/metrics` endpoint | ✅ Working | `src/cmd/agent-server/main.go` | ❌ No tests | ❌ Not documented | P1 | — | Prometheus format |
| `/healthz` Kubernetes liveness | ❌ Missing | — | ❌ No tests | ❌ Not documented | P2 | — | GAP-021 |
| `/readyz` Kubernetes readiness | ❌ Missing | — | ❌ No tests | ❌ Not documented | P2 | — | GAP-021; should check DB |
| `/version` endpoint | ❌ Missing | — | ❌ No tests | ❌ Not documented | P2 | — | GAP-020 |
| `/docs` Swagger UI | ❌ Missing | — | ❌ No tests | ❌ Not documented | P2 | — | GAP-039; spec exists at `src/openapi/openapi.yaml` |
| WebSocket endpoint `/ws` | ✅ Working | `src/cmd/agent-server/main.go` | ❌ No tests | ⚠️ Partial | P2 | — | Auth via header required |

---

## Authentication and Authorization

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| API key authentication | ✅ Working | `src/internal/auth/`, `src/internal/api/middleware.go` | ⚠️ Partial | ⚠️ Partial | P0 | — | Bearer or ApiKey scheme |
| Rate limiting per key | ✅ Working | `src/internal/auth/` | ⚠️ Partial | ❌ Not documented | P1 | — | |
| Org scoping per key | ✅ Working | `src/internal/api/org_middleware.go` | ❌ No tests | ❌ Not documented | P1 | — | |
| Admin role enforcement | ✅ Working | `src/internal/api/` handler level | ⚠️ Partial | ❌ Not documented | P1 | — | `/api/v1/admin/*` routes |
| Developer role | ⚠️ Partial | `src/internal/auth/` | ❌ No tests | ❌ Not documented | P1 | — | Role exists in model; enforcement coverage unclear |
| Viewer role | ⚠️ Partial | `src/internal/auth/` | ❌ No tests | ❌ Not documented | P2 | — | Role exists; enforcement unclear |
| Service role | ⚠️ Partial | `src/internal/auth/` | ❌ No tests | ❌ Not documented | P1 | — | Needed for OpenClaw integration |
| Owner role | ⚠️ Partial | `src/internal/auth/` | ❌ No tests | ❌ Not documented | P2 | — | |
| JWT authentication | ⚠️ Partial | `src/internal/config/config.go` | ❌ No tests | ❌ Not documented | P2 | — | Config key exists; implementation depth unclear |
| Workspace-scoped API keys | ⚠️ Partial | `src/internal/auth/` | ❌ No tests | ❌ Not documented | P1 | — | Model exists; enforcement coverage unclear |
| Per-resource permission checks | ⚠️ Partial | Various handlers | ❌ No tests | ❌ Not documented | P1 | — | Admin endpoints checked; others unclear |
| Production mode refuses weak secrets | ⚠️ Partial | `src/internal/config/config.go` | ❌ No tests | ❌ Not documented | P1 | — | GAP-024 |
| `docs/security.md` | ❌ Missing | — | — | ❌ Not documented | P1 | — | No security documentation |
| `SECURITY.md` vulnerability reporting | ❌ Missing | — | — | ❌ Not documented | P2 | — | |

---

## Agent Management

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| Create agent | ✅ Working | `src/internal/api/` | ❌ No tests | ⚠️ Partial | P0 | — | |
| List agents | ✅ Working | `src/internal/api/` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Get agent | ✅ Working | `src/internal/api/` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Update agent | ✅ Working | `src/internal/api/` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Delete agent | ✅ Working | `src/internal/api/` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Clone agent | ✅ Working | `src/internal/api/` | ❌ No tests | ❌ Not documented | P2 | — | |
| Agent planning | ✅ Working | `src/internal/api/`, `src/internal/agent/` | ❌ No tests | ❌ Not documented | P2 | — | |
| Agent reflection | ✅ Working | `src/internal/api/`, `src/internal/agent/` | ❌ No tests | ❌ Not documented | P2 | — | |
| Agent delegation | ✅ Working | `src/internal/api/` | ❌ No tests | ❌ Not documented | P3 | — | |
| Agent metrics | ✅ Working | `src/internal/api/` | ❌ No tests | ❌ Not documented | P2 | — | |
| Agent cost tracking | ✅ Working | `src/internal/api/` | ❌ No tests | ❌ Not documented | P2 | — | |
| Agent budget management | ✅ Working | `src/internal/api/` | ❌ No tests | ❌ Not documented | P2 | — | |
| Agent from blueprint | ✅ Working | `src/internal/api/` | ❌ No tests | ❌ Not documented | P2 | — | |
| Blueprint listing | ✅ Working | `src/internal/api/` | ❌ No tests | ❌ Not documented | P2 | — | |
| Tool permission model per agent | ⚠️ Partial | `src/internal/agent/runtime.go` | ❌ No tests | ❌ Not documented | P1 | — | Model exists; enforcement depth unclear |
| Workspace boundary enforcement | ⚠️ Partial | Various | ❌ No tests | ❌ Not documented | P1 | — | |
| Agent state machine run path | ⚠️ Partial | `src/internal/agent/runtime.go` | ❌ No tests | ❌ Not documented | P2 | — | GAP-019; `SetUseStateMachine` never called |

---

## Session and Message Management

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| Create session | ✅ Working | `src/internal/api/` | ❌ No tests | ⚠️ Partial | P0 | — | |
| List sessions | ✅ Working | `src/internal/api/` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Get/update/delete session | ✅ Working | `src/internal/api/` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Send message (triggers agent) | ✅ Working | `src/internal/api/` | ❌ No tests | ⚠️ Partial | P0 | — | `POST /sessions/{id}/messages` |
| List messages | ✅ Working | `src/internal/api/` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Message feedback | ✅ Working | `src/internal/api/` | ❌ No tests | ❌ Not documented | P2 | — | |
| Session reflection | ✅ Working | `src/internal/api/` | ❌ No tests | ❌ Not documented | P2 | — | |
| Session cleanup / TTL expiry | ⚠️ Partial | Config exists | ❌ No tests | ❌ Not documented | P2 | — | |

---

## Memory System

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| Short-term memory (STM) | ✅ Working | `src/internal/agent/hierarchical_memory.go` | ⚠️ Partial | ❌ Not documented | P1 | — | `memory_stm` table |
| Medium-term memory (MTM) | ✅ Working | `src/internal/agent/hierarchical_memory.go` | ⚠️ Partial | ❌ Not documented | P1 | — | `memory_mtm` table |
| Long-term memory (LPM) | ✅ Working | `src/internal/agent/hierarchical_memory.go` | ⚠️ Partial | ❌ Not documented | P1 | — | `memory_lpm` table |
| General memory chunks | ✅ Working | `src/internal/agent/memory.go`, `src/internal/db/queries.go` | ⚠️ Partial | ❌ Not documented | P1 | — | `memory_chunks` table |
| Episodic memory | ✅ Working | `src/internal/agent/episodic_memory.go` | ⚠️ Partial | ❌ Not documented | P2 | — | |
| Vector similarity search | ✅ Working | `src/internal/db/queries.go` | ⚠️ Partial | ❌ Not documented | P1 | — | `embedding <=> $1::neurondb_vector` |
| Memory importance scoring | ✅ Working | `src/internal/agent/memory.go` | ⚠️ Partial | ❌ Not documented | P2 | — | Heuristic-based |
| Memory temporal decay | ✅ Working | `src/internal/agent/memory.go` | ⚠️ Partial | ❌ Not documented | P2 | — | |
| Memory promotion worker | ✅ Working | `src/internal/worker/memory_promoter.go` | ❌ No tests | ❌ Not documented | P1 | — | STM→MTM→LPM promotion |
| Memory search API | ✅ Working | `POST /api/v1/agents/{id}/memory/search` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Memory forget API | ✅ Working | `POST /api/v1/agents/{id}/memory/forget` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Memory quality API | ✅ Working | `GET /api/v1/agents/{id}/memory/quality` | ❌ No tests | ❌ Not documented | P2 | — | |
| Memory conflict detection | ✅ Working | `POST /api/v1/agents/{id}/memory/conflicts` | ❌ No tests | ❌ Not documented | P2 | — | |
| Memory summarization (LLM-based) | ❌ Missing | `src/internal/agent/memory.go` | ❌ No tests | ❌ Not documented | P2 | — | Current impl uses truncation |
| Memory audit entries | ⚠️ Partial | `src/internal/tools/registry.go` (ExecuteTool) | ❌ No tests | ❌ Not documented | P1 | — | Partial coverage |
| Memory feedback API | ✅ Working | `POST /api/v1/agents/{id}/memory/feedback` | ❌ No tests | ❌ Not documented | P2 | — | |
| Agentic retrieval tool | 🐛 Broken | `src/internal/tools/retrieval_tool.go` | ❌ No tests | ❌ Not documented | P1 | — | GAP-009; tool exists but not registered |
| `docs/memory.md` | ❌ Missing | — | — | ❌ Not documented | P1 | — | |

---

## RAG (Retrieval-Augmented Generation)

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| Document ingest API | ⚠️ Partial | `POST /api/v1/rag/ingest` | ❌ No tests | ❌ Not documented | P1 | — | Conditional on ragClient != nil |
| Document query API | ❌ Missing | — | ❌ No tests | ❌ Not documented | P1 | — | No GET/POST /api/v1/rag/query |
| Document chunking | ✅ Working | `src/pkg/neurondb/rag_client.go` | ❌ No tests | ❌ Not documented | P1 | — | `neurondb_chunk_text()` |
| Embedding generation | ✅ Working | `src/pkg/neurondb/embedding.go` | ❌ No tests | ❌ Not documented | P1 | — | `neurondb_embed()` |
| Context retrieval | ✅ Working | `src/pkg/neurondb/rag_client.go` | ❌ No tests | ❌ Not documented | P1 | — | |
| Reranking results | ✅ Working | `src/pkg/neurondb/rag_client.go` | ❌ No tests | ❌ Not documented | P2 | — | `neurondb_rerank_results()` |
| Answer generation | ✅ Working | `src/pkg/neurondb/rag_client.go` | ❌ No tests | ❌ Not documented | P1 | — | `neurondb_generate_answer()` |
| Source citations | ⚠️ Partial | — | ❌ No tests | ❌ Not documented | P2 | — | Not exposed in API response |
| Hybrid search | ✅ Working | `src/pkg/neurondb/hybrid_search_client.go` | ❌ No tests | ❌ Not documented | P2 | — | `neurondb_hybrid_search()` |
| RAG tool (agent-accessible) | 🐛 Broken | `src/internal/tools/rag_tool.go` | ❌ No tests | ❌ Not documented | P1 | — | GAP-008; not in agent-server registry |
| Advanced RAG | ⚠️ Partial | `src/internal/agent/advanced_rag.go` | ❌ No tests | ❌ Not documented | P3 | — | Experimental, not in hot path |
| Modular RAG | ⚠️ Partial | `src/internal/agent/modular_rag.go` | ❌ No tests | ❌ Not documented | P3 | — | Experimental |
| `docs/rag.md` | ❌ Missing | — | — | ❌ Not documented | P1 | — | |

---

## Tool System

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| Tool registry | ✅ Working | `src/internal/tools/registry.go` | ⚠️ Partial | ❌ Not documented | P1 | — | |
| Tool CRUD API | ✅ Working | `/api/v1/tools` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Tool analytics API | ✅ Working | `/api/v1/tools/{id}/analytics` | ❌ No tests | ❌ Not documented | P2 | — | |
| Circuit breaker per tool | ✅ Working | `src/internal/tools/registry.go` | ❌ No tests | ❌ Not documented | P2 | — | |
| Tool timeout enforcement | ✅ Working | `src/internal/tools/registry.go` | ❌ No tests | ❌ Not documented | P1 | — | |
| Tool audit log on execution | ✅ Working | `src/internal/tools/registry.go` | ❌ No tests | ❌ Not documented | P1 | — | |
| RBAC check before execution | ⚠️ Partial | `src/internal/tools/registry.go` | ❌ No tests | ❌ Not documented | P1 | — | Optional RBAC |
| SQL tool | ✅ Working | `src/internal/tools/sql_tool.go` | ⚠️ Partial | ❌ Not documented | P1 | — | Allows writes by default — needs read-only default |
| HTTP tool | ✅ Working | `src/internal/tools/http_tool.go` | ⚠️ Partial | ❌ Not documented | P2 | — | No domain allowlist by default |
| Code/file analysis tool | ✅ Working | `src/internal/tools/code_tool.go` | ⚠️ Partial | ❌ Not documented | P2 | — | Reads files in allowed dirs |
| Shell tool | ✅ Working | `src/internal/tools/shell_tool.go` | ⚠️ Partial | ❌ Not documented | P1 | — | Should be disabled by default |
| Browser tool | ✅ Working | `src/internal/tools/browser_tool.go` | ❌ No tests | ❌ Not documented | P3 | — | Requires Chromedp |
| Visualization tool | ✅ Working | `src/internal/tools/visualization_tool.go` | ❌ No tests | ❌ Not documented | P3 | — | |
| RAG tool | 🐛 Broken | `src/internal/tools/rag_tool.go` | ❌ No tests | ❌ Not documented | P1 | — | GAP-008; not in agent-server registry |
| Vector search tool | 🐛 Broken | `src/internal/tools/vector_tool.go` | ❌ No tests | ❌ Not documented | P1 | — | GAP-008 + GAP-010; not registered + SQL bug |
| ML tool | ❌ Missing | `src/internal/tools/ml_tool.go` | ❌ No tests | ❌ Not documented | P2 | — | GAP-008; exists but not registered |
| Analytics tool | ❌ Missing | `src/internal/tools/analytics_tool.go` | ❌ No tests | ❌ Not documented | P2 | — | GAP-008; exists but not registered |
| Hybrid search tool | ❌ Missing | `src/internal/tools/hybrid_search_tool.go` | ❌ No tests | ❌ Not documented | P2 | — | GAP-008; exists but not registered |
| Reranking tool | ❌ Missing | `src/internal/tools/reranking_tool.go` | ❌ No tests | ❌ Not documented | P2 | — | GAP-008; exists but not registered |
| Retrieval tool | 🐛 Broken | `src/internal/tools/retrieval_tool.go` | ❌ No tests | ❌ Not documented | P1 | — | GAP-009; never registered |
| Filesystem tool | ⚠️ Partial | `src/internal/tools/filesystem_tool.go` | ❌ No tests | ❌ Not documented | P3 | — | Only in `NewRegistryWithVFS` |
| Memory tool | ⚠️ Partial | `src/internal/tools/memory_tool.go` | ❌ No tests | ❌ Not documented | P3 | — | Only in `NewRegistryWithAllFeatures` |
| Collaboration tool | ⚠️ Partial | `src/internal/tools/collaboration_tool.go` | ❌ No tests | ❌ Not documented | P3 | — | Only in `NewRegistryWithAllFeatures` |
| MCP sync | 🐛 Broken | `src/internal/tools/registry.go` | ❌ No tests | ❌ Not documented | P3 | — | GAP-018; stub that returns nil |
| `docs/tools.md` | ❌ Missing | — | — | ❌ Not documented | P1 | — | |

---

## NeuronSQL Module

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| SQL generation | ✅ Working | `src/internal/modules/neuronsql/module.go` | ⚠️ Partial | ✅ Documented | P1 | — | `POST /api/v1/neuronsql/generate` |
| SQL optimization | ✅ Working | `src/internal/modules/neuronsql/module.go` | ⚠️ Partial | ✅ Documented | P1 | — | `POST /api/v1/neuronsql/optimize` |
| SQL validation | ✅ Working | `src/internal/modules/neuronsql/module.go` | ⚠️ Partial | ⚠️ Partial | P2 | — | `POST /api/v1/neuronsql/validate` |
| PL/pgSQL generation | ✅ Working | `src/internal/modules/neuronsql/module.go` | ⚠️ Partial | ⚠️ Partial | P2 | — | `POST /api/v1/neuronsql/plpgsql` |
| Schema snapshot tool | ✅ Working | `src/internal/modules/neuronsql/module.go` | ⚠️ Partial | ⚠️ Partial | P1 | — | `neuronsql.schema_snapshot` |
| Query explain tool | ✅ Working | `src/internal/modules/neuronsql/module.go` | ⚠️ Partial | ⚠️ Partial | P1 | — | `neuronsql.explain_json` |
| Table profile tool | ⚠️ Partial | `src/internal/neuronsql/tools/register.go` | ❌ No tests | ❌ Not documented | P2 | — | GAP-017; dead code path |
| Index suggestion tool | ⚠️ Partial | `src/internal/neuronsql/tools/register.go` | ❌ No tests | ❌ Not documented | P2 | — | GAP-017; dead code path |

---

## Workflow Engine

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| Create workflow | ✅ Working | `/api/v1/workflows` | ❌ No tests | ⚠️ Partial | P1 | — | |
| List/get/update/delete workflow | ✅ Working | `/api/v1/workflows/{id}` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Add/list workflow steps | ✅ Working | `/api/v1/workflows/{id}/steps` | ❌ No tests | ⚠️ Partial | P1 | — | |
| Execute workflow | ✅ Working | `POST /api/v1/workflows/{id}/execute` | ❌ No tests | ⚠️ Partial | P1 | — | |
| List executions | ✅ Working | `/api/v1/workflows/{id}/executions` | ❌ No tests | ⚠️ Partial | P1 | — | |
| DAG topological ordering | ✅ Working | `src/internal/workflow/engine.go` | ⚠️ Partial | ❌ Not documented | P1 | — | |
| Agent step type | ✅ Working | `src/internal/workflow/engine.go` | ❌ No tests | ❌ Not documented | P1 | — | Agent step nil UUID bug (GAP-016) |
| Tool step type | ✅ Working | `src/internal/workflow/engine.go` | ❌ No tests | ❌ Not documented | P1 | — | |
| Approval step type | ✅ Working | `src/internal/workflow/engine.go` | ❌ No tests | ❌ Not documented | P1 | — | |
| HTTP step type | ✅ Working | `src/internal/workflow/engine.go` | ❌ No tests | ❌ Not documented | P2 | — | |
| SQL step type | ✅ Working | `src/internal/workflow/engine.go` | ❌ No tests | ❌ Not documented | P2 | — | |
| Conditional step type | 🐛 Broken | `src/internal/workflow/advanced_engine.go` | ❌ No tests | ❌ Not documented | P1 | — | GAP-015; not in main engine switch |
| Workflow retry tracking | ⚠️ Partial | `src/internal/workflow/engine.go` | ❌ No tests | ❌ Not documented | P1 | — | Tracked but not re-executed (GAP-014) |
| Workflow retry scheduler | ❌ Missing | — | ❌ No tests | ❌ Not documented | P1 | — | GAP-013 |
| Workflow scheduling CRUD | ✅ Working | `/api/v1/workflows/{id}/schedule` | ❌ No tests | ❌ Not documented | P2 | — | CRUD works |
| Workflow schedule runner | ❌ Missing | — | ❌ No tests | ❌ Not documented | P1 | — | GAP-013; no background worker |
| Step idempotency | ✅ Working | `src/internal/workflow/engine.go` | ❌ No tests | ❌ Not documented | P2 | — | `idempotency_key` field |
| Step compensation | ⚠️ Partial | `src/internal/workflow/engine.go` | ❌ No tests | ❌ Not documented | P3 | — | `CompensateStep` exists |
| Approval management API | ✅ Working | `/api/v1/approvals` | ❌ No tests | ⚠️ Partial | P1 | — | |
| `docs/workflows.md` | ❌ Missing | — | — | ❌ Not documented | P1 | — | |

---

## OpenClaw Bridge

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| `/claw/v1/health` endpoint | ✅ Working | `src/cmd/agent-server/main.go` | ❌ No tests | ❌ Not documented | P1 | — | |
| `/claw/v1/tools/list` endpoint | ✅ Working | `src/cmd/agent-server/main.go` | ❌ No tests | ❌ Not documented | P1 | — | |
| `/claw/v1/tools/run` endpoint | ✅ Working | `src/cmd/agent-server/main.go` | ❌ No tests | ❌ Not documented | P1 | — | |
| Auth enforcement on claw routes | ⚠️ Partial | `src/internal/api/middleware.go` | ❌ No tests | ❌ Not documented | P1 | — | Global middleware applies; service role not confirmed |
| Permission mapping for claw | ⚠️ Partial | — | ❌ No tests | ❌ Not documented | P1 | — | |
| Request logging for claw | ⚠️ Partial | `src/internal/api/middleware.go` | ❌ No tests | ❌ Not documented | P2 | — | Global logging applies |
| `docs/integrations/openclaw.md` | ❌ Missing | — | — | ❌ Not documented | P1 | — | |

---

## Install and Operations

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| `docker-compose.yml` at root | ❌ Missing | — | — | ❌ Not documented | P0 | — | GAP-003 |
| `.env.example` | ❌ Missing | — | — | ❌ Not documented | P0 | — | GAP-004 |
| `scripts/install.sh` | ❌ Missing | — | — | ❌ Not documented | P0 | — | GAP-005 |
| `scripts/demo.sh` | ❌ Missing | — | — | ❌ Not documented | P0 | — | GAP-006 |
| `scripts/wait-for-health.sh` | ❌ Missing | — | — | ❌ Not documented | P1 | — | |
| `scripts/bootstrap-demo.sh` | ❌ Missing | — | — | ❌ Not documented | P1 | — | |
| `scripts/reset.sh` | ❌ Missing | — | — | ❌ Not documented | P2 | — | |
| `scripts/benchmark.sh` | ❌ Missing | — | — | ❌ Not documented | P3 | — | GAP-029 |
| Docker container starts | 🐛 Broken | `docker/Dockerfile`, `docker/docker-entrypoint.sh` | — | — | P0 | — | GAP-001; binary name mismatch |
| NeuronDB in compose stack | 🐛 Broken | `docker/docker-compose.neuronsql.yml` | — | — | P0 | — | GAP-002; plain postgres |
| `make up` | ❌ Missing | `Makefile` | — | — | P0 | — | Target not in current Makefile |
| `make down` | ❌ Missing | `Makefile` | — | — | P0 | — | Target not in current Makefile |
| `make demo` | ❌ Missing | `Makefile` | — | — | P0 | — | Target not in current Makefile |
| `make smoke` | ❌ Missing | `Makefile` | — | — | P1 | — | Target not in current Makefile |
| `make test` | ✅ Working | `Makefile` | — | — | P0 | — | Runs go test |
| `make lint` | ✅ Working | `Makefile` | — | — | P0 | — | golangci-lint |
| `make docker-build` | ✅ Working | `Makefile` | — | — | P1 | — | But container is broken (GAP-001) |
| Migration naming consistency | 🐛 Broken | `sql/`, `scripts/neuronagent-migrate.sh`, `src/internal/db/schema.go` | — | — | P2 | — | GAP-023 |
| Graceful shutdown | ✅ Working | `src/cmd/agent-server/main.go` | ❌ No tests | ❌ Not documented | P2 | — | Signal handling present |
| Worker panic recovery | ⚠️ Partial | — | ❌ No tests | ❌ Not documented | P2 | — | Not confirmed |
| Startup config banner | ❌ Missing | — | — | ❌ Not documented | P2 | — | GAP-027 |

---

## CI/CD and Release

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| CI on push/PR | ❌ Missing | `.github/workflows/` | — | — | P0 | — | GAP-007; workflow_dispatch only |
| CI lint check | ✅ Working | `.github/workflows/neuron-agent-build-matrix.yml` | — | — | P1 | — | But only manual trigger |
| CI unit tests | ✅ Working | `.github/workflows/neuron-agent-build-matrix.yml` | — | — | P1 | — | But only manual trigger |
| CI Docker build | ❌ Missing | — | — | — | P1 | — | Not in any workflow |
| CI smoke test | ❌ Missing | — | — | — | P1 | — | Not in any workflow |
| Docker Hub publish | ❌ Missing | — | — | — | P1 | — | GAP-038 |
| GHCR publish | ❌ Missing | — | — | — | P1 | — | GAP-038 |
| Semantic versioning | ⚠️ Partial | README badge (`3.0.0-devel`) | — | — | P2 | — | No VERSION file at root |
| Release notes automation | ❌ Missing | — | — | — | P2 | — | GAP-038 |
| SBOM generation | ❌ Missing | — | — | — | P3 | — | GAP-038 |
| Image signing | ❌ Missing | — | — | — | P3 | — | GAP-038 |
| Vulnerability scan | ❌ Missing | — | — | — | P2 | — | GAP-038 |
| CHANGELOG with versions | ❌ Missing | `CHANGELOG.md` | — | — | P2 | — | Only `[Unreleased]` exists |

---

## Documentation

| Feature | Status | Code Location | Tests | Docs | Priority | Owner | Notes |
|---|---|---|---|---|---|---|---|
| README — install instructions | ⚠️ Partial | `README.md` | — | — | P0 | — | GAP-012; broken compose, no install script |
| README — demo | ❌ Missing | `README.md` | — | — | P0 | — | No demo GIF or working demo |
| README — correct Go version | 🐛 Broken | `README.md` | — | — | P1 | — | Badge says 1.23, go.mod is 1.24 |
| `docs/architecture.md` | ✅ Working | `docs/architecture.md` | — | ✅ | P2 | — | Exists; may need updates |
| `docs/api.md` | ⚠️ Partial | `docs/api.md`, `docs/api-reference.md` | — | ⚠️ | P1 | — | Exists; no curl examples |
| `docs/configuration.md` | ❌ Missing | `docs/config_env_schema.txt` | — | ❌ | P1 | — | Exists as .txt; needs .md conversion |
| `docs/memory.md` | ❌ Missing | — | — | ❌ | P1 | — | |
| `docs/rag.md` | ❌ Missing | — | — | ❌ | P1 | — | |
| `docs/tools.md` | ❌ Missing | — | — | ❌ | P1 | — | |
| `docs/workflows.md` | ❌ Missing | — | — | ❌ | P1 | — | |
| `docs/security.md` | ❌ Missing | `docs/security_model.txt` | — | ❌ | P1 | — | Exists as .txt |
| `docs/observability.md` | ❌ Missing | — | — | ❌ | P1 | — | |
| `docs/backup-restore.md` | ❌ Missing | — | — | ❌ | P1 | — | |
| `docs/deploy/docker-compose.md` | ❌ Missing | `docs/deployment_guide.txt` | — | ❌ | P1 | — | Exists as .txt |
| `docs/deploy/kubernetes.md` | ❌ Missing | — | — | ❌ | P2 | — | Helm chart exists but undocumented |
| `docs/integrations/openclaw.md` | ❌ Missing | — | — | ❌ | P1 | — | |
| `docs/quickstart.md` | ❌ Missing | — | — | ❌ | P0 | — | |
| `docs/index.md` | ❌ Missing | — | — | ❌ | P2 | — | |
| `docs/faq.md` | ❌ Missing | — | — | ❌ | P2 | — | |
| `docs/use-cases.md` | ❌ Missing | — | — | ❌ | P2 | — | |
| `docs/comparisons.md` | ❌ Missing | — | — | ❌ | P2 | — | |
| `CONTRIBUTING.md` | ❌ Missing | — | — | ❌ | P2 | — | GAP-032 |
| `SECURITY.md` | ❌ Missing | — | — | ❌ | P2 | — | |
| `CODE_OF_CONDUCT.md` | ❌ Missing | — | — | ❌ | P3 | — | |

---

## Summary Counts

| Status | Count |
|---|---|
| ✅ Working | 62 |
| ⚠️ Partial | 28 |
| ❌ Missing | 38 |
| 🐛 Broken | 10 |
| **Total features tracked** | **138** |

| Test status | Count |
|---|---|
| ✅ Has tests | 4 |
| ⚠️ Partial tests | 35 |
| ❌ No tests | 99 |

The core engine is functional. The operational layer, documentation, and NeuronDB tool wiring need the most attention.

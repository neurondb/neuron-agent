# NeuronAgent OpenClaw Bridge Plan

> Last updated: 2026-04-30  
> Target completion: Month 4, Week 4  
> Status: Partial implementation exists; auth, permissions, docs, and tests needed

---

## What OpenClaw Is

OpenClaw is an agent channel gateway. It handles user-facing communication: Slack messages, Teams threads, web chat interfaces, API gateway routing, and message fan-out. OpenClaw knows how to talk to users across channels.

OpenClaw does not store memory. OpenClaw does not reason over data. OpenClaw does not execute tools with business logic. OpenClaw routes.

## What NeuronAgent Is

NeuronAgent is an agent runtime with a transactional brain. It stores durable memory, executes RAG over private documents, runs SQL-native reasoning, executes tools with audit logging, manages multi-step workflows with approval gates, and maintains complete audit trails. NeuronAgent knows how to think reliably.

## The Integration

```
User → Slack / Teams / Web Chat
         │
         ▼
      OpenClaw (channel gateway)
         │ POST /claw/v1/tools/run
         ▼
      NeuronAgent (agent runtime)
         │ NeuronDB functions
         ▼
      PostgreSQL + NeuronDB (database-native AI)
```

OpenClaw calls NeuronAgent as a tool provider. When a user sends a message that requires intelligence — memory recall, document search, database inspection, workflow triggers — OpenClaw delegates to NeuronAgent. NeuronAgent returns structured results. OpenClaw formats and delivers them back to the user.

---

## Current State

The claw gateway routes are already registered in `src/cmd/agent-server/main.go`:

```
GET  /claw/v1/health
POST /claw/v1/tools/list
POST /claw/v1/tools/run
```

These endpoints exist in the router. What is missing:
- Auth is applied via global middleware (correct) but `service` role is not confirmed
- The tools that these endpoints expose are not fully defined
- No documentation exists for OpenClaw integration
- No example configuration exists
- No tests exist for the claw endpoints

---

## Month 4 Week 4 Implementation Plan

### Step 1: Auth and permissions (Day 1)

**Goal:** Claw endpoint calls require a `service` role API key and nothing else.

**Tasks:**
1. Confirm that global `AuthMiddleware` applies to `/claw/v1` routes (it does — `claw` routes are under the main router)
2. Verify that the `service` role is defined in `src/internal/auth/`
3. Add explicit `RequireRole(..., RoleService)` check in the claw handler functions
4. Add service role to the API key creation flow with documentation
5. Write a test: claw endpoint with service key returns 200; claw endpoint with developer key returns 403

**Files to change:**
- `src/internal/api/` — claw handler files
- `src/internal/auth/` — verify service role definition

---

### Step 2: Define the tool catalog (Day 1–2)

**Goal:** Define exactly which tools OpenClaw can invoke via the bridge, with input/output schemas.

The claw bridge exposes a subset of NeuronAgent capabilities as named tools. Each tool maps to a specific NeuronAgent API call.

#### Tool catalog

| Tool name | Maps to | Description |
|---|---|---|
| `agent.chat` | `POST /api/v1/sessions/{id}/messages` | Send a message to a named agent, return response |
| `rag.answer` | RAG query path | Ask a question over the knowledge base |
| `rag.ingest` | `POST /api/v1/rag/ingest` | Add a document to the knowledge base |
| `memory.search` | `POST /api/v1/agents/{id}/memory/search` | Search agent memory by query string |
| `schema.inspect` | `neuronsql.schema_snapshot` tool | Return schema summary for the workspace database |
| `sql.read` | `sql` tool with read-only enforcement | Run a read-only SQL query (requires explicit allow) |
| `workflow.trigger` | `POST /api/v1/workflows/{id}/execute` | Start a named workflow with input parameters |
| `query.explain` | `neuronsql.explain_json` tool | Run EXPLAIN on a SQL query |

**Input schema for `POST /claw/v1/tools/run`:**
```json
{
  "tool": "agent.chat",
  "workspace_id": "ws_123",
  "agent_name": "support-agent",
  "inputs": {
    "message": "What is the status of order #12345?",
    "session_id": "sess_abc"
  }
}
```

**Output schema:**
```json
{
  "tool": "agent.chat",
  "success": true,
  "result": {
    "response": "Order #12345 was shipped on April 28...",
    "session_id": "sess_abc",
    "memory_used": true
  },
  "request_id": "req_xyz",
  "duration_ms": 842
}
```

**For `POST /claw/v1/tools/list`:**
```json
{
  "tools": [
    {
      "name": "agent.chat",
      "description": "Send a message to a named agent and receive a response with memory context",
      "inputs": {
        "message": "string (required)",
        "agent_name": "string (required)",
        "session_id": "string (optional)"
      }
    },
    ...
  ]
}
```

---

### Step 3: Implement tool routing (Day 2–3)

**Goal:** `POST /claw/v1/tools/run` routes to the correct NeuronAgent handler based on `tool` field.

**Implementation approach:**

Create a `ClawRouter` in `src/internal/api/claw_handlers.go` (or extend existing claw handler):

```go
type ClawRouter struct {
    agents   *AgentHandlers
    sessions *SessionHandlers
    tools    *ToolHandlers
    rag      *RAGHandlers
    workflow *WorkflowHandlers
    queries  *db.Queries
}

func (cr *ClawRouter) RunTool(w http.ResponseWriter, r *http.Request) {
    var req ClawToolRunRequest
    // validate auth, decode request, route to handler
    switch req.Tool {
    case "agent.chat":
        cr.runAgentChat(w, r, req)
    case "rag.answer":
        cr.runRAGAnswer(w, r, req)
    // ...
    default:
        writeError(w, 400, "unknown_tool", "Tool not found: "+req.Tool)
    }
}
```

Key requirements for each routed call:
- Validate that the requesting API key has permission to use this tool
- Validate `workspace_id` matches the API key scope
- Validate all required inputs are present
- Execute the underlying NeuronAgent capability
- Return structured result in the claw response schema
- Write an audit log entry for every claw tool run

---

### Step 4: Permission mapping (Day 3)

**Goal:** Define which tools a service API key can invoke, and enforce it.

**Default service role permissions:**

| Tool | Service role can invoke |
|---|---|
| `agent.chat` | Yes — if agent belongs to the key's workspace |
| `rag.answer` | Yes |
| `rag.ingest` | Yes — with write permission enabled on key |
| `memory.search` | Yes — read only |
| `schema.inspect` | Yes — read only |
| `sql.read` | Yes — only if explicitly allowed in bridge config (`CLAW_ALLOW_SQL_READ=true`) |
| `workflow.trigger` | Yes — if workflow belongs to the key's workspace |
| `query.explain` | Yes — read only |

**SQL tool safety:** `sql.read` must use the SQL tool with `allow_writes: false`. The bridge config flag `CLAW_ALLOW_SQL_READ` defaults to `false`. It must be explicitly enabled.

---

### Step 5: Request logging and audit (Day 3)

**Goal:** Every claw tool call produces a structured log entry and an audit event.

**Log entry fields:**
- `request_id` — correlates with HTTP middleware
- `claw_tool` — tool name
- `workspace_id`
- `actor_id` — API key ID
- `inputs_hash` — SHA-256 of input JSON (not the inputs themselves)
- `status` — `success`, `error`, `denied`
- `duration_ms`
- `result_size_bytes`

**Audit event:** type `claw_tool_run`, includes tool name, workspace, status. Uses the standard audit log schema.

---

### Step 6: Tests (Day 4)

**Required tests:**

| Test | What it verifies |
|---|---|
| `TestClawHealth` | `GET /claw/v1/health` returns 200 |
| `TestClawToolsList` | `POST /claw/v1/tools/list` returns tool catalog |
| `TestClawRunAgentChat` | `agent.chat` tool routes to agent and returns response |
| `TestClawRunRAGAnswer` | `rag.answer` returns answer for test question |
| `TestClawAuthRequired` | Without API key, claw endpoints return 401 |
| `TestClawServiceRoleRequired` | Developer key on claw endpoint returns 403 |
| `TestClawWorkspaceBoundary` | Cannot invoke tool on agent from different workspace |
| `TestClawUnknownTool` | Unknown tool name returns 400 with clear error |
| `TestClawSQLDefaultDisabled` | `sql.read` returns 403 when `CLAW_ALLOW_SQL_READ=false` |
| `TestClawAuditEntry` | Tool run produces an audit log entry |

---

### Step 7: Documentation and example (Day 4–5)

**Create `docs/integrations/openclaw.md` covering:**

1. What OpenClaw is and how it relates to NeuronAgent
2. Setting up the integration:
   - Create a service role API key
   - Configure OpenClaw with the NeuronAgent base URL and API key
   - Test with `GET /claw/v1/health`
3. Available tools — full catalog with input/output schemas and examples
4. Permission model — what a service key can and cannot do
5. SQL safety — why `sql.read` is disabled by default
6. Request correlation — using `X-Request-ID` for tracing calls across systems
7. Error handling — error response schema and how to handle tool failures
8. Rate limits — the service key uses the same rate limiter as other keys
9. Example OpenClaw configuration YAML

**Create `examples/openclaw-bridge/`:**
- `README.md` — setup instructions and what the example shows
- `test-claw-endpoints.sh` — curl commands testing each endpoint
- `sample-openclaw-config.yaml` — example OpenClaw tool configuration pointing to NeuronAgent
- `demo.sh` — end-to-end demo: send user message through simulated OpenClaw → NeuronAgent → return response

---

## Architecture Notes

### Why these tools, not the full API

The claw bridge exposes a read-oriented, action-constrained subset of NeuronAgent. This is intentional:
- OpenClaw is a routing gateway, not an admin interface
- Service accounts should have minimal permissions
- The bridge is not a general HTTP proxy — it is a typed tool interface with validation

### Session management

When `agent.chat` is called, the session ID should be:
1. Provided by the OpenClaw caller (preferred — maintains conversation continuity across calls)
2. Created by NeuronAgent if not provided (creates a new session — no memory continuity)

OpenClaw should pass a consistent session ID for each conversation thread.

### Workspace scoping

All claw tool calls include a `workspace_id`. The service API key must be scoped to that workspace. Cross-workspace calls are rejected.

### Idempotency

`workflow.trigger` should accept an `idempotency_key` field that maps to NeuronAgent's `X-Idempotency-Key` header. This prevents duplicate workflow executions if OpenClaw retries a failed request.

---

## Configuration Reference

Environment variables for the claw bridge:

```
CLAW_ALLOW_SQL_READ=false          # Whether sql.read tool is available via claw bridge
CLAW_MAX_AGENT_RESPONSE_BYTES=65536  # Maximum response size for agent.chat
CLAW_WORKFLOW_TIMEOUT_MS=30000     # Maximum time to wait for workflow.trigger
CLAW_AUDIT_ALL_REQUESTS=true       # Whether to audit all requests (not just tool runs)
```

---

## Acceptance Criteria

The Month 4 Week 4 implementation is complete when:

- [ ] `GET /claw/v1/health` returns 200 and is documented
- [ ] `POST /claw/v1/tools/list` returns the defined tool catalog
- [ ] `POST /claw/v1/tools/run` correctly routes all 8 defined tools
- [ ] Service role API key is required; developer/viewer keys are rejected with 403
- [ ] Workspace boundaries are enforced
- [ ] `sql.read` is disabled by default and requires explicit `CLAW_ALLOW_SQL_READ=true`
- [ ] Every tool run produces an audit log entry
- [ ] All 10 defined tests pass
- [ ] `docs/integrations/openclaw.md` is complete
- [ ] `examples/openclaw-bridge/` has working demo scripts
- [ ] Integration is described in README under "OpenClaw bridge" section

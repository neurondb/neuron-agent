# NeuronAgent Gap Analysis

> Last updated: 2026-04-30  
> Source: `docs/roadmap/current-state-audit.md`  
> Priority levels: P0 = critical blocker, P1 = high, P2 = medium, P3 = low

This document lists every known gap, bug, and missing piece in NeuronAgent. Each item has a severity, effort estimate, target month, owner placeholder, and implementation notes. Nothing is invented here — every gap was found by reading code.

---

## P0: Critical Blockers (Fix in Month 1)

These items prevent the project from being installed or used. They must be fixed before any other work has value.

---

### GAP-001: Docker entrypoint binary name mismatch

**Severity:** P0 — Critical  
**Effort:** Small (1–2 hours)  
**Target month:** M1 Week 2  
**Owner:** —

**Description:**  
`docker/Dockerfile` builds the binary and names it `neuron-agent` (with hyphen):
```dockerfile
RUN go build -o /app/neuron-agent ./cmd/agent-server
```
`docker/docker-entrypoint.sh` checks for and executes `/app/neuronagent` (no hyphen):
```sh
if [ ! -f /app/neuronagent ]; then
    echo "Binary not found"
    exit 1
fi
exec /app/neuronagent "$@"
```
The Docker container exits immediately with "Binary not found." No one can run NeuronAgent via Docker.

**Fix:**  
Pick one name and use it consistently. Recommendation: `neuron-agent` (hyphenated, matches repo name convention). Update `docker-entrypoint.sh` to use `neuron-agent`.

**Files to change:**
- `docker/docker-entrypoint.sh`

**Test:** `docker build -f docker/Dockerfile -t test . && docker run --rm test /app/neuron-agent --help`

---

### GAP-002: Docker Compose stack has no NeuronDB extension

**Severity:** P0 — Critical  
**Effort:** Medium (half day)  
**Target month:** M1 Week 2  
**Owner:** —

**Description:**  
`docker/docker-compose.neuronsql.yml` uses `postgres:17-alpine` — plain PostgreSQL with no NeuronDB extension. The primary schema file `sql/neuron-agent.sql` begins with:
```sql
CREATE EXTENSION IF NOT EXISTS neurondb;
```
Running `neuronagent-migrate.sh` against a plain postgres container will fail or produce errors. Vector operations, embeddings, RAG functions, and hybrid search all depend on NeuronDB extension functions being present.

**Fix options:**
1. Build and publish a `neurondb/postgres:17` Docker image that includes the NeuronDB extension (preferred for public use)
2. Add an init-script volume mount that installs the extension from a local `.so` file (for internal/private use)
3. Provide clear documentation that NeuronDB extension must be installed separately, with instructions

For the root compose file, at minimum provide option 3 with a clear error if the extension is not found.

**Files to change:**
- Create `docker-compose.yml` at repository root
- Potentially create `docker/init-scripts/01-neurondb.sh`
- Update `scripts/install.sh` to check for NeuronDB before migration

**Test:** `docker compose up -d && docker compose exec db psql -U postgres -c "SELECT neurondb_version();"`

---

### GAP-003: No root-level docker-compose.yml

**Severity:** P0 — Critical (for install experience)  
**Effort:** Small (2–4 hours)  
**Target month:** M1 Week 2  
**Owner:** —

**Description:**  
The only compose file is `docker/docker-compose.neuronsql.yml`. There is no `docker-compose.yml` at the repository root. The README's quickstart tells users to run `docker compose up -d` from the repository root — this command fails with "no configuration file provided."

**Fix:**  
Create `docker-compose.yml` at repository root with:
- NeuronDB-enabled PostgreSQL service
- NeuronAgent service
- Named volumes
- Health checks with `depends_on` + `service_healthy`
- `.env` file loading

**Files to change:**
- Create `docker-compose.yml` at repository root
- Create `.env.example` at repository root

---

### GAP-004: No .env.example file

**Severity:** P0 — Critical  
**Effort:** Small (1–2 hours)  
**Target month:** M1 Week 2  
**Owner:** —

**Description:**  
There is no `.env.example` at the repository root. New users cannot know what environment variables to set. The install flow requires users to read the source code to find required config.

**Fix:**  
Create `.env.example` with all required variables, safe local defaults, inline comments explaining each variable, and production guidance.

**Minimum required variables:**
```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=neuronagent
DB_USER=neuronagent
DB_PASSWORD=local-only-password
AUTH_API_KEY_SECRET=change-me-before-production
SERVER_PORT=8080
SERVER_MODE=local
LOG_LEVEL=info
LOG_FORMAT=json
```

---

### GAP-005: No install script

**Severity:** P0 — Critical (for adoption)  
**Effort:** Medium (half day)  
**Target month:** M1 Week 2  
**Owner:** —

**Description:**  
There is no `scripts/install.sh` that a new developer can run to set up NeuronAgent. The README implies a simple install but requires manual steps. The target install command is:
```bash
curl -fsSL https://raw.githubusercontent.com/neurondb/neuron-agent/main/scripts/install.sh | bash
```
This script does not exist.

**Fix:**  
Create `scripts/install.sh` that:
1. Checks for Docker 24+ and docker compose (not docker-compose)
2. Checks for git
3. Clones the repo if not already cloned
4. Copies `.env.example` to `.env`
5. Runs `docker compose up -d`
6. Waits for health endpoint via `scripts/wait-for-health.sh`
7. Prints success message with next steps

If any prerequisite is missing, the script must print a clear error with installation instructions. No silent failures.

---

### GAP-006: No demo script

**Severity:** P0 — Critical (for adoption)  
**Effort:** Medium (half day to one day)  
**Target month:** M1 Week 4  
**Owner:** —

**Description:**  
There is no `scripts/demo.sh` that demonstrates NeuronAgent working. The `make demo` target has no backing implementation. A developer who successfully installs NeuronAgent has no guided first experience.

**Fix:**  
Create `scripts/demo.sh` that runs the golden path:
1. Wait for health
2. Generate or load an API key
3. Create an agent
4. Ingest a fixture document
5. Ask a question about the document
6. Run a read-only SQL query
7. Print clean, formatted output at each step
8. Print final pass/fail summary

---

### GAP-007: No CI automation on push or pull_request

**Severity:** P0 — Critical  
**Effort:** Medium (half day)  
**Target month:** M1 Week 4  
**Owner:** —

**Description:**  
`.github/workflows/neuron-agent-build-matrix.yml` has trigger `workflow_dispatch` only. No push or pull_request trigger exists. Every pull request can merge broken code without any automated check.

**Fix:**  
Add `push` and `pull_request` triggers to the existing workflow, or create a new lightweight `ci.yml` for fast PR checks:
```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```
Minimum PR checks: lint + unit tests. Docker build can be a nightly or separate workflow due to startup time.

---

## P1: High Priority (Fix in Month 1–2)

These items prevent NeuronAgent from working correctly for its documented capabilities.

---

### GAP-008: NeuronDB tool handlers not registered in agent-server

**Severity:** P1 — High  
**Effort:** Small (2–4 hours)  
**Target month:** M2 Week 4  
**Owner:** —

**Description:**  
`src/cmd/agent-server/main.go` calls `tools.NewRegistry(...)`, which only registers base handlers (sql, http, code, shell, browser, visualization). The NeuronDB tool handlers — `rag`, `vector`, `ml`, `analytics`, `hybrid_search`, `reranking` — are only registered in `NewRegistryWithNeuronDB`, which is used by the `bench` binary but not `agent-server`.

This means agents cannot use RAG, vector search, ML, analytics, hybrid search, or reranking tools even though the server is connected to a NeuronDB-enabled database.

**Fix:**  
In `src/cmd/agent-server/main.go`, after constructing the NeuronDB client, call `NewRegistryWithNeuronDB` instead of `NewRegistry` when NeuronDB client is successfully initialized:
```go
if neurondbClient != nil {
    registry = tools.NewRegistryWithNeuronDB(db, queries, neurondbClient)
} else {
    registry = tools.NewRegistry(db, queries)
}
```

**Files to change:**
- `src/cmd/agent-server/main.go`

---

### GAP-009: RetrievalTool is never registered

**Severity:** P1 — High  
**Effort:** Small (1–2 hours)  
**Target month:** M2 Week 4  
**Owner:** —

**Description:**  
`src/internal/tools/retrieval_tool.go` defines `RetrievalTool` and `NewRetrievalTool`. The agent runtime, when `agentic_retrieval_enabled` is in config, generates prompts that describe using a "retrieval tool." But `NewRetrievalTool` is never called in any registry construction function. The retrieval tool is referenced in prompts but does not exist in the registry.

**Fix:**  
Register `RetrievalTool` in `NewRegistryWithNeuronDB` (and potentially `NewRegistryWithAllFeatures`). Alternatively, wire it directly in `agent-server/main.go` when agentic retrieval is enabled.

**Files to change:**
- `src/internal/tools/registry.go`
- or `src/cmd/agent-server/main.go`

---

### GAP-010: VectorClient.VectorSearch SQL is broken

**Severity:** P1 — High  
**Effort:** Small (1 hour)  
**Target month:** M2 Week 4  
**Owner:** —

**Description:**  
In `src/pkg/neurondb/vector_client.go`, the `VectorSearch` function builds a SQL query that includes an `ORDER BY embedding <=> $1` clause for similarity search but does not bind the query vector to `$1` — only `limit` is passed as a parameter. The SQL will fail with "wrong number of parameters" or return incorrect results.

**Fix:**  
Pass the query vector as the first parameter in the query execution:
```go
rows, err := v.db.QueryContext(ctx, sql, queryVector, limit)
```

**Files to change:**
- `src/pkg/neurondb/vector_client.go`

---

### GAP-011: Makefile integration-test silently passes failures

**Severity:** P1 — High  
**Effort:** Trivial (5 minutes)  
**Target month:** M1 Week 2  
**Owner:** —

**Description:**  
The `integration-test` Makefile target uses `|| true`:
```makefile
integration-test:
    cd src && go test -short ./... || true
```
This means test failures are silently ignored. CI and local `make integration-test` always exit 0.

**Fix:**  
Remove `|| true`:
```makefile
integration-test:
    cd src && go test -short ./...
```

**Files to change:**
- `Makefile`

---

### GAP-012: README has false claims and version drift

**Severity:** P1 — High  
**Effort:** Small (2–4 hours for a full audit)  
**Target month:** M1 Week 3  
**Owner:** —

**Description:**  
- Badge claims `Go 1.23+`; `src/go.mod` declares `go 1.24.0`
- Prerequisites section says `1.24+` — inconsistent with badge
- Troubleshooting references `docker compose logs agent-server` — service is named `neuronagent` in compose
- `CHANGELOG.md` references `Docs/` (capital D) — path is `docs/` (lowercase)
- Any feature claims that require NeuronDB (and note that NeuronDB may not be available) need to be clearly qualified

**Fix:**  
Audit every claim in README against actual code. Fix version badge. Fix service name. Qualify NeuronDB-dependent features. Restructure README to the target structure.

---

### GAP-013: Workflow schedule runner does not exist

**Severity:** P1 — High  
**Effort:** Large (2–3 days)  
**Target month:** M4 Week 3  
**Owner:** —

**Description:**  
Workflow schedule CRUD endpoints exist (`POST/GET/PUT/DELETE /api/v1/workflows/{id}/schedule`). `ListWorkflowSchedulesByNextRun` exists in the DB layer. But no background worker calls this function. Schedules are stored but never executed. The scheduling feature is non-functional.

**Fix:**  
Create a background worker in `src/internal/worker/workflow_scheduler.go` that:
1. On startup, queries `ListWorkflowSchedulesByNextRun` with `now()` as the cutoff
2. For each due schedule, calls `workflow.Engine.Execute`
3. Updates `next_run` after execution
4. Sleeps until the next due schedule
5. Handles panics and restarts with backoff

Requires a distributed lock or single-writer guarantee to avoid duplicate execution in multi-instance deployments.

---

### GAP-014: Workflow retry scheduling is incomplete

**Severity:** P1 — High  
**Effort:** Medium (1 day)  
**Target month:** M4 Week 3  
**Owner:** —

**Description:**  
The workflow engine tracks retry counts but explicitly does not re-execute failed steps. The code comment states: "retry scheduling is not implemented — would need retry scheduler."

**Fix:**  
Add retry scheduler as part of the workflow scheduler worker (GAP-013). When a step fails and `retry_count < max_retries`, schedule re-execution after `retry_delay_ms`. Use the existing `retry_config` JSONB field on `WorkflowStep`.

---

### GAP-015: Workflow conditional step type not in main engine

**Severity:** P1 — High  
**Effort:** Medium (half day to one day)  
**Target month:** M4 Week 3  
**Owner:** —

**Description:**  
`conditional` step type is defined in:
- `src/internal/workflow/advanced_engine.go` (not wired to HTTP)
- TypeScript SDK types (`type: 'conditional'`)

But the main engine `switch` in `src/internal/workflow/engine.go` does not handle `conditional`. Any workflow with a conditional step will hit the default case and error.

**Fix:**  
Port the conditional step logic from `advanced_engine.go` into the main `engine.go` switch case. Test with a workflow that branches based on a previous step's output.

---

### GAP-016: Workflow agent step uses nil UUID for session lookup

**Severity:** P1 — High  
**Effort:** Small (2–4 hours)  
**Target month:** M4 Week 3  
**Owner:** —

**Description:**  
In `executeAgentStep` in `src/internal/workflow/engine.go`:
```go
session, err := r.queries.GetSession(ctx, uuid.Nil)
```
Using `uuid.Nil` (all zeros) as a session ID will either return a random session, return an error, or create sessions with a zero UUID as parent. This is likely a bug that will cause incorrect behavior in production workflows.

**Fix:**  
Create a new session for each workflow execution's agent step, using the workflow execution ID as the parent context. Store the session ID in the step execution record.

---

### GAP-017: RegisterNeuronSQLTools is dead code

**Severity:** P1 — Medium  
**Effort:** Small (1–2 hours)  
**Target month:** M2 Week 4  
**Owner:** —

**Description:**  
`RegisterNeuronSQLTools` in `src/internal/neuronsql/tools/register.go` registers unprefixed NeuronSQL tools: `optimize_candidates`, `table_profile`, and similar. This function has no callers. These tools are unreachable via any agent.

**Fix options:**
1. Call `RegisterNeuronSQLTools` from within the NeuronSQL module's tool registration (preferred)
2. Remove the function if its tools are redundant with the module-registered tools

---

### GAP-018: SyncFromMCP is a stub

**Severity:** P1 — Medium  
**Effort:** Large (1 week)  
**Target month:** M6  
**Owner:** —

**Description:**  
`SyncFromMCP` in `src/internal/tools/registry.go` returns nil without doing anything. Any code path that calls `SyncFromMCP` silently does nothing.

**Fix:**  
Either implement MCP sync (import external tool definitions from an MCP endpoint) or remove the function and any callers. Do not leave dead stubs in the public API.

---

## P2: Medium Priority (Fix in Month 2–4)

These items limit functionality but do not prevent the core use case from working.

---

### GAP-019: SetUseStateMachine never called

**Severity:** P2  
**Effort:** Small to medium  
**Target month:** M2  
**Owner:** —

**Description:**  
`agent.Runtime.SetUseStateMachine(true)` is never called from `agent-server/main.go`. The state machine run path, which persists run records via `StartRun`, is inactive. Agent execution uses the legacy pipeline.

**Decision needed:** Either enable the state machine (test it first) or document that the legacy path is intentional.

---

### GAP-020: No /version endpoint

**Severity:** P2  
**Effort:** Small (1–2 hours)  
**Target month:** M2 Week 1  
**Owner:** —

**Description:**  
There is no `GET /version` endpoint. Operators cannot query the running server to determine its version, build date, Git commit, or migration version without accessing logs.

**Fix:**  
Add `GET /version` returning:
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
Populate these from build-time linker flags: `-ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)"`.

---

### GAP-021: No /healthz or /readyz Kubernetes probes

**Severity:** P2  
**Effort:** Small (1 hour)  
**Target month:** M3 Week 3  
**Owner:** —

**Description:**  
Kubernetes liveness and readiness probes use `/healthz` and `/readyz` by convention. Only `/health` exists. `/readyz` should perform an actual DB connectivity check and return 503 until the database is connected and migrations are applied.

---

### GAP-022: No /docs or /swagger OpenAPI UI

**Severity:** P2  
**Effort:** Medium (half day)  
**Target month:** M5 Week 1  
**Owner:** —

**Description:**  
`src/openapi/openapi.yaml` exists but there is no runtime endpoint that serves an interactive API explorer. Developers must find the YAML file manually.

**Fix:**  
Serve Swagger UI or Redoc at `/docs` using embedded static assets. This is a one-time setup with no ongoing maintenance cost.

---

### GAP-023: Migration approach is inconsistent

**Severity:** P2  
**Effort:** Medium (half to full day)  
**Target month:** M1 Week 2  
**Owner:** —

**Description:**  
Two migration approaches exist and conflict:
1. `scripts/neuronagent-migrate.sh` — runs all `*.sql` files in `sql/` via `psql`, glob order
2. `src/internal/db/schema.go` — expects `001_name.sql` pattern, version tracking

The main schema file is `neuron-agent.sql`, which does not match the `001_name.sql` pattern expected by the Go migration system.

**Fix:**  
Choose one approach:
- Recommended: rename SQL files to `001_initial_schema.sql`, `002_partitioning.sql`, etc. Use the Go migration system as authoritative.
- Alternative: keep psql script approach but add explicit ordering via filenames.

Update both `scripts/neuronagent-migrate.sh` and `Makefile` `migrate` target to use the chosen approach consistently.

---

### GAP-024: No production config validation for weak secrets

**Severity:** P2  
**Effort:** Small (2–4 hours)  
**Target month:** M2 Week 1  
**Owner:** —

**Description:**  
`src/internal/config/config.go` has `ValidateConfig` but it is unclear whether it refuses weak production defaults (e.g., empty `AUTH_API_KEY_SECRET`, default DB passwords). A server should refuse to start in `production` mode with insecure configuration.

**Fix:**  
Add to `ValidateConfig` for `SERVER_MODE=production`:
- Refuse `AUTH_API_KEY_SECRET` shorter than 32 characters
- Refuse `DB_PASSWORD` that matches common defaults (`password`, `postgres`, `admin`, empty)
- Refuse `SERVER_HOST=0.0.0.0` without explicit acknowledgment

---

### GAP-025: Advanced workflow engine is dead code

**Severity:** P2  
**Effort:** Medium (half day to evaluate)  
**Target month:** M4 Week 3  
**Owner:** —

**Description:**  
`src/internal/workflow/advanced_engine.go` defines an advanced engine with conditional branching. It is not referenced from HTTP handlers or the main engine. It is unreachable code.

**Fix:**  
Either port conditional step support into the main engine (see GAP-015) and delete the advanced engine file, or remove it entirely.

---

### GAP-026: Makefile uses deprecated docker-compose command

**Severity:** P2  
**Effort:** Trivial (10 minutes)  
**Target month:** M1 Week 2  
**Owner:** —

**Description:**  
The Makefile references `docker-compose` (V1, deprecated and removed in many systems). Modern Docker includes `docker compose` (V2, built-in subcommand).

**Fix:**  
Replace all `docker-compose` with `docker compose` in the Makefile.

---

### GAP-027: No startup config summary or banner

**Severity:** P2  
**Effort:** Small (1–2 hours)  
**Target month:** M2 Week 1  
**Owner:** —

**Description:**  
The server starts silently. Operators cannot tell from logs which config was loaded, which mode is running, or what database it connected to without searching through logs.

**Fix:**  
At startup, after config validation, print a summary:
```
NeuronAgent v3.0.0-devel [local]
Listening on 0.0.0.0:8080
Database: neuronagent@localhost:5432/neuronagent
Migrations: v12 applied
NeuronDB: compatible
Workers: memory_promoter=running
```
Never print secret values, passwords, or API keys.

---

### GAP-028: Python test suite coverage target may not work

**Severity:** P2  
**Effort:** Medium (half day to audit)  
**Target month:** M2  
**Owner:** —

**Description:**  
`src/tests/pytest.ini` sets `--cov=NeuronAgent` as coverage target. This assumes a Python package named `NeuronAgent` is installed in the test environment. The repository is Go-first; the Python client library is in `src/examples/neurondb_client/`. The coverage target is likely pointing at a package that is not installed.

**Fix:**  
Either install the Python client as a proper package (`pip install -e src/examples/neurondb_client/`) in the test environment, or update `pytest.ini` to point at the correct module path.

---

## P3: Low Priority (Fix in Month 4–6)

These items are nice-to-have or represent future capabilities.

---

### GAP-029: No benchmark script

**Severity:** P3  
**Effort:** Medium (1 day)  
**Target month:** M6 Week 3  
**Owner:** —

**Description:**  
No `scripts/benchmark.sh` exists. `src/cmd/bench/main.go` exists but is not accessible as a user-facing script. Performance characteristics are undocumented.

---

### GAP-030: No demo GIF

**Severity:** P3  
**Effort:** Small (2–4 hours once demo.sh works)  
**Target month:** M5 Week 4  
**Owner:** —

**Description:**  
README has no demo GIF. First impression requires reading text. A 30-second GIF showing install → demo → output would significantly improve adoption.

---

### GAP-031: Go API handler test coverage is thin

**Severity:** P3  
**Effort:** Large (several days)  
**Target month:** M2  
**Owner:** —

**Description:**  
`src/internal/api/` has only three test files: `errors_test.go`, `context_test.go`, `request_id_test.go`. None of the HTTP handlers for agents, sessions, memory, tools, or workflows have Go unit tests. The API surface is largely untested at the unit level.

---

### GAP-032: No CONTRIBUTING.md

**Severity:** P3  
**Effort:** Small (2–4 hours)  
**Target month:** M5 Week 3  
**Owner:** —

**Description:**  
No `CONTRIBUTING.md` exists. Potential contributors have no guidance on development setup, PR process, coding standards, or test requirements.

---

### GAP-033: No GitHub issue templates

**Severity:** P3  
**Effort:** Small (1–2 hours)  
**Target month:** M5 Week 3  
**Owner:** —

**Description:**  
`.github/ISSUE_TEMPLATE/` does not exist or is empty. Issues filed by community members will lack structure, making triage harder.

---

### GAP-034: Docs mix .md and .txt files

**Severity:** P3  
**Effort:** Small (half day)  
**Target month:** M1  
**Owner:** —

**Description:**  
`docs/` contains `.txt` files (`deployment_guide.txt`, `operations_runbook.txt`, `config_env_schema.txt`, etc.). These are not rendered by GitHub, not searchable as Markdown, and break the docs site structure.

**Fix:**  
Convert `.txt` files to `.md` or move their content into proper doc files. Remove or archive originals.

---

### GAP-035: No agent template files

**Severity:** P3  
**Effort:** Medium (half day per template)  
**Target month:** M5 Week 2  
**Owner:** —

**Description:**  
No `templates/` directory exists. There are YAML examples in `src/cli/examples/` but they are not labeled as templates and have no documentation.

---

### GAP-036: TypeScript SDK conditional step type ahead of Go engine

**Severity:** P3  
**Effort:** Part of GAP-015  
**Target month:** M4 Week 3  
**Owner:** —

**Description:**  
TypeScript SDK `src/sdks/typescript/neurondb-agent/src/types.ts` defines `type: 'conditional'` as a valid workflow step type. Go engine does not support it. The SDK is ahead of the implementation.

---

### GAP-037: Helm chart completeness unverified

**Severity:** P3  
**Effort:** Medium (half day to audit and fix)  
**Target month:** M6 Week 2  
**Owner:** —

**Description:**  
`helm/` directory exists but completeness and correctness have not been verified. Missing items likely include: readiness probe pointing to `/readyz`, resource limits, HPA, PDB, and ServiceMonitor.

---

### GAP-038: No release automation

**Severity:** P3  
**Effort:** Large (1–2 days)  
**Target month:** M6 Week 1  
**Owner:** —

**Description:**  
No GitHub Actions workflow exists for automated release. Docker Hub and GHCR images must be pushed manually. CHANGELOG is not auto-generated. No SBOM. No image signing.

---

### GAP-039: No Swagger UI at runtime

**Severity:** P3  
**Effort:** Small (2–4 hours)  
**Target month:** M5 Week 1  
**Owner:** —

**Description:**  
The OpenAPI spec exists at `src/openapi/openapi.yaml` but is not served at any HTTP endpoint. Developers must find and read the YAML manually.

---

## Gap Summary by Month

### Month 1 (Must fix)
- GAP-001: Docker entrypoint binary name
- GAP-002: Compose has no NeuronDB
- GAP-003: No root-level docker-compose.yml
- GAP-004: No .env.example
- GAP-005: No install script
- GAP-006: No demo script
- GAP-007: No CI on push/PR
- GAP-011: Makefile integration-test silently passes failures
- GAP-012: README false claims and version drift
- GAP-023: Migration approach inconsistency
- GAP-026: Makefile uses deprecated docker-compose
- GAP-034: Docs mix .md and .txt

### Month 2
- GAP-008: NeuronDB tools not registered in agent-server
- GAP-009: RetrievalTool never registered
- GAP-010: VectorClient.VectorSearch SQL broken
- GAP-017: RegisterNeuronSQLTools dead code
- GAP-019: SetUseStateMachine never called (decision needed)
- GAP-020: No /version endpoint
- GAP-024: No production config validation for weak secrets
- GAP-027: No startup config summary
- GAP-028: Python test suite coverage target
- GAP-031: Go API handler test coverage thin

### Month 3
- GAP-021: No /healthz or /readyz
- GAP-032: No CONTRIBUTING.md (may move to M5)

### Month 4
- GAP-013: Workflow schedule runner missing
- GAP-014: Workflow retry scheduling incomplete
- GAP-015: Workflow conditional step not in main engine
- GAP-016: Workflow agent step nil UUID bug
- GAP-025: Advanced workflow engine dead code
- GAP-036: TypeScript SDK conditional step ahead of Go engine

### Month 5
- GAP-022: No /docs Swagger UI
- GAP-030: No demo GIF
- GAP-032: No CONTRIBUTING.md
- GAP-033: No GitHub issue templates
- GAP-035: No agent template files
- GAP-039: No Swagger UI at runtime

### Month 6
- GAP-018: SyncFromMCP stub
- GAP-029: No benchmark script
- GAP-037: Helm chart completeness
- GAP-038: No release automation

---

## Total Gap Count

| Priority | Count |
|---|---|
| P0 — Critical blocker | 7 |
| P1 — High | 11 |
| P2 — Medium | 10 |
| P3 — Low | 11 |
| **Total** | **39** |

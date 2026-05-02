# NeuronAgent Month 1 Execution Board

> Last updated: 2026-04-30  
> Month: May 2026  
> Theme: Foundation, install, and product clarity  
> Goal: A new user installs NeuronAgent locally and runs a working demo in under 5 minutes

Each ticket below is immediately actionable. Every ticket has a file location, implementation notes, test plan, and acceptance criteria. No vague tasks.

---

## How to use this board

- Pick a ticket, assign yourself in the Owner field, move it to In Progress
- When done, check all acceptance criteria
- Open a PR that references this ticket
- Move to Done after PR merges and acceptance criteria are confirmed

---

## Ticket Index

| ID | Title | Priority | Effort | Week | Status |
|---|---|---|---|---|---|
| M1-001 | Fix Docker entrypoint binary name mismatch | P0 | XS | 1 | Open |
| M1-002 | Create root-level docker-compose.yml | P0 | S | 2 | Open |
| M1-003 | Create .env.example | P0 | XS | 2 | Open |
| M1-004 | Fix Makefile: replace docker-compose with docker compose | P0 | XS | 1 | Open |
| M1-005 | Fix Makefile: remove || true from integration-test | P0 | XS | 1 | Open |
| M1-006 | Add CI workflow on push and pull_request | P0 | S | 4 | Open |
| M1-007 | Create scripts/install.sh | P1 | S | 2 | Open |
| M1-008 | Create scripts/wait-for-health.sh | P1 | XS | 2 | Open |
| M1-009 | Create scripts/demo.sh | P1 | M | 4 | Open |
| M1-010 | Create scripts/bootstrap-demo.sh | P1 | S | 4 | Open |
| M1-011 | Add Makefile targets: up, down, demo, smoke, reset | P1 | XS | 2 | Open |
| M1-012 | Create examples/quickstart-chat | P1 | S | 4 | Open |
| M1-013 | Create examples/rag-with-docs | P1 | M | 4 | Open |
| M1-014 | Create examples/sql-agent | P1 | S | 4 | Open |
| M1-015 | Rewrite README to adoption structure | P1 | M | 3 | Open |
| M1-016 | Fix README Go version badge (1.23 → 1.24) | P0 | XS | 1 | Open |
| M1-017 | Create tests/fixtures/sample-doc.txt | P1 | XS | 4 | Open |
| M1-018 | Resolve migration naming inconsistency | P2 | M | 2 | Open |
| M1-019 | Convert docs .txt files to .md | P2 | S | 1 | Open |
| M1-020 | Add smoke test to CI | P1 | M | 4 | Open |

---

## Week 1 Tickets

---

### M1-001: Fix Docker entrypoint binary name mismatch

**Title:** Fix Docker entrypoint binary name mismatch  
**Priority:** P0  
**Effort:** XS (< 1 hour)  
**Week:** 1  
**Gap reference:** GAP-001  
**Owner:** —

**Goal:**  
The Docker container currently exits immediately because `docker/Dockerfile` builds the binary as `neuron-agent` but `docker/docker-entrypoint.sh` looks for `neuronagent`. Fix the mismatch so the container starts.

**Files to inspect:**
- `docker/Dockerfile` — look for the `go build -o` line
- `docker/docker-entrypoint.sh` — look for the binary path check and exec call

**Files to change:**
- `docker/docker-entrypoint.sh` — change `/app/neuronagent` to `/app/neuron-agent`

**Implementation notes:**  
Do not rename the Dockerfile output — changing the entrypoint is safer because it touches fewer downstream references. If there are other scripts or docs that reference `/app/neuronagent`, search for them and update together:
```bash
grep -r "neuronagent" docker/ scripts/ docs/
```

**Test plan:**
1. `docker build -f docker/Dockerfile -t test-fix .`
2. `docker run --rm test-fix --help` — must print help text, not "Binary not found"
3. `docker run --rm test-fix /app/neuron-agent --version` — must succeed

**Acceptance criteria:**
- [ ] `docker/docker-entrypoint.sh` references `/app/neuron-agent` (with hyphen)
- [ ] `docker build` succeeds
- [ ] `docker run --rm test-fix --help` exits with code 0 and prints help
- [ ] No other references to the old binary path remain in Docker-related files

**Dependencies:** None

---

### M1-004: Fix Makefile — replace docker-compose with docker compose

**Title:** Fix Makefile: replace deprecated docker-compose with docker compose  
**Priority:** P0  
**Effort:** XS (10 minutes)  
**Week:** 1  
**Gap reference:** GAP-026  
**Owner:** —

**Goal:**  
`docker-compose` (V1) is deprecated and removed from recent Docker installations. Replace all instances with `docker compose` (V2 subcommand).

**Files to inspect:**
- `Makefile` — search for `docker-compose`

**Files to change:**
- `Makefile` — replace all `docker-compose` with `docker compose`

**Implementation notes:**
```bash
sed -i 's/docker-compose/docker compose/g' Makefile
```
Verify the result: `grep "docker-compose" Makefile` should return nothing.

Also check `scripts/` for the same issue:
```bash
grep -r "docker-compose" scripts/
```

**Test plan:**
1. `make docker-up` (or equivalent target) — must not fail with "docker-compose: command not found"
2. `make docker-down` — must succeed

**Acceptance criteria:**
- [ ] No `docker-compose` (hyphenated) remains in `Makefile`
- [ ] No `docker-compose` remains in `scripts/` shell scripts

**Dependencies:** None

---

### M1-005: Fix Makefile — remove || true from integration-test

**Title:** Fix Makefile: integration-test must fail on test failures  
**Priority:** P0  
**Effort:** XS (5 minutes)  
**Week:** 1  
**Gap reference:** GAP-011  
**Owner:** —

**Goal:**  
The `integration-test` Makefile target silently passes even when tests fail because of `|| true`. Remove it so test failures surface correctly.

**Files to inspect:**
- `Makefile` — find the `integration-test` target

**Files to change:**
- `Makefile` — remove `|| true` from `integration-test` target

**Implementation notes:**
Current (broken):
```makefile
integration-test:
    cd src && go test -short ./... || true
```
Fixed:
```makefile
integration-test:
    cd src && go test -short ./...
```

**Test plan:**
1. Temporarily add a failing test: `func TestForceFail(t *testing.T) { t.Fatal("force fail") }`
2. Run `make integration-test` — must exit with non-zero code
3. Remove the temporary test

**Acceptance criteria:**
- [ ] `|| true` is removed from `integration-test` target
- [ ] Running `make integration-test` with a failing test exits with non-zero code
- [ ] `make integration-test` with all passing tests exits with code 0

**Dependencies:** None

---

### M1-016: Fix README Go version badge

**Title:** Fix README Go version badge from 1.23 to 1.24  
**Priority:** P0  
**Effort:** XS (5 minutes)  
**Week:** 1  
**Gap reference:** GAP-012  
**Owner:** —

**Goal:**  
README badge claims `Go 1.23+` but `src/go.mod` declares `go 1.24.0`. Fix to match.

**Files to inspect:**
- `README.md` — find the Go version badge
- `src/go.mod` — confirm version

**Files to change:**
- `README.md` — update badge URL and any text references

**Implementation notes:**  
Find the badge line: `![Go Version]` or similar. Change `1.23` to `1.24`. Also fix any prose references to Go version in prerequisites section.

**Acceptance criteria:**
- [ ] README badge shows `Go 1.24+`
- [ ] Prerequisites section says `Go 1.24+`
- [ ] Badge and prose are consistent

**Dependencies:** None

---

### M1-019: Convert docs .txt files to .md

**Title:** Convert documentation .txt files to Markdown  
**Priority:** P2  
**Effort:** S (2–4 hours)  
**Week:** 1  
**Gap reference:** GAP-034  
**Owner:** —

**Goal:**  
Files like `docs/deployment_guide.txt`, `docs/operations_runbook.txt`, and `docs/config_env_schema.txt` are not rendered by GitHub, not indexed as docs, and not searchable as Markdown. Convert them.

**Files to inspect:**
```
docs/deployment_guide.txt
docs/operations_runbook.txt
docs/config_env_schema.txt
docs/compliance_profiles.txt
docs/security_model.txt
docs/product_gap_report.txt
docs/product_readiness_audit.txt
docs/architecture_v2.txt
docs/neuronsql_design.txt
docs/neuronsql/quickstart.txt
docs/neuronsql/api.txt
docs/neuronsql/architecture.txt
docs/neuronsql/repo_map.txt
docs/neuronsql/security.txt
docs/neuronsql/eval.txt
```

**Files to change:**
- Rename each `.txt` to `.md` or create proper `.md` replacements
- Content that is superseded by new docs (from roadmap targets) can be archived or removed

**Implementation notes:**  
Do not blindly rename — read each file and decide:
- If the content is good and just needs formatting, rename to `.md` and add proper Markdown headers
- If the content is superseded by a roadmap doc, move relevant parts and delete
- If the content is outdated, archive to `docs/archive/` and note it in a commit message

**Acceptance criteria:**
- [ ] No `.txt` files remain in `docs/`
- [ ] All content is preserved or explicitly archived
- [ ] New `.md` files have proper Markdown formatting (headers, code blocks, lists)

**Dependencies:** None

---

## Week 2 Tickets

---

### M1-002: Create root-level docker-compose.yml

**Title:** Create root-level docker-compose.yml with NeuronDB-enabled PostgreSQL  
**Priority:** P0  
**Effort:** S (half day)  
**Week:** 2  
**Gap reference:** GAP-002, GAP-003  
**Owner:** —

**Goal:**  
Create `docker-compose.yml` at the repository root that starts NeuronAgent and a NeuronDB-enabled PostgreSQL together. This is the foundation for one-command local install.

**Files to inspect:**
- `docker/docker-compose.neuronsql.yml` — reference for service names and ports
- `docker/Dockerfile` — understand what environment variables the server needs
- `sql/neuron-agent.sql` — understand what extensions and schema are needed

**Files to change:**
- Create `docker-compose.yml` at repository root
- Create `docker/init-scripts/01-extensions.sh` (if using init-script approach for NeuronDB)

**Implementation notes:**

```yaml
services:
  db:
    image: postgres:17-alpine          # replace with neurondb/postgres:17 when available
    environment:
      POSTGRES_DB: ${DB_NAME:-neuronagent}
      POSTGRES_USER: ${DB_USER:-neuronagent}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./docker/init-scripts:/docker-entrypoint-initdb.d:ro
      - ./sql:/sql:ro
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-neuronagent}"]
      interval: 5s
      timeout: 5s
      retries: 12
    ports:
      - "5432:5432"

  neuronagent:
    build:
      context: .
      dockerfile: docker/Dockerfile
    image: neurondb/neuron-agent:local
    env_file: .env
    environment:
      DB_HOST: db
      DB_PORT: 5432
    ports:
      - "${SERVER_PORT:-8080}:8080"
    depends_on:
      db:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 6

volumes:
  postgres_data:
```

**NeuronDB note:**  
If `neurondb/postgres:17` is not yet published, use the init-script approach: create `docker/init-scripts/01-extensions.sh` that runs `psql -c "CREATE EXTENSION IF NOT EXISTS neurondb;"`. This requires the NeuronDB shared library to be present in the postgres image, which may not be possible without a custom image. Document this limitation clearly.

If NeuronDB extension is not available in any Docker image, the compose file must still work for development (without vector features) and must print a clear warning about what is missing.

**Test plan:**
1. `cp .env.example .env`
2. `docker compose up -d`
3. `docker compose ps` — both services show `healthy`
4. `curl http://localhost:8080/health` — returns `{"status":"ok"}`
5. `docker compose down -v` — removes cleanly

**Acceptance criteria:**
- [ ] `docker-compose.yml` exists at repository root
- [ ] `docker compose up -d` starts both services
- [ ] Both services reach healthy state within 60 seconds
- [ ] `curl localhost:8080/health` returns 200
- [ ] `docker compose down` stops cleanly
- [ ] `docker compose down -v` removes volumes

**Dependencies:** M1-001 (Docker entrypoint must be fixed first)

---

### M1-003: Create .env.example

**Title:** Create .env.example with all required variables and comments  
**Priority:** P0  
**Effort:** XS (1–2 hours)  
**Week:** 2  
**Gap reference:** GAP-004  
**Owner:** —

**Goal:**  
Create `.env.example` at repository root so any developer can start with `cp .env.example .env`.

**Files to inspect:**
- `src/internal/config/env.go` — authoritative list of all env vars
- `src/internal/config/defaults.go` — default values

**Files to change:**
- Create `.env.example` at repository root

**Implementation notes:**

```bash
# NeuronAgent Configuration
# Copy this file to .env and fill in the values.
# For local development, the defaults here are safe.
# For production, change ALL values marked CHANGE ME.

# ─── Database ────────────────────────────────────────────────────────────────
DB_HOST=localhost
DB_PORT=5432
DB_NAME=neuronagent
DB_USER=neuronagent
DB_PASSWORD=local-dev-password-CHANGE-ME-IN-PRODUCTION

# ─── Server ──────────────────────────────────────────────────────────────────
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
# Mode: local | development | test | production
# In production mode, weak secrets cause a startup error.
SERVER_MODE=local

# ─── Auth ────────────────────────────────────────────────────────────────────
# Secret used to sign API keys. Must be at least 32 characters in production.
AUTH_API_KEY_SECRET=local-dev-secret-CHANGE-ME-IN-PRODUCTION

# ─── Logging ─────────────────────────────────────────────────────────────────
# Level: debug | info | warn | error
LOG_LEVEL=info
# Format: json | text
LOG_FORMAT=json

# ─── NeuronSQL LLM Sidecar (optional) ────────────────────────────────────────
# Only needed if you use the NeuronSQL SQL generation feature.
# LLM_SQL_BASE_URL=http://localhost:11434
# LLM_SQL_API_KEY=your-api-key

# ─── CORS ────────────────────────────────────────────────────────────────────
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

**Test plan:**
1. `cp .env.example .env`
2. `docker compose up -d` with the new `.env` — services start without errors

**Acceptance criteria:**
- [ ] `.env.example` exists at repository root
- [ ] Every environment variable read in `src/internal/config/env.go` has an entry
- [ ] Each variable has an inline comment explaining what it does
- [ ] Production-sensitive variables are marked with "CHANGE ME IN PRODUCTION"
- [ ] File has no actual secrets (no real passwords, no real API keys)
- [ ] `.env` is in `.gitignore`

**Dependencies:** None

---

### M1-007: Create scripts/install.sh

**Title:** Create scripts/install.sh — one-command installer with prerequisite checks  
**Priority:** P1  
**Effort:** S (3–4 hours)  
**Week:** 2  
**Gap reference:** GAP-005  
**Owner:** —

**Goal:**  
Enable the target install command:
```bash
curl -fsSL https://raw.githubusercontent.com/neurondb/neuron-agent/main/scripts/install.sh | bash
```

**Files to inspect:**
- `scripts/neuronagent-setup.sh` — reference for existing setup logic
- `Makefile` — understand what targets exist

**Files to change:**
- Create `scripts/install.sh`

**Implementation notes:**

The script must:
1. Check Docker is installed: `command -v docker` — if not, print `Docker is required. Install it from https://docs.docker.com/get-docker/` and exit 1
2. Check Docker is running: `docker info > /dev/null 2>&1` — if not, print `Docker is not running. Start Docker Desktop and try again.` and exit 1
3. Check docker compose (not docker-compose): `docker compose version > /dev/null 2>&1`
4. Check git is installed
5. Clone repo if not already in a NeuronAgent directory
6. Copy `.env.example` to `.env` if `.env` does not exist
7. Print: "Starting NeuronAgent stack..."
8. Run `docker compose up -d`
9. Wait for health using `scripts/wait-for-health.sh`
10. Print success message with next steps

Every error must print:
- What went wrong
- Why it matters
- What to do about it

```bash
#!/usr/bin/env bash
set -euo pipefail

REPO_URL="https://github.com/neurondb/neuron-agent.git"
HEALTH_URL="http://localhost:8080/health"

# Color output
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
info()  { echo -e "${GREEN}[neuron-agent]${NC} $*"; }
warn()  { echo -e "${YELLOW}[neuron-agent]${NC} $*"; }
error() { echo -e "${RED}[neuron-agent]${NC} ERROR: $*"; exit 1; }

# Prerequisite checks
command -v docker >/dev/null 2>&1 || \
    error "Docker is not installed. Install it from https://docs.docker.com/get-docker/"
docker info >/dev/null 2>&1 || \
    error "Docker is not running. Start Docker and run this script again."
docker compose version >/dev/null 2>&1 || \
    error "docker compose (V2) is required. Update Docker to version 24 or later."
command -v git >/dev/null 2>&1 || \
    error "git is required. Install it from https://git-scm.com/"

# ... clone, cp .env.example .env, docker compose up -d ...
```

**Test plan:**
1. Run on a machine that has never had NeuronAgent installed
2. Verify each error message by temporarily removing prerequisites
3. Verify success path completes and prints next steps

**Acceptance criteria:**
- [ ] Script fails with a clear message if Docker is not installed
- [ ] Script fails with a clear message if Docker is not running
- [ ] Script fails with a clear message if git is not installed
- [ ] Script copies `.env.example` to `.env` only if `.env` does not exist (does not overwrite)
- [ ] Script runs `docker compose up -d`
- [ ] Script waits for health before declaring success
- [ ] Script prints next steps on success

**Dependencies:** M1-002 (compose file), M1-003 (.env.example), M1-008 (wait-for-health.sh)

---

### M1-008: Create scripts/wait-for-health.sh

**Title:** Create scripts/wait-for-health.sh — polls health endpoint until ready  
**Priority:** P1  
**Effort:** XS (30 minutes)  
**Week:** 2  
**Gap reference:** None (new capability)  
**Owner:** —

**Goal:**  
A reusable script that other scripts (install.sh, demo.sh, CI) can source to wait for NeuronAgent to be healthy before proceeding.

**Files to change:**
- Create `scripts/wait-for-health.sh`

**Implementation notes:**

```bash
#!/usr/bin/env bash
set -euo pipefail

HEALTH_URL="${1:-http://localhost:8080/health}"
MAX_ATTEMPTS="${2:-30}"
SLEEP_SECONDS=2

for i in $(seq 1 "$MAX_ATTEMPTS"); do
    if curl -sf "$HEALTH_URL" > /dev/null 2>&1; then
        echo "NeuronAgent is healthy."
        exit 0
    fi
    echo "Waiting for NeuronAgent... ($i/$MAX_ATTEMPTS)"
    sleep "$SLEEP_SECONDS"
done

echo "ERROR: NeuronAgent did not become healthy after $((MAX_ATTEMPTS * SLEEP_SECONDS)) seconds."
echo "Check logs with: docker compose logs neuronagent"
exit 1
```

**Acceptance criteria:**
- [ ] Script polls the health endpoint until it returns 200
- [ ] Script exits 0 when healthy
- [ ] Script exits 1 after timeout with a helpful message including the log command
- [ ] Timeout is configurable via argument
- [ ] Default timeout is 60 seconds (30 attempts × 2 seconds)

**Dependencies:** M1-001 (container must start)

---

### M1-011: Add Makefile targets: up, down, demo, smoke, reset

**Title:** Add missing Makefile targets for common developer workflows  
**Priority:** P1  
**Effort:** XS (30 minutes)  
**Week:** 2  
**Gap reference:** None  
**Owner:** —

**Goal:**  
`make up`, `make down`, `make demo`, `make smoke`, and `make reset` must work after this ticket.

**Files to inspect:**
- `Makefile` — current targets

**Files to change:**
- `Makefile`

**Implementation notes:**

```makefile
# Stack management
up:
    docker compose up -d

down:
    docker compose down

logs:
    docker compose logs -f

reset:
    docker compose down -v && docker compose up -d

# Development workflows
demo:
    @bash scripts/demo.sh

smoke:
    @bash scripts/smoke-test.sh

# Existing targets remain unchanged
```

**Acceptance criteria:**
- [ ] `make up` starts the compose stack
- [ ] `make down` stops the compose stack
- [ ] `make logs` tails compose logs
- [ ] `make reset` stops, removes volumes, and restarts (fresh state)
- [ ] `make demo` runs `scripts/demo.sh` (script need not exist yet — prints helpful error if missing)
- [ ] `make smoke` runs smoke test (script need not exist yet — prints helpful error if missing)

**Dependencies:** M1-002, M1-004

---

### M1-018: Resolve migration naming inconsistency

**Title:** Resolve migration naming inconsistency between psql script and Go migration system  
**Priority:** P2  
**Effort:** M (half day)  
**Week:** 2  
**Gap reference:** GAP-023  
**Owner:** —

**Goal:**  
The Go migration system (`src/internal/db/schema.go`) expects files named `001_name.sql`. The main schema file is `neuron-agent.sql`. Align them.

**Files to inspect:**
- `sql/` — list all SQL files and their names
- `src/internal/db/schema.go` — understand version parsing logic
- `src/internal/db/migrations.go` — understand migration runner
- `scripts/neuronagent-migrate.sh` — understand shell migration approach

**Files to change:**
- Rename SQL files in `sql/` to match `001_name.sql` pattern, OR
- Update `schema.go` to handle the existing naming, OR
- Document that the shell script is authoritative and remove the Go migration runner

**Implementation notes:**

Recommended approach:
1. Rename `sql/neuron-agent.sql` to `sql/001_initial_schema.sql`
2. Rename other schema files to `002_partitioning.sql`, `003_rls.sql`, etc. based on logical dependency order
3. Update `scripts/neuronagent-migrate.sh` to iterate in numeric order: `ls sql/*.sql | sort`
4. Verify the Go migration runner can now parse version numbers correctly
5. Update the compose `init-scripts` to use the numbered approach

**Test plan:**
1. Start a fresh PostgreSQL container
2. Run the migration script
3. Connect and verify all tables exist: `\dt neurondb_agent.*`
4. Run migrations again — must be idempotent (no errors on re-run)

**Acceptance criteria:**
- [ ] All SQL files in `sql/` follow `NNN_description.sql` naming
- [ ] Migration order is deterministic (no glob ambiguity)
- [ ] Running migrations on a fresh DB creates all required tables
- [ ] Running migrations again on an existing DB does not error
- [ ] Go migration runner and shell script use the same approach

**Dependencies:** None

---

## Week 3 Tickets

---

### M1-015: Rewrite README to adoption structure

**Title:** Rewrite README to lead with install and demo, not architecture  
**Priority:** P1  
**Effort:** M (half day to full day)  
**Week:** 3  
**Gap reference:** GAP-012  
**Owner:** —

**Goal:**  
A visitor who has never heard of NeuronAgent reads the README and within 90 seconds understands what it is and wants to try it.

**Files to inspect:**
- `README.md` — full current contents
- `docs/architecture.md` — content to move deep architecture to
- `docs/api.md` or `docs/api-reference.md` — content to move API details to

**Files to change:**
- `README.md` — complete rewrite following the target structure
- `docs/architecture.md` — add content moved from README

**Target README structure:**
```
# NeuronAgent

[Badge: CI status — only if CI workflow exists on push trigger]
[Badge: Go 1.24 — correct version]
[Badge: License]

One-sentence description.

## Install

[4 commands — git clone, cp .env.example .env, docker compose up -d, make demo]

## What you get

[6-8 bullet points of concrete capabilities — no architecture names]

## Why NeuronAgent

[3-4 sentences on positioning vs alternatives]

## Use cases

[4-6 one-line scenarios]

## Examples

[Links to examples/ with one-line description each]

## OpenClaw bridge

[One paragraph, link to docs]

## Production deployment

[Link to docs/deploy/]

## Docs

[docs/quickstart.md, docs/api.md, docs/concepts.md]

## Contributing

[Link to CONTRIBUTING.md]

## Security

[Link to SECURITY.md or security section]

## License

[License statement]
```

**Rules (all must be followed):**
- No architecture diagrams in README
- No internal package names (no `src/internal/agent/`)
- No fake badges (remove CI badge if CI is workflow_dispatch only)
- No feature claims without backing code
- Every command in README is tested before PR merges
- README must not require scrolling more than 3 screens to see the install section

**Test plan:**
1. Find a person who has never seen NeuronAgent
2. Give them the README
3. Ask: "What is this?" — must be answerable in one sentence
4. Ask: "How do you install it?" — must be answerable from README
5. Run every command in the README on a clean machine — all must succeed

**Acceptance criteria:**
- [ ] README opens with a project statement (1–2 sentences)
- [ ] Install section appears before any architecture or feature explanation
- [ ] Install uses 4 or fewer commands
- [ ] README has no fake badges
- [ ] README has no dead links
- [ ] README Go version badge says `1.24`
- [ ] README has no commands that fail on a clean machine
- [ ] Deep architecture is in `docs/architecture.md`, not README
- [ ] Full API reference is in `docs/api.md`, not README

**Dependencies:** M1-002, M1-003, M1-007

---

## Week 4 Tickets

---

### M1-009: Create scripts/demo.sh

**Title:** Create scripts/demo.sh — golden path end-to-end demo  
**Priority:** P1  
**Effort:** M (half day)  
**Week:** 4  
**Gap reference:** GAP-006  
**Owner:** —

**Goal:**  
`make demo` must run a complete, clean demo that covers every major NeuronAgent capability and prints readable output.

**Files to inspect:**
- `src/cmd/generate-key/main.go` — understand key generation
- `docs/api.md` or `docs/api-reference.md` — understand API request shapes
- Any existing demo scripts in `scripts/`

**Files to change:**
- Create `scripts/demo.sh`

**Implementation notes:**

The demo script must:
1. Check that the stack is healthy (`scripts/wait-for-health.sh`)
2. Generate or load an API key (use `docker compose exec neuronagent /app/generate-key` or call admin API)
3. Create an agent:
   ```bash
   curl -sf -X POST http://localhost:8080/api/v1/agents \
     -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"name":"demo-agent","system_prompt":"You are a helpful assistant."}'
   ```
4. Ingest a fixture document:
   ```bash
   curl -sf -X POST http://localhost:8080/api/v1/rag/ingest \
     -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d "{\"title\":\"NeuronAgent Overview\",\"content\":$(cat tests/fixtures/sample-doc.txt | jq -Rs .)}"
   ```
5. Create a session and send a message
6. Print the response
7. Run a read-only SQL query via the SQL tool
8. Print pass/fail for each step

Output format — each step should print:
```
[1/6] Creating agent...     ✓
[2/6] Ingesting document... ✓
[3/6] Creating session...   ✓
[4/6] Asking question...    ✓
[5/6] Running SQL query...  ✓
[6/6] Demo complete!        ✓
```

If any step fails, print:
```
[3/6] Creating session...   ✗
  Error: HTTP 401 - API key required
  Check: Is your .env file correct? Run: cat .env | grep AUTH
```

**Test plan:**
1. Start a clean stack with `make reset && make up`
2. Run `make demo` — must complete all steps with ✓
3. Run `make demo` again — must work (idempotent)
4. Stop the stack, run `make demo` — must fail with clear error about stack not running

**Acceptance criteria:**
- [ ] All 6 demo steps complete successfully on a clean stack
- [ ] Demo prints pass/fail for each step
- [ ] Demo exits 0 on success, non-zero on any failure
- [ ] Demo is idempotent (running twice does not fail)
- [ ] Demo fails with helpful error if stack is not running

**Dependencies:** M1-001, M1-002, M1-003, M1-008, M1-017

---

### M1-010: Create scripts/bootstrap-demo.sh

**Title:** Create scripts/bootstrap-demo.sh — create demo workspace and API key  
**Priority:** P1  
**Effort:** S (2–3 hours)  
**Week:** 4  
**Gap reference:** None  
**Owner:** —

**Goal:**  
Create a script that initializes NeuronAgent with a demo workspace, API key, and basic agent on first run. This is called by `demo.sh` and also useful for CI.

**Files to inspect:**
- `src/cmd/generate-key/main.go` — understand key generation binary
- `scripts/neuronagent-generate-keys.sh` — reference existing key generation

**Files to change:**
- Create `scripts/bootstrap-demo.sh`

**Implementation notes:**  
The script should:
1. Check if bootstrap has already been done (presence of `.bootstrap-complete` file)
2. If not done: generate an API key, save it to `.demo.env`
3. Create a default demo agent via the API
4. Print the API key (once only)
5. Mark bootstrap as complete

The key output should be:
```
Bootstrap complete!
API Key: na_xxxxxxxxxxxxxxxxxxxx
Saved to: .demo.env

To use in scripts: source .demo.env && echo $API_KEY
```

**Acceptance criteria:**
- [ ] Bootstrap script creates an API key
- [ ] API key is saved to `.demo.env` (not `.env` — don't overwrite user config)
- [ ] Running bootstrap twice is idempotent (does not create duplicate keys)
- [ ] Key is usable for subsequent API calls

**Dependencies:** M1-001, M1-002, M1-003

---

### M1-012: Create examples/quickstart-chat

**Title:** Create examples/quickstart-chat with README and working demo script  
**Priority:** P1  
**Effort:** S (3–4 hours)  
**Week:** 4  
**Gap reference:** None  
**Owner:** —

**Goal:**  
The simplest possible working example: create an agent, send a message, see the response.

**Files to change:**
- Create `examples/quickstart-chat/README.md`
- Create `examples/quickstart-chat/demo.sh`
- Create `examples/quickstart-chat/agent-config.json`

**Implementation notes:**

`README.md` must cover:
- What this example shows (2 sentences)
- Prerequisites (running stack, API key)
- Run instructions (2–3 commands)
- Expected output
- Cleanup

`demo.sh` must:
1. Create an agent (using config from `agent-config.json`)
2. Create a session
3. Send 2–3 messages
4. Print the responses
5. Clean up the agent

`agent-config.json`:
```json
{
  "name": "quickstart-agent",
  "system_prompt": "You are a helpful assistant. Be concise."
}
```

**Acceptance criteria:**
- [ ] `cd examples/quickstart-chat && bash demo.sh` works on a running stack
- [ ] README has correct prerequisites
- [ ] README has accurate expected output
- [ ] Demo creates, uses, and deletes the agent
- [ ] Demo exits 0 on success

**Dependencies:** M1-001, M1-002, M1-003

---

### M1-013: Create examples/rag-with-docs

**Title:** Create examples/rag-with-docs with document ingest and RAG query demo  
**Priority:** P1  
**Effort:** M (half day)  
**Week:** 4  
**Gap reference:** None  
**Owner:** —

**Goal:**  
Demonstrate document ingest and RAG question answering with a concrete, verifiable example.

**Files to change:**
- Create `examples/rag-with-docs/README.md`
- Create `examples/rag-with-docs/demo.sh`
- Create `examples/rag-with-docs/sample-document.md`

**Implementation notes:**

`sample-document.md` — use a small, factual document with specific information that can be verified in the demo output:
```markdown
# NeuronAgent Technical Overview

NeuronAgent version 3.0 was released in 2026. It uses PostgreSQL as its
primary database and NeuronDB as its AI substrate. The default HTTP port
is 8080. The health endpoint is available at /health.
```

`demo.sh` must:
1. Ingest `sample-document.md`
2. Ask a question whose answer is in the document: "What port does NeuronAgent use by default?"
3. Verify the answer contains "8080"
4. Print pass/fail based on verification

This makes the demo self-validating — the expected answer is known.

**Acceptance criteria:**
- [ ] Demo ingests a document successfully
- [ ] Demo asks a question and receives an answer containing the expected value
- [ ] Demo script verifies the answer (not just prints it)
- [ ] Demo exits 0 on success, 1 on failure

**Dependencies:** M1-001, M1-002, M1-003

---

### M1-014: Create examples/sql-agent

**Title:** Create examples/sql-agent with read-only SQL inspection demo  
**Priority:** P1  
**Effort:** S (3 hours)  
**Week:** 4  
**Gap reference:** None  
**Owner:** —

**Goal:**  
Demonstrate SQL schema inspection and read-only query execution through NeuronAgent.

**Files to change:**
- Create `examples/sql-agent/README.md`
- Create `examples/sql-agent/demo.sh`
- Create `examples/sql-agent/agent-config.json`

**Implementation notes:**

`demo.sh` must:
1. Create an agent with NeuronSQL tools enabled
2. Ask the agent to describe the `neurondb_agent` schema
3. Ask the agent to count the number of agents in the agents table
4. Print the responses
5. Verify no DDL or DML was executed

Use `neuronsql.schema_snapshot` tool to inspect the schema.

**Acceptance criteria:**
- [ ] Demo creates an agent with SQL tools
- [ ] Agent returns schema information
- [ ] Demo exits 0 on success

**Dependencies:** M1-001, M1-002, M1-003

---

### M1-017: Create tests/fixtures/sample-doc.txt

**Title:** Create test fixture document for demo and smoke test  
**Priority:** P1  
**Effort:** XS (15 minutes)  
**Week:** 4  
**Gap reference:** None  
**Owner:** —

**Goal:**  
Create a small, factual fixture document used by `demo.sh`, smoke test, and RAG examples.

**Files to change:**
- Create `tests/fixtures/sample-doc.txt`

**Content requirements:**
- Contains 3–5 specific facts about NeuronAgent that can be queried
- Short enough to be processed quickly (< 500 words)
- Facts must be verifiable (specific version numbers, port numbers, endpoint names)
- Not copyrighted material

**Acceptance criteria:**
- [ ] File exists at `tests/fixtures/sample-doc.txt`
- [ ] Contains at least 3 verifiable facts
- [ ] File is under 500 words

**Dependencies:** None

---

### M1-006: Add CI workflow on push and pull_request

**Title:** Add GitHub Actions CI with push and PR triggers  
**Priority:** P0  
**Effort:** S (3–4 hours)  
**Week:** 4  
**Gap reference:** GAP-007  
**Owner:** —

**Goal:**  
Every PR must run lint and unit tests automatically. Broken code cannot merge.

**Files to inspect:**
- `.github/workflows/neuron-agent-build-matrix.yml` — existing workflow to extend or replace

**Files to change:**
- Create `.github/workflows/ci.yml` (new, focused CI workflow)
- Keep existing build-matrix workflow for manual full-matrix testing

**Implementation notes:**

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          working-directory: src

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Run tests
        run: cd src && go test -short ./...

  docker-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build Docker image
        run: docker build -f docker/Dockerfile -t neuron-agent:ci .
      - name: Smoke test
        run: |
          cp .env.example .env
          docker compose up -d
          bash scripts/wait-for-health.sh
          curl -sf http://localhost:8080/health
          docker compose down -v
```

**Test plan:**
1. Open a PR to main — CI must trigger automatically
2. Introduce a failing test — CI must fail and block merge
3. Fix the test — CI must pass

**Acceptance criteria:**
- [ ] CI triggers on every push to `main`
- [ ] CI triggers on every PR to `main`
- [ ] CI includes: lint, unit tests, Docker build
- [ ] CI smoke test calls `GET /health` and verifies 200
- [ ] Failing test causes CI to fail (exit non-zero)
- [ ] CI completes in under 10 minutes for the standard path

**Dependencies:** M1-001, M1-002, M1-003, M1-008

---

### M1-020: Add smoke test to CI

**Title:** Add automated smoke test to CI pipeline  
**Priority:** P1  
**Effort:** M (part of M1-006)  
**Week:** 4  
**Gap reference:** None  
**Owner:** —

**Goal:**  
CI must verify that the golden path (start stack → create agent → send message → verify response) works on every PR.

**Files to change:**
- Create `scripts/smoke-test.sh` — lightweight golden path verification
- Update `.github/workflows/ci.yml` to run smoke test

**Implementation notes:**

`scripts/smoke-test.sh` must:
1. Call `GET /health` — verify 200
2. Call `POST /api/v1/agents` — verify 201, capture agent ID
3. Call `POST /api/v1/sessions` — verify 201, capture session ID
4. Call `POST /api/v1/sessions/{id}/messages` — verify response is received
5. Verify response has the expected shape: `{"id": "...", "content": "...", "role": "assistant"}`
6. Delete the test agent
7. Exit 0 on success, 1 on any failure

The smoke test uses a test API key that must be created during stack bootstrap. The key value is read from `TEST_API_KEY` env var.

**Acceptance criteria:**
- [ ] `make smoke` runs the smoke test against the local stack
- [ ] Smoke test covers: health, create agent, create session, send message, verify response shape
- [ ] Smoke test exits 0 on success, 1 on failure
- [ ] Smoke test is in CI pipeline
- [ ] Smoke test completes in under 60 seconds

**Dependencies:** M1-006 (CI workflow), M1-009 (or subset of demo.sh)

---

## Month 1 Acceptance Criteria

At the end of Month 1, all of the following must be true:

- [ ] **GAP-001 resolved:** Docker container starts without "Binary not found"
- [ ] **GAP-002 resolved:** `docker-compose.yml` includes NeuronDB-enabled PostgreSQL
- [ ] **GAP-003 resolved:** `docker-compose.yml` exists at repository root
- [ ] **GAP-004 resolved:** `.env.example` exists and is accurate
- [ ] **GAP-005 resolved:** `scripts/install.sh` exists and works
- [ ] **GAP-006 resolved:** `scripts/demo.sh` runs the golden path end-to-end
- [ ] **GAP-007 resolved:** CI runs on push and PR
- [ ] **GAP-011 resolved:** `integration-test` fails on test failures
- [ ] **GAP-012 partial:** README Go version badge is correct
- [ ] **GAP-026 resolved:** Makefile uses `docker compose` not `docker-compose`
- [ ] **Install time:** Under 5 minutes from `git clone` to `make demo` success
- [ ] **Demo success:** Demo runs end-to-end on a clean machine with only Docker
- [ ] **Health check:** `GET /health` returns 200 on a running stack
- [ ] **Smoke test:** `make smoke` passes and is in CI
- [ ] **Examples:** `examples/quickstart-chat/`, `examples/rag-with-docs/`, `examples/sql-agent/` all work
- [ ] **README:** Explains the project clearly, no fake badges, no broken commands

---

## Effort Summary

| Ticket | Effort |
|---|---|
| M1-001 | XS |
| M1-002 | S |
| M1-003 | XS |
| M1-004 | XS |
| M1-005 | XS |
| M1-006 | S |
| M1-007 | S |
| M1-008 | XS |
| M1-009 | M |
| M1-010 | S |
| M1-011 | XS |
| M1-012 | S |
| M1-013 | M |
| M1-014 | S |
| M1-015 | M |
| M1-016 | XS |
| M1-017 | XS |
| M1-018 | M |
| M1-019 | S |
| M1-020 | M |
| **Total** | ~8 person-days |

XS = < 1 hour | S = 2–4 hours | M = half day | L = full day

Month 1 is achievable in 2 weeks of focused work for one engineer, or 1 week with two engineers working in parallel on independent tickets.

The critical path is: M1-001 → M1-002 → M1-003 → M1-007 → M1-009 → M1-006/M1-020

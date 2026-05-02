# NeuronAgent Release Readiness Checklist

> Last updated: 2026-04-30  
> Use this checklist before every versioned release. Check every box. Do not skip items.  
> For v1.0 release, all items must pass. For pre-releases, mark skipped items explicitly.

---

## How to use this checklist

Before cutting a release tag:
1. Create a release branch: `git checkout -b release/vX.Y.Z`
2. Work through every section below
3. Mark each item `[x]` when verified
4. If an item cannot pass, document why in the notes field
5. Create a PR for the release branch and get a second set of eyes on the checklist
6. Only tag after the PR is approved and all checks pass

---

## Section 1: Install and Quickstart

### 1.1 Clean machine test

- [ ] Tested on a machine that has never had NeuronAgent installed
- [ ] Only Docker and git were present before starting
- [ ] Ran every command in the README quickstart exactly as written
- [ ] All commands completed without errors
- [ ] Time from `git clone` to `make demo` success: _____ minutes (must be under 5)

### 1.2 Install script

- [ ] `scripts/install.sh` exists and is executable
- [ ] Script checks Docker version (requires 24+)
- [ ] Script checks docker compose (not docker-compose)
- [ ] Script fails with a clear error if Docker is not installed
- [ ] Script fails with a clear error if Docker is not running
- [ ] Script prints success message when complete
- [ ] Script can be run with `curl -fsSL ... | bash` from a fresh system

### 1.3 Docker compose

- [ ] `docker-compose.yml` exists at repository root
- [ ] `docker compose up -d` starts both NeuronAgent and a NeuronDB-enabled PostgreSQL
- [ ] Both services pass their health checks within 60 seconds
- [ ] NeuronAgent container does not exit immediately (GAP-001 is fixed)
- [ ] `docker compose down` stops all services cleanly
- [ ] `docker compose down -v` removes all volumes and leaves a clean state
- [ ] Named volumes are used (not anonymous volumes)

### 1.4 Environment setup

- [ ] `.env.example` exists at repository root
- [ ] Every required environment variable is listed in `.env.example`
- [ ] Each variable has a comment explaining what it does
- [ ] Copying `.env.example` to `.env` produces a working local configuration
- [ ] No secrets are committed to the repository

### 1.5 Demo

- [ ] `scripts/demo.sh` exists and is executable
- [ ] `make demo` runs `scripts/demo.sh` and succeeds
- [ ] Demo creates an agent
- [ ] Demo ingests a fixture document
- [ ] Demo asks a question and receives an answer
- [ ] Demo runs a read-only SQL query
- [ ] Demo prints clear pass/fail output at each step
- [ ] Demo completes in under 2 minutes

---

## Section 2: API Correctness

### 2.1 Core endpoints

- [ ] `GET /health` returns 200 with `{"status": "ok"}`
- [ ] `GET /metrics` returns 200 with Prometheus text format
- [ ] `GET /version` returns version JSON (must exist before v1.0)
- [ ] `GET /healthz` returns 200 when server is healthy
- [ ] `GET /readyz` returns 200 when DB is connected and migrations are applied
- [ ] `GET /readyz` returns 503 when DB is not connected

### 2.2 Agent API

- [ ] `POST /api/v1/agents` creates an agent with required fields
- [ ] `POST /api/v1/agents` returns 400 for missing required fields
- [ ] `GET /api/v1/agents/{id}` returns 404 for unknown ID
- [ ] `DELETE /api/v1/agents/{id}` returns 404 for unknown ID
- [ ] All agent endpoints return `{"error": "...", "code": "...", "request_id": "..."}` on error

### 2.3 Session and message API

- [ ] `POST /api/v1/sessions/{id}/messages` returns a response (may be async or sync)
- [ ] `GET /api/v1/sessions` returns a list (empty list if none, not 404)
- [ ] Session messages are persisted across restarts

### 2.4 Memory API

- [ ] `POST /api/v1/agents/{id}/memory/search` returns results for a known query
- [ ] `POST /api/v1/agents/{id}/memory/forget` removes the specified memory chunk
- [ ] Memory persists across container restarts

### 2.5 Auth

- [ ] Requests without Authorization header return 401
- [ ] Requests with invalid API key return 401
- [ ] Admin endpoints return 403 for non-admin keys
- [ ] Rate limiting returns 429 after threshold is exceeded

### 2.6 OpenClaw bridge

- [ ] `GET /claw/v1/health` returns 200
- [ ] `POST /claw/v1/tools/list` returns list of available tools
- [ ] `POST /claw/v1/tools/run` executes a tool and returns result

---

## Section 3: Tests

### 3.1 Unit tests

- [ ] `make test` passes with exit code 0
- [ ] No tests are skipped (unless explicitly marked with reason)
- [ ] `make test` runs in under 5 minutes
- [ ] Test output shows coverage percentage for core packages

### 3.2 Integration tests

- [ ] `make integration-test` passes (must not use `|| true`)
- [ ] Integration tests run against a real database (not mocks only)
- [ ] Integration tests cover agent create → session → message → memory path

### 3.3 Smoke test

- [ ] `make smoke` passes against the local compose stack
- [ ] Smoke test covers: health, create agent, send message, verify response
- [ ] Smoke test is in CI and runs on every PR

### 3.4 Test coverage

- [ ] `src/internal/agent/` coverage is above 60%
- [ ] `src/internal/api/` has at least one test per handler group
- [ ] `src/internal/tools/` has tests for success, error, and timeout paths
- [ ] `src/internal/workflow/` has tests for each step type

---

## Section 4: Documentation

### 4.1 README

- [ ] README explains what NeuronAgent is in under 3 sentences
- [ ] README shows install commands (4 or fewer)
- [ ] README shows a demo or links to a demo
- [ ] README has no broken links
- [ ] README has no commands that don't work
- [ ] README Go version badge matches `src/go.mod`
- [ ] README has no badges for workflows that don't exist

### 4.2 Required docs pages

- [ ] `docs/quickstart.md` exists and works
- [ ] `docs/configuration.md` exists with all environment variables
- [ ] `docs/api.md` exists with curl examples
- [ ] `docs/memory.md` exists
- [ ] `docs/security.md` exists
- [ ] `docs/deploy/docker-compose.md` exists
- [ ] `docs/observability.md` exists

### 4.3 Docs quality

- [ ] All internal links in docs resolve (run `markdown-link-check` or equivalent)
- [ ] No `.txt` files in the `docs/` directory
- [ ] No doc page claims a feature that doesn't work
- [ ] No doc page is empty or contains only placeholder text

### 4.4 CHANGELOG

- [ ] `CHANGELOG.md` has an entry for this release version
- [ ] CHANGELOG entry lists changes, bug fixes, and breaking changes
- [ ] Breaking changes are clearly marked

---

## Section 5: Security

### 5.1 Authentication

- [ ] Auth middleware is applied to all routes except `/health`, `/metrics`, `/healthz`, `/readyz`
- [ ] No API key is stored in logs or error messages
- [ ] No secret values appear in the version endpoint response
- [ ] Production mode refuses to start with `AUTH_API_KEY_SECRET` shorter than 32 characters

### 5.2 Production config validation

- [ ] `SERVER_MODE=production` triggers strict config validation
- [ ] Server refuses to start in production mode with default DB password
- [ ] Server refuses to start in production mode with empty API key secret
- [ ] Startup prints config summary without any secret values

### 5.3 Security headers

- [ ] `X-Content-Type-Options: nosniff` is set on all responses
- [ ] `X-Frame-Options: DENY` is set on all responses
- [ ] `X-Request-ID` is set on all responses
- [ ] CORS origins are not `*` by default

### 5.4 Vulnerability scan

- [ ] Docker image has been scanned with Trivy (or equivalent)
- [ ] No HIGH or CRITICAL CVEs in the image (or all are documented with justification)
- [ ] Go dependencies scanned with `govulncheck` or equivalent
- [ ] No known vulnerabilities in direct dependencies

---

## Section 6: Docker and Packaging

### 6.1 Docker image

- [ ] `docker build -f docker/Dockerfile -t neuron-agent:test .` succeeds
- [ ] Built container starts without exiting immediately (GAP-001 confirmed fixed)
- [ ] `docker run --rm neuron-agent:test --help` returns help text
- [ ] Image size is documented (target: under 200 MB for production image)
- [ ] Image uses a non-root user for the application process

### 6.2 Multi-arch build

- [ ] `linux/amd64` image builds and starts
- [ ] `linux/arm64` image builds and starts (required before v1.0)

### 6.3 Published images

- [ ] Image is published to Docker Hub: `neurondb/neuron-agent:VERSION`
- [ ] Image is published to Docker Hub: `neurondb/neuron-agent:latest`
- [ ] Image is published to GHCR: `ghcr.io/neurondb/neuron-agent:VERSION`
- [ ] Image is published to GHCR: `ghcr.io/neurondb/neuron-agent:latest`
- [ ] Image pull and run works without authentication (public images)

---

## Section 7: CI/CD

### 7.1 CI triggers

- [ ] CI runs automatically on push to `main`
- [ ] CI runs automatically on every pull request
- [ ] CI failure blocks PR merge (branch protection enabled on `main`)

### 7.2 CI jobs

- [ ] `golangci-lint` runs and passes in CI
- [ ] `go test ./...` runs and passes in CI
- [ ] Docker build runs in CI (at least smoke-build)
- [ ] No `|| true` silencing test failures anywhere in CI or Makefile

### 7.3 Release workflow

- [ ] Pushing a `v*` tag triggers the release workflow
- [ ] Release workflow builds Docker images
- [ ] Release workflow publishes Docker images
- [ ] Release workflow creates a GitHub Release
- [ ] Release notes are generated automatically from commit history or CHANGELOG

---

## Section 8: Observability

### 8.1 Health endpoints

- [ ] `GET /health` returns correct status
- [ ] `GET /healthz` returns correct status
- [ ] `GET /readyz` returns 503 before DB is ready, 200 after

### 8.2 Metrics

- [ ] `GET /metrics` returns Prometheus metrics in correct format
- [ ] At minimum these metrics are present:
  - `neuronagent_http_requests_total`
  - `neuronagent_http_request_duration_seconds`
  - `neuronagent_agent_requests_total`
  - `neuronagent_tool_executions_total`
  - `neuronagent_memory_writes_total`

### 8.3 Logging

- [ ] All log output is valid JSON when `LOG_FORMAT=json`
- [ ] No secrets appear in log output
- [ ] Error logs include request ID
- [ ] Startup logs show version and config summary (no secrets)

---

## Section 9: Examples and Templates

### 9.1 Examples

- [ ] `examples/quickstart-chat/` has a working demo script
- [ ] `examples/rag-with-docs/` has a working demo script
- [ ] `examples/sql-agent/` has a working demo script
- [ ] All examples have a README
- [ ] All example scripts are tested on a clean stack

### 9.2 Templates (required for v1.0)

- [ ] At least 3 agent templates exist in `templates/`
- [ ] `make load-template TEMPLATE=<name>` works for each template
- [ ] Each template has documentation

---

## Section 10: v1.0 Specific Requirements

These items are only required for the v1.0 release. Pre-releases may skip them with documentation.

- [ ] All P0 and P1 gaps from `docs/roadmap/gap-analysis.md` are resolved
- [ ] API versioning is stable — no breaking changes to `/api/v1` without a new version
- [ ] All 5 required examples work
- [ ] At least 5 agent templates exist
- [ ] `docs/benchmarks.md` exists with real measured numbers
- [ ] Kubernetes Helm chart is complete and documented
- [ ] SBOM is generated and attached to the release
- [ ] Image is signed with cosign
- [ ] All items in `docs/roadmap/github-adoption-checklist.md` pass
- [ ] All items in `docs/roadmap/enterprise-readiness-checklist.md` pass at baseline level

---

## Release Sign-off

```
Release version: ___________
Release date: ___________
Prepared by: ___________
Reviewed by: ___________

Section 1 (Install): PASS / FAIL / PARTIAL
Section 2 (API): PASS / FAIL / PARTIAL
Section 3 (Tests): PASS / FAIL / PARTIAL
Section 4 (Docs): PASS / FAIL / PARTIAL
Section 5 (Security): PASS / FAIL / PARTIAL
Section 6 (Docker): PASS / FAIL / PARTIAL
Section 7 (CI/CD): PASS / FAIL / PARTIAL
Section 8 (Observability): PASS / FAIL / PARTIAL
Section 9 (Examples): PASS / FAIL / PARTIAL
Section 10 (v1.0): PASS / FAIL / N/A

Overall: APPROVED / NEEDS WORK

Notes:
```

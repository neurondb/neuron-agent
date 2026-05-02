# NeuronAgent Enterprise Readiness Checklist

> Last updated: 2026-04-30  
> Purpose: Track progress toward enterprise evaluation readiness  
> Audience: Engineering teams preparing NeuronAgent for enterprise customer evaluation

Enterprise evaluators ask hard questions. They do not accept "planned" or "coming soon." Every item below must have working code, tested behavior, or documented limitation with a tracked issue. Nothing is faked.

---

## How to use this checklist

Before an enterprise evaluation or demo:
1. Work through every section
2. Mark `[x]` for items that work and are tested
3. Mark `[-]` for items that are not applicable to this deployment
4. For every unchecked item, either fix it or document why it doesn't apply
5. Never present this checklist with unchecked items unless you have explained each one

---

## Section 1: Authentication and Access Control

### 1.1 Authentication methods

- [ ] API key authentication is implemented and required on all endpoints
- [ ] API keys can be rotated without downtime
- [ ] API keys can be revoked immediately
- [ ] Revoked keys are rejected within 1 request (no cache delay beyond 1 second)
- [ ] JWT authentication is supported for service-to-service calls
- [ ] Multiple authentication methods can coexist

### 1.2 Role-based access control (RBAC)

- [ ] Role model is defined and documented:
  - [ ] `owner` — full control including workspace deletion
  - [ ] `admin` — full control within workspace, no deletion
  - [ ] `developer` — create and run agents, read audit logs
  - [ ] `viewer` — read-only access to agents, sessions, memory
  - [ ] `service` — tool execution and session creation (for integration accounts)
- [ ] Roles are enforced at every API endpoint, not just admin routes
- [ ] Role assignment is documented in `docs/security.md`
- [ ] Permission denied responses return 403 (not 401 or 500)
- [ ] Audit log records the role of every actor

### 1.3 Workspace isolation

- [ ] Agents in workspace A cannot read sessions or memory from workspace B
- [ ] API keys scoped to workspace A cannot access workspace B resources
- [ ] Workspace boundaries are enforced at the database query level (not just application logic)
- [ ] Cross-workspace access attempts are logged as security events

### 1.4 Key management

- [ ] API key creation is documented
- [ ] API key rotation procedure is documented
- [ ] API key compromise response procedure is documented in `SECURITY.md`
- [ ] Keys are stored as hashes, not plaintext, in the database
- [ ] Keys are never logged, never returned in error messages, never in API responses

---

## Section 2: Audit Logging

### 2.1 Audit log schema

- [ ] Durable audit log table exists in `neurondb_agent` schema
- [ ] Audit log schema includes all required fields:
  - [ ] `id` — unique event ID
  - [ ] `timestamp` — UTC timestamp with microsecond precision
  - [ ] `actor_id` — API key ID or user ID
  - [ ] `actor_type` — `api_key`, `user`, `system`
  - [ ] `workspace_id` — workspace scope
  - [ ] `action` — event type (see event coverage below)
  - [ ] `resource_type` — `agent`, `session`, `memory`, `tool`, `workflow`, etc.
  - [ ] `resource_id` — UUID of the affected resource
  - [ ] `request_id` — correlates to HTTP request ID
  - [ ] `status` — `success`, `failure`, `denied`
  - [ ] `reason` — human-readable reason for failure or denial
  - [ ] `metadata` — JSONB for event-specific fields
  - [ ] `ip_address` — client IP (if applicable)

### 2.2 Audit event coverage

Every action below must produce an audit entry:

- [ ] API key used successfully
- [ ] API key authentication failure
- [ ] Rate limit exceeded
- [ ] Permission denied
- [ ] Agent created
- [ ] Agent updated
- [ ] Agent deleted
- [ ] Session created
- [ ] Message sent
- [ ] Tool executed (success)
- [ ] Tool executed (failure)
- [ ] Tool execution denied (permission)
- [ ] Memory chunk written
- [ ] Memory chunk deleted
- [ ] Memory search performed
- [ ] Workflow started
- [ ] Workflow step executed
- [ ] Workflow step failed
- [ ] Workflow approval requested
- [ ] Workflow approved
- [ ] Workflow rejected
- [ ] Workflow completed
- [ ] Admin config read
- [ ] Admin diagnostics read

### 2.3 Audit log security

- [ ] Audit log entries cannot be deleted via the standard API
- [ ] Audit log entries cannot be modified after creation
- [ ] Audit log is in a separate table with restricted write access
- [ ] Secrets, API keys, and PII beyond workspace-scoped IDs are never written to audit log
- [ ] Audit log retention policy is configurable (`AUDIT_LOG_RETENTION_DAYS`)

### 2.4 Audit log access

- [ ] `GET /api/v1/admin/audit` endpoint exists and is limited to admin role
- [ ] Audit log is paginated
- [ ] Audit log can be filtered by: time range, actor, workspace, action, resource type, status
- [ ] Audit log can be exported (JSON format minimum)
- [ ] Audit log access itself produces an audit entry

---

## Section 3: Data Security

### 3.1 Encryption in transit

- [ ] All HTTP traffic is served over TLS in production (documented in deployment guide)
- [ ] TLS termination is documented (reverse proxy or direct)
- [ ] Minimum TLS version is 1.2 (prefer 1.3)
- [ ] Certificate management procedure is documented

### 3.2 Encryption at rest

- [ ] PostgreSQL data directory encryption is documented (OS or filesystem level)
- [ ] No sensitive data is written to unencrypted locations
- [ ] Memory chunks (which may contain PII) are protected by workspace access controls

### 3.3 Data retention and deletion

- [ ] Individual memory chunks can be deleted via API
- [ ] All agent data can be deleted when an agent is deleted
- [ ] All workspace data can be deleted when a workspace is deleted
- [ ] Deletion is permanent (no soft-delete with recoverable PII)
- [ ] GDPR data deletion procedure is documented
- [ ] Data retention policies are configurable per workspace (or organization)

### 3.4 Input validation

- [ ] All API inputs are validated before processing
- [ ] SQL tool defaults to read-only (no DDL or DML without explicit permission)
- [ ] HTTP tool has domain allowlist support
- [ ] Shell tool is disabled by default and requires explicit opt-in
- [ ] Maximum payload size is enforced (10 MiB default)

---

## Section 4: Reliability and Operations

### 4.1 Health and availability

- [ ] `GET /health` returns 200 when server is healthy
- [ ] `GET /healthz` is available for Kubernetes liveness probes
- [ ] `GET /readyz` returns 503 until DB is connected and migrations are applied
- [ ] Health endpoints do not require authentication
- [ ] Server starts with a clear error if DB is unreachable (does not silently retry forever)

### 4.2 Graceful shutdown

- [ ] Server handles SIGTERM and SIGINT
- [ ] In-flight requests complete before shutdown
- [ ] Background workers stop cleanly
- [ ] Maximum shutdown timeout is configurable
- [ ] Kubernetes can safely terminate pods using readiness probe + preStop hook

### 4.3 Restart and recovery

- [ ] Server restarts cleanly after crash (no corrupted state)
- [ ] Agent state in PostgreSQL survives server restart
- [ ] Session history survives server restart
- [ ] Memory chunks survive server restart
- [ ] Workflow execution state survives server restart
- [ ] In-progress workflows resume on restart (or are marked failed with retry)

### 4.4 Database reliability

- [ ] Connection pool is configured appropriately for production load
- [ ] DB connection errors produce clear log messages with connection details (no passwords)
- [ ] DB connection retry behavior is documented
- [ ] Migration rollback procedure is documented

---

## Section 5: Backup and Restore

### 5.1 Backup documentation

- [ ] `docs/backup-restore.md` exists
- [ ] Document covers: what to back up, how to back up, how often
- [ ] Full database backup procedure using `pg_dump` is documented
- [ ] Backup file format is documented
- [ ] Backup size estimation guidance is provided

### 5.2 Restore documentation

- [ ] Full restore procedure is documented step by step
- [ ] Restore testing procedure is documented
- [ ] Restore time estimate is provided (or how to estimate it)
- [ ] Point-in-time recovery (PITR) configuration notes are included

### 5.3 Backup verification

- [ ] Backup procedure has been tested at least once (not just documented)
- [ ] Restore procedure has been tested at least once
- [ ] Backup test is part of release checklist

---

## Section 6: Observability

### 6.1 Metrics

- [ ] Prometheus metrics are available at `GET /metrics`
- [ ] Key metrics are instrumented:
  - [ ] `neuronagent_http_requests_total{method, path, status}`
  - [ ] `neuronagent_http_request_duration_seconds{method, path}`
  - [ ] `neuronagent_agent_requests_total{status}`
  - [ ] `neuronagent_tool_executions_total{handler_type, status}`
  - [ ] `neuronagent_memory_writes_total`
  - [ ] `neuronagent_memory_search_duration_seconds`
  - [ ] `neuronagent_rag_requests_total{status}`
  - [ ] `neuronagent_workflow_runs_total{status}`
  - [ ] `neuronagent_neurondb_query_duration_seconds{operation}`
- [ ] Grafana dashboard is available (starter dashboard at minimum)
- [ ] `docs/observability.md` covers all metrics and alert recommendations

### 6.2 Logging

- [ ] All logs are structured JSON when `LOG_FORMAT=json`
- [ ] Log level is configurable (`LOG_LEVEL`: debug, info, warn, error)
- [ ] Every log entry includes request ID for correlation
- [ ] No secrets, API keys, or passwords appear in logs
- [ ] Error logs include enough context to diagnose without source code access

### 6.3 Tracing

- [ ] OpenTelemetry is instrumented (dependency exists — verify instrumentation depth)
- [ ] Trace context propagates through HTTP requests
- [ ] Trace context propagates through agent execution
- [ ] OTLP export endpoint is configurable
- [ ] `docs/observability.md` covers tracing configuration

### 6.4 Alerting guidance

- [ ] Recommended alert thresholds are documented:
  - [ ] High HTTP error rate (5xx > 1% of requests)
  - [ ] High agent execution latency (p99 > 5 seconds)
  - [ ] High tool failure rate (> 5% of executions)
  - [ ] DB query latency (p99 > 500ms)
  - [ ] Worker queue depth (> 100 pending jobs)
  - [ ] Memory search latency (p99 > 200ms)

---

## Section 7: Deployment

### 7.1 Docker deployment

- [ ] Docker image is published and publicly pullable
- [ ] `docs/deploy/docker-compose.md` covers production Docker Compose setup
- [ ] Production compose file uses:
  - [ ] Named volumes for persistence
  - [ ] Resource limits
  - [ ] Health checks
  - [ ] Non-root user
  - [ ] No hardcoded secrets (env or secrets file only)
  - [ ] Network isolation between services

### 7.2 Kubernetes deployment

- [ ] Helm chart exists and is tested
- [ ] `docs/deploy/kubernetes.md` covers Helm chart installation
- [ ] Helm chart includes:
  - [ ] Liveness probe (`/healthz`)
  - [ ] Readiness probe (`/readyz`)
  - [ ] Resource requests and limits
  - [ ] Configurable replica count
  - [ ] Secret management for DB credentials and API key secret
  - [ ] HorizontalPodAutoscaler example
  - [ ] PodDisruptionBudget
- [ ] NeuronDB as external dependency is documented

### 7.3 Configuration management

- [ ] All configuration is via environment variables or config file (no compiled-in secrets)
- [ ] `docs/configuration.md` lists every environment variable with type, default, and description
- [ ] Sensitive variables are clearly marked
- [ ] Production mode refuses to start with insecure defaults

### 7.4 Upgrade procedure

- [ ] Database migration is safe to run while the previous version is still running
- [ ] Upgrade procedure is documented (stop old, run migrations, start new)
- [ ] Rollback procedure is documented
- [ ] Breaking changes between versions are documented in CHANGELOG

---

## Section 8: Security Policies and Compliance

### 8.1 Vulnerability management

- [ ] `SECURITY.md` exists with vulnerability reporting instructions
- [ ] Response time commitment is stated (e.g., acknowledge within 48 hours)
- [ ] Security advisories are published for critical vulnerabilities
- [ ] Docker images are scanned for CVEs before release
- [ ] Go dependencies are scanned for known vulnerabilities

### 8.2 Compliance documentation

- [ ] `docs/security.md` covers the full security model
- [ ] Data processing locations are documented (database is where you deploy it)
- [ ] No telemetry or data is sent to external services without explicit configuration
- [ ] Outbound network calls are documented:
  - [ ] LLM API calls (to configured endpoint — enterprise-configurable)
  - [ ] HTTP tool calls (to agent-configured endpoints — allowlist-controllable)
  - [ ] No required calls to neurondb.ai or any SaaS
- [ ] Air-gapped deployment is possible if LLM endpoint is self-hosted

### 8.3 Secrets management

- [ ] No secrets are stored in the repository
- [ ] No secrets are in Docker images
- [ ] No secrets appear in API responses
- [ ] No secrets appear in logs
- [ ] API keys are stored as hashes in the database

---

## Section 9: Performance

### 9.1 Benchmarks

- [ ] Performance benchmarks exist (`docs/benchmarks.md`)
- [ ] Benchmarks are measured on real hardware (not invented)
- [ ] Benchmark methodology is documented (hardware, software, configuration)
- [ ] Benchmarks cover: agent request latency, memory search, RAG, tool execution, workflow start
- [ ] Results are reproducible via `scripts/benchmark.sh`

### 9.2 Scalability

- [ ] Horizontal scaling is possible (stateless server + shared PostgreSQL)
- [ ] Maximum concurrent requests are documented or benchmarked
- [ ] Database connection pool limits are documented
- [ ] Memory worker scaling behavior is documented (single-writer or multi-writer safe)

### 9.3 SLA guidance

- [ ] Recommended uptime SLA is documented (based on PostgreSQL reliability)
- [ ] Maximum memory search latency under load is measured
- [ ] Maximum agent execution timeout is configurable
- [ ] Rate limits are configurable per API key

---

## Enterprise Readiness Score

| Section | Items | Passing | Score |
|---|---|---|---|
| Auth and RBAC | 20 | ___ | ___% |
| Audit logging | 20 | ___ | ___% |
| Data security | 15 | ___ | ___% |
| Reliability | 16 | ___ | ___% |
| Backup and restore | 8 | ___ | ___% |
| Observability | 18 | ___ | ___% |
| Deployment | 18 | ___ | ___% |
| Security policies | 12 | ___ | ___% |
| Performance | 9 | ___ | ___% |
| **Total** | **136** | ___ | ___% |

**Target for enterprise evaluation:** 80% or above.  
**Current estimate (2026-04-30):** ~35% (auth baseline exists; audit coverage partial; no backup docs; no /readyz; Helm chart unverified; no security policy docs).

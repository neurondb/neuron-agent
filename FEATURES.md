# NeuronAgent Features

Capability matrix and feature reference for NeuronAgent.

| Status | Meaning |
|--------|--------|
| **Supported** | Feature is implemented and supported. |
| **Partial** | Feature is implemented with limitations or schema-only (see section). |
| **Not supported** | Feature is out of scope. |

---

## Scope

NeuronAgent is an agent runtime with REST/WebSocket API, tools, workflows, and memory. This document lists its capabilities and support level by area.

---

## Summary


| Area | Capability | Status |
|------|------------|--------|
| **Agent runtime**              | State machine, runs, plans, steps, tool invocations, model calls                                                                                                                                                                                                                                                                    | Supported |
| **Built-in tool types**        | SQL, HTTP, code, shell, browser, visualization, filesystem, memory, collaboration, multimodal, web search; ML, vector, RAG, analytics, hybrid search, reranking when compatible client configured; NeuronSQL (7 tools: schema_snapshot, validate_sql, explain_json, optimize_candidates, table_profile, index_profile, sample_rows) | Supported |
| **Custom tools**               | DB-stored tools, dynamic registration, versioning, tool-level permissions                                                                                                                                                                                                                                                           | Supported |
| **Workspaces & collaboration** | Workspaces, participant management, collaboration tool                                                                                                                                                                                                                                                                              | Supported |
| **Workflow engine**            | DAG workflows, agent/tool/HTTP/SQL/approval/custom steps, retries, idempotency, compensation, scheduling, audit                                                                                                                                                                                                                     | Supported |
| **Human-in-the-loop**          | Approval gates, feedback, email/webhook notifications                                                                                                                                                                                                                                                                               | Supported |
| **Memory**                     | Hierarchical memory (STM/MTM/LPM), vector search (HNSW), promotion, summarization, feedback, quality metrics, forgetting, conflict resolution; episodic memory schema                                                                                                                                                               | Supported / Partial                    |
| **Budget & cost**              | Per-agent and per-session budgets, tracking, alerts, period-based, cost analytics                                                                                                                                                                                                                                                   | Supported |
| **REST API**                   | CRUD for agents, sessions, tools, workflows, etc.; pagination, filtering; API key auth                                                                                                                                                                                                                                              | Supported |
| **WebSocket**                  | Streaming responses, message queue, keepalive                                                                                                                                                                                                                                                                                       | Supported |
| **Auth & RBAC**                | API keys, bcrypt, RBAC, principal/session tool permissions, workspace policies                                                                                                                                                                                                                                                      | Supported |
| **Database schemas**           | Runtime, tenancy, RBAC, audit (v2), compliance, marketplace, tool versioning, self-improvement, reliability, observability, advanced memory, distributed                                                                                                                                                                            | Supported (schema)                     |
| **Audit & compliance**         | Audit tables (audit_events, policy_decisions, tool_invocations, workflow_executions_audit, approval_actions); compliance_reports (report_type: soc2, iso27001, gdpr)                                                                                                                                                                | Supported (schema) / Partial (reports) |
| **Marketplace**                | Schema only (marketplace_tools, tool_ratings, marketplace_agents, marketplace_workflows)                                                                                                                                                                                                                                            | Partial                           |
| **Claw gateway**               | `/claw/v1` tools/list, tools/run, health — NeuronSQL tools only                                                                                                                                                                                                                                                                     | Supported |
| **NeuronSQL**                  | SQL validation, explain, optimize, schema/table/index profile, sample rows; policy engine; workflow templates                                                                                                                                                                                                                       | Supported |


---

## Supported

**Agent runtime:** State machine for task execution; persistent runs, plans, steps; tool invocations and model calls stored; execution traces, retrieval events, context builds.

**Built-in tool handlers:** `sql`, `http`, `code`, `shell`, `browser`, `visualization`, `filesystem`, `memory`, `collaboration`, `multimodal`, `web_search`. When a compatible database client is configured: `ml`, `vector`, `rag`, `analytics`, `hybrid_search`, `reranking`. NeuronSQL (Claw): `schema_snapshot`, `validate_sql`, `explain_json`, `optimize_candidates`, `table_profile`, `index_profile`, `sample_rows`. Custom tools can be registered from the database.

**Workspaces:** Create workspaces; add users/agents; collaboration tool for shared context.

**Workflow engine:** DAG-based workflows; step types: agent, tool, HTTP, SQL, approval, custom; conditional logic, retry, idempotency, compensation steps; cron scheduling; execution history and audit (e.g. `workflow_executions_audit`).

**Human-in-the-loop:** Approval gates in workflows; feedback collection; email and webhook notifications for approvals.

**Memory:** Hierarchical memory (short/medium/long-term); HNSW-based vector search; promotion, summarization; user feedback; quality metrics; forgetting strategies; conflict resolution. Schema for episodic memory and related features.

**Budget & cost:** Per-agent and per-session budgets; real-time cost tracking; budget alerts; daily/weekly/monthly/yearly or total budgets; cost analytics.

**API:** REST CRUD for agents, sessions, tools, workflows, webhooks, etc.; pagination and filtering; API key authentication. WebSocket for streaming and message queue.

**Auth:** API keys, bcrypt hashing, RBAC (admin/user), principal and session tool permissions, workspace policies.

**Database:** Schemas for runtime, tenancy (organizations, quotas), RBAC (principal_tool_permissions, workflow_permissions, workspace_policies), audit v2, compliance (audit_logs, compliance_reports), marketplace (marketplace_tools, tool_ratings, marketplace_agents, marketplace_workflows), tool versioning, self-improvement (execution_results, performance_feedback, ab_tests), reliability (dead_letter_queue), observability (execution_trace, performance_profiles), advanced memory, distributed (cluster_nodes, events, cache).

**Claw:** `POST /claw/v1/tools/list`, `POST /claw/v1/tools/run`, `GET /claw/v1/health`; only `neuronsql.`* tools are listed and executable.

**NeuronSQL:** Policy engine, sensitive tables, validation/explain/optimize/profile/sample tools; YAML workflow templates with audit and retries.

---

## Partial support

**Compliance reports:** Tables exist for `audit_logs` and `compliance_reports` with `report_type IN ('soc2', 'iso27001', 'gdpr')`. Report generation (e.g. GDPR, HIPAA, SOX) may be basic or placeholder; HIPAA/SOX are not in the schema enum.

**Marketplace:** Schema is present (tools, agents, workflows, ratings). ROADMAP states “No marketplace” as current state — publish/discover/rate platform UI or services may be incomplete or planned.

**Memory:** Advanced memory features (e.g. consolidation, compression, cross-session sharing) may vary by configuration and implementation maturity.

**NeuronDB integration:** ML, vector, RAG, analytics, hybrid search, reranking tools work when a compatible database client is configured; otherwise those tool types are not registered.

---

## Out of scope

- **Tool/Agent marketplace platform:** Schema exists; full marketplace product (publish, discover, rate, categories) is not fully implemented per ROADMAP.

---

## Counts (reference)

- **Built-in tool handler types:** ~18 (11 generic + 6 when compatible client configured + 7 NeuronSQL; some overlap with “tool types”).
- **NeuronSQL tools (Claw):** 7.
- **Compliance report types in schema:** 3 (soc2, iso27001, gdpr).

---

## Documentation

- [README](README.md)
- [API reference](docs/api-reference.md)
- [NeuronSQL design](docs/neuronsql_design.txt)
- [Workflow templates](docs/workflow_templates/README.md)
- [ROADMAP](ROADMAP.md) for planned vs current state

---

[Back to top](#neuronagent-features) · [README](README.md)
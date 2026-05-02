# Reliability

- **Graceful shutdown** — SIGINT/SIGTERM triggers server shutdown and module stop (`cmd/agent-server/main.go`).
- **Readiness** — use **`GET /readyz`** in Kubernetes before marking pods ready.
- **Workflows** — executions persist in SQL; partial failure surfaces in workflow execution records.

For operational runbooks see [`operations_runbook.txt`](operations_runbook.txt).

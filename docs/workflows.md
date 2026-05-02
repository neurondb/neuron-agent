# Workflows

Workflows are DAGs stored in PostgreSQL with steps executed by `internal/workflow/engine.go`.

## Step types

| Type | Behavior |
|------|----------|
| `agent` | Runs agent runtime with a dedicated session per execution |
| `tool` | Executes a named tool |
| `sql` | Executes configured SQL |
| `http` | HTTP request step |
| `approval` | Human approval gate |
| `conditional` | Evaluates `condition` expression via `AdvancedWorkflowEngine.EvaluateCondition` |

Retries: steps with `retry_config` run synchronously with backoff (`retry_delay_seconds` or capped linear delay). Cron schedules in `workflow_schedules` are polled by the server (`WORKFLOW_SCHEDULE_INTERVAL`, disable with `WORKFLOW_SCHEDULE_ENABLED=false`).

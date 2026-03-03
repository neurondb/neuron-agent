# NeuronSQL Workflow Templates (Part D)

These YAML files define **ready workflows** for NeuronSQL: audit, metrics, retries, and idempotency are supported by the workflow engine.

## Templates

| Template | Description |
|----------|-------------|
| `neuronsql_generate_and_validate.yaml` | Generate SELECT with `neuronsql.generate_select`, validate with `neuronsql.validate_sql` (retries, idempotency). |
| `neuronsql_optimize_with_approval_gate.yaml` | Optimize SELECT with `neuronsql.optimize_select`, then approval gate before follow-up. |
| `neuronsql_plpgsql_generate_with_review.yaml` | Generate PL/pgSQL with `neuronsql.plpgsql_generate`, then review/approval step. |

## Loading via API

1. **Create workflow**  
   `POST /api/v1/workflows` with `name` and `dag_definition` (e.g. copy from the YAML).

2. **Add steps**  
   For each step in the YAML, call  
   `POST /api/v1/workflows/{workflow_id}/steps`  
   with `step_name`, `step_type`, `inputs`, `outputs`, `dependencies`, `retry_config`, and optional `idempotency_key`.

3. **Execute**  
   `POST /api/v1/workflows/{workflow_id}/execute` with `inputs` and optional `trigger_type` / `trigger_data`.

Workflow runs are recorded in the audit log when the engine has an audit logger configured. Step-level `timeout_seconds` in `retry_config` enforces a per-step time budget.

## Engine behavior

- **Audit**: `workflow_run` and approval events are logged when an audit logger is set on the workflow engine.
- **Metrics**: Workflow and step execution metrics are emitted by the engine.
- **Retries**: Use `retry_config.max_attempts` (and optional `timeout_seconds`) on each step.
- **Idempotency**: Set `idempotency_key` on steps that should return cached results when re-run with the same key.

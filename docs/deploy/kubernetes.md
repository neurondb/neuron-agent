# Kubernetes deployment

Helm chart skeleton: [`helm/`](../helm/) at repository root.

## Checklist

- Set **secrets** for `DB_PASSWORD` and any LLM sidecar keys.
- Point NeuronAgent at a **NeuronDB-capable** PostgreSQL service (or external managed Postgres).
- Configure **liveness** → `GET /healthz`, **readiness** → `GET /readyz`.
- Scrape **`GET /metrics`** from your Prometheus.

Review `helm/values.yaml` and templates for resource limits and ingress; adjust for your cluster policy.

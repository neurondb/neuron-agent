# Observability

- **`GET /metrics`** — Prometheus exposition format (`internal/metrics`).
- **`GET /health`**, **`GET /healthz`** — process up.
- **`GET /readyz`** — database ping; returns 503 until PostgreSQL is reachable.
- **`GET /version`** — build metadata for support tickets.
- **Logging** — structured logging via middleware (`internal/api/middleware.go`); JSON format recommended in production.

Sample Prometheus scrape target:

```yaml
scrape_configs:
  - job_name: neuronagent
    static_configs:
      - targets: ["localhost:8080"]
    metrics_path: /metrics
```

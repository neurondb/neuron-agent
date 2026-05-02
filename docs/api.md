# HTTP API

The authoritative contract is the **OpenAPI 3** specification served by the running server:

- **Interactive docs:** `GET /docs` (Redoc)
- **Raw spec:** `GET /docs/openapi.yaml`

Same-origin path prefix for REST routes is typically **`/api/v1`**. OpenClaw-compatible endpoints live under **`/claw/v1`**.

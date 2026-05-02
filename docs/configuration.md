# Configuration

NeuronAgent loads YAML from `CONFIG_PATH` if set, then applies environment overrides (see `internal/config/env.go`).

## Common environment variables

| Variable | Purpose |
|----------|---------|
| `SERVER_HOST`, `SERVER_PORT` | HTTP bind address |
| `SERVER_READ_TIMEOUT`, `SERVER_WRITE_TIMEOUT` | HTTP timeouts |
| `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` | PostgreSQL connection |
| `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS` | Pool sizing |
| `LOG_LEVEL`, `LOG_FORMAT` | `debug`/`info`/… and `json` or `console` |
| `CORS_ALLOWED_ORIGINS` | Comma-separated origins |
| `WEBSOCKET_ALLOWED_ORIGINS` | WebSocket origins (defaults to CORS if unset) |
| `CONFIG_PROFILE` | `development`, `staging`, `production` — production validation rejects weak default DB passwords |
| `ENV` | Alias used with profile for some middleware behavior |
| `MODULE_NEURONSQL_ENABLED` | `true`/`false` NeuronSQL module |
| `TOOLS_TIMEOUT` | Per-tool execution timeout |
| `WORKFLOW_MAX_DURATION` | Workflow run wall-clock limit |
| `WORKFLOW_SCHEDULE_ENABLED`, `WORKFLOW_SCHEDULE_INTERVAL` | Cron workflow scheduler (`true` / e.g. `30s`) |
| `NEURONDB_REPO`, `NEURONDB_IMAGE`, `PG_MAJOR` | Build Postgres from a local NeuronDB clone (`docker-compose.neurondb.yml` / `make up-neurondb`) |
| `LLM_SQL_BASE_URL`, `LLM_SQL_API_KEY` | Sidecar for `/api/v1/llm/sql/*` routes |

See also `docs/config_env_schema.txt` for extended notes. After editing `src/openapi/openapi.yaml`, run **`make sync-openapi`** so the embedded spec under `internal/api/specdata/` stays in sync.

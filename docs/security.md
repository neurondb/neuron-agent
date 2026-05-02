# Security

- **Authentication** — API keys via `Authorization: Bearer` or `ApiKey`; middleware in `internal/api/middleware.go`.
- **Roles** — principals resolved per key; admin routes require elevated role where enforced.
- **Production hardening** — `ValidateConfig` rejects common default database passwords when profile/env indicates production (`internal/config/config.go`).
- **Secrets** — never commit `.env`; use strong `DB_PASSWORD` and rotate API keys.

See also [`security_model.txt`](security_model.txt) for extended notes.

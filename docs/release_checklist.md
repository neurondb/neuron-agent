# Release checklist

Use this checklist for each release to ensure compatibility and safe rollout.

## Pre-release

- [ ] **Version bump** – Update version in Chart.yaml, package.json / go.mod / VERSION as appropriate.
- [ ] **Migrations** – Ensure all new SQL migrations are in `sql/` or `migrations/` and documented. Run migrations against a copy of production schema and verify no errors.
- [ ] **Compat checks** – Confirm existing `/api/v1` endpoints and behavior are unchanged. Run existing integration tests.
- [ ] **Smoke tests** – Run smoke tests against a staging deployment (health, create agent, send message, NeuronSQL generate/validate, Claw tools/list).

## Release

- [ ] **Tag** – Create git tag (e.g. `v1.2.0`).
- [ ] **Build** – Build and push container image for the tag.
- [ ] **Changelog** – Update CHANGELOG.md with the release notes.

## Post-release / Rollback

- [ ] **Rollback plan** – Document rollback steps: revert to previous image, re-run previous migrations if any were not backward-compatible (prefer additive-only migrations).
- [ ] **Monitor** – After deploy, monitor metrics (latency, errors, policy_blocks) and audit logs.
- [ ] **Graceful shutdown** – Server uses `Shutdown(ctx)` to drain in-flight requests; audit writes are synchronous so no flush step is required.

## Optional CI automation

- Run migrations in CI against a fixture DB.
- Run eval suite and publish report artifact: `neuronsql eval --dsn ... --suite saas_basic --out eval_report.json`.
- Run compatibility test suite before tagging.

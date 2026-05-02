# Contributing to NeuronAgent

## Development setup

1. Install **Go 1.24+** and **PostgreSQL** (or use Docker Compose).
2. Clone the repo and copy env: `cp .env.example .env`
3. Build from repo root: `make build` (binary at `bin/neuron-agent`).
4. Apply schema: `make migrate` (requires `psql` and a running database).
5. Run tests: `make test-fast` or `cd src && go test ./... -short`

After editing [`src/openapi/openapi.yaml`](src/openapi/openapi.yaml), sync the embedded copy used by `go:embed`:

```bash
make sync-openapi
```

Go code lives under [`src/`](src/). The Go module path is `github.com/neurondb/NeuronAgent`.

## Pull requests

- Keep changes focused; match existing style and imports.
- Run `make fmt` and `make test-fast` before sending a PR.
- Document user-visible behavior in README or `docs/` when you change defaults or APIs.

## Security

Report vulnerabilities according to [SECURITY.md](SECURITY.md).

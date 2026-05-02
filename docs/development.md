# Development

```bash
make build          # bin/neuron-agent
make test-fast      # all Go tests (no race)
make fmt && make lint
make sync-openapi   # after editing src/openapi/openapi.yaml
```

Primary module directory: **`src/`**. HTTP server entry: **`cmd/agent-server`**.

See [CONTRIBUTING.md](../CONTRIBUTING.md).

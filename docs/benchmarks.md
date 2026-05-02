# Benchmarks

Use [`scripts/benchmark.sh`](../scripts/benchmark.sh) for a minimal latency sample on `/health`.

For load testing, install [hey](https://github.com/rakyll/hey) or [wrk](https://github.com/wg/wrk) and target:

```bash
hey -n 1000 -c 50 http://127.0.0.1:8080/health
```

Document hardware, NeuronDB version, and connection pool settings when publishing results.

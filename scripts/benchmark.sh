#!/usr/bin/env bash
# Placeholder benchmark harness — extend with hey/wrk against /health and authenticated routes.
set -euo pipefail
BASE_URL="${NEURON_AGENT_URL:-http://127.0.0.1:8080}"
echo "Benchmark placeholder: GET ${BASE_URL}/health"
for i in 1 2 3 4 5; do
  curl -s -o /dev/null -w "request %{time_total}s HTTP %{http_code}\n" "${BASE_URL}/health"
done
echo "Install 'hey' or 'wrk' for load tests; see docs/benchmarks.md"

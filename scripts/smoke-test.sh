#!/usr/bin/env bash
# Smoke test: public endpoints only (no API key required).
set -euo pipefail

BASE_URL="${NEURON_AGENT_URL:-http://127.0.0.1:8080}"
BASE_URL="${BASE_URL%/}"

fail() { echo "FAIL: $*" >&2; exit 1; }

echo "== smoke: GET ${BASE_URL}/health"
curl -sf "${BASE_URL}/health" | grep -q ok || fail "/health"

echo "== smoke: GET ${BASE_URL}/healthz"
curl -sf "${BASE_URL}/healthz" | grep -q ok || fail "/healthz"

echo "== smoke: GET ${BASE_URL}/version"
curl -sf "${BASE_URL}/version" | grep -q api_version || fail "/version"

echo "== smoke: GET ${BASE_URL}/docs/openapi.yaml"
n=$(curl -sf "${BASE_URL}/docs/openapi.yaml" | wc -c | tr -d ' ')
test "${n}" -gt 1000 || fail "/docs/openapi.yaml too small"

echo "== smoke: GET ${BASE_URL}/readyz"
code=$(curl -s -o /tmp/readyz.json -w "%{http_code}" "${BASE_URL}/readyz")
if [[ "${code}" != "200" && "${code}" != "503" ]]; then
  fail "/readyz unexpected HTTP ${code}"
fi
echo "readyz HTTP ${code} body: $(cat /tmp/readyz.json)"

echo "OK: smoke tests passed"

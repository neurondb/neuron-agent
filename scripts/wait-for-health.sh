#!/usr/bin/env bash
# Wait until NeuronAgent responds on /health (or exit 1).
set -euo pipefail

BASE_URL="${1:-http://127.0.0.1:8080}"
MAX_ATTEMPTS="${2:-60}"
SLEEP_SECS="${3:-2}"

for i in $(seq 1 "${MAX_ATTEMPTS}"); do
  if curl -sf "${BASE_URL%/}/health" >/dev/null; then
    echo "NeuronAgent is healthy at ${BASE_URL}"
    exit 0
  fi
  echo "waiting for NeuronAgent... (${i}/${MAX_ATTEMPTS})"
  sleep "${SLEEP_SECS}"
done

echo "ERROR: NeuronAgent did not become healthy. Try: docker compose logs neuronagent" >&2
exit 1

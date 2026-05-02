#!/usr/bin/env bash
# Golden-path demo when NEURON_AGENT_API_KEY is set; otherwise runs smoke only.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
# shellcheck source=/dev/null
[[ -f "${ROOT}/.demo.env" ]] && source "${ROOT}/.demo.env"
[[ -f "${ROOT}/.env" ]] && set -a && source "${ROOT}/.env" && set +a

BASE_URL="${NEURON_AGENT_URL:-http://127.0.0.1:8080}"
BASE_URL="${BASE_URL%/}"
API_KEY="${NEURON_AGENT_API_KEY:-}"

bash "${ROOT}/scripts/wait-for-health.sh" "${BASE_URL}" 30 2

if [[ -z "${API_KEY}" ]]; then
  echo "NEURON_AGENT_API_KEY not set — running smoke-only (public endpoints)."
  bash "${ROOT}/scripts/smoke-test.sh"
  echo "To run authenticated demo: add key to .demo.env (see scripts/bootstrap-demo.sh and cmd/generate-key)"
  exit 0
fi

echo "== demo: create agent"
resp=$(curl -sf -X POST "${BASE_URL}/api/v1/agents" \
  -H "Authorization: Bearer ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{"name":"demo-agent","system_prompt":"You are a concise assistant."}') || {
  echo "Agent create failed — check API key and DB migrations" >&2
  exit 1
}
echo "${resp}" | head -c 500
echo
echo "Demo agent created (see response above). OK."

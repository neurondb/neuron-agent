#!/usr/bin/env bash
set -euo pipefail
BASE="${NEURON_AGENT_URL:-http://127.0.0.1:8080}"
KEY="${NEURON_AGENT_API_KEY:?set NEURON_AGENT_API_KEY}"
AGENT_ID="${AGENT_ID:?set AGENT_ID to an existing agent uuid}"

resp=$(curl -sf -X POST "${BASE}/api/v1/sessions" \
  -H "Authorization: Bearer ${KEY}" -H "Content-Type: application/json" \
  -d "{\"agent_id\":\"${AGENT_ID}\"}")
SID=$(printf '%s' "$resp" | python3 -c "import sys,json; print(json.load(sys.stdin).get('id',''))")
[[ -n "$SID" ]] || { echo "$resp"; exit 1; }

curl -sf -X POST "${BASE}/api/v1/sessions/${SID}/messages" \
  -H "Authorization: Bearer ${KEY}" -H "Content-Type: application/json" \
  -d '{"role":"user","content":"Hello from quickstart example"}'
echo

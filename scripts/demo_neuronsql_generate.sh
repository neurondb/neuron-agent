#!/usr/bin/env bash
# Demo: POST /api/v1/neuronsql/generate (requires NeuronAgent running with auth disabled or valid API key)
set -e
API="${NEURONAGENT_URL:-http://localhost:8080}"
DSN="${DB_DSN:-host=localhost port=5433 user=neurondb password=neurondb dbname=neurondb sslmode=disable}"
curl -s -X POST "$API/api/v1/neuronsql/generate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${NEURONAGENT_API_KEY:-demo}" \
  -d "{\"db_dsn\": \"$DSN\", \"question\": \"List all products\"}" | jq .

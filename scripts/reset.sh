#!/usr/bin/env bash
# Tear down stack and remove volumes (destructive).
set -euo pipefail
cd "$(dirname "$0")/.."
docker compose down -v "$@"
echo "Stack removed. Run: docker compose up -d"

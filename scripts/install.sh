#!/usr/bin/env bash
# One-command style setup from repo root (expects clone already present).
set -euo pipefail

RED='\033[0;31m'; GREEN='\033[0;32m'; NC='\033[0m'
info() { echo -e "${GREEN}[install]${NC} $*"; }
err() { echo -e "${RED}[install]${NC} $*" >&2; exit 1; }

command -v docker >/dev/null 2>&1 || err "Docker is required: https://docs.docker.com/get-docker/"
docker info >/dev/null 2>&1 || err "Docker daemon is not running."
docker compose version >/dev/null 2>&1 || err "docker compose (v2) is required."

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "${ROOT}"

[[ -f .env ]] || { cp .env.example .env && info "Created .env from .env.example"; }
set -a
# shellcheck source=/dev/null
[[ -f .env ]] && . ./.env
set +a
PORT="${SERVER_HOST_PORT:-8080}"

info "Building and starting stack..."
docker compose build neuronagent
docker compose up -d

bash "${ROOT}/scripts/wait-for-health.sh" "http://127.0.0.1:${PORT}" 60 2

info "Done. Next: apply DB schema if needed (make migrate), generate API key, run scripts/demo.sh"

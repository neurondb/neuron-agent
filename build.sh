#!/usr/bin/env bash
# build.sh - Build NeuronAgent (Go binary and runtime assets)
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROGNAME="${0##*/}"
cd "$SCRIPT_DIR"

usage() {
    cat <<EOF
NeuronAgent Build Script

Build the NeuronAgent Go binary and runtime assets into bin/.

Usage: $PROGNAME [OPTIONS]

Options:
  -h, --help    Show this help and exit

Runs: make build
Produces: bin/neuron-agent, bin/scripts/, bin/conf/, bin/sql/

EOF
}

case "${1:-}" in
    -h|--help) usage; exit 0 ;;
esac

exec make build

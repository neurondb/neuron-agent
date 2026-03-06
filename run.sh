#!/usr/bin/env bash
# run.sh - List or run a built binary
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROGNAME="${0##*/}"
cd "$SCRIPT_DIR"

usage() {
    cat <<EOF
NeuronAgent Run Script

List or run built binaries from this repo.

Usage: $PROGNAME [OPTIONS]
   or: $PROGNAME [OPTIONS] --binary NAME [ARGS...]

Options:
  -d                Run as daemon (background)
  -h, --help        Show this help and exit
  --binary NAME     Run binary NAME (remaining args passed to it)

Binaries (after build): neuron-agent

Examples:
  $PROGNAME --binary neuron-agent      Run agent server

EOF
}

list_binaries() {
    [[ -d bin ]] || return 0
    for f in bin/*; do
        [[ -f "$f" ]] && [[ -x "$f" ]] && [[ "$f" != *.sh ]] && echo "$f"
    done
}

# Load .env from repo root if present (e.g. JWT_SECRET, DATABASE_URL)
load_env() {
    if [[ -f "$SCRIPT_DIR/.env" ]]; then
        set -a
        # shellcheck source=/dev/null
        source "$SCRIPT_DIR/.env"
        set +a
    fi
}

run_binary() {
    local name="$1"
    local as_daemon="${2:-0}"
    shift 2
    local path
    while IFS= read -r path; do
        if [[ "$(basename "$path")" == "$name" ]]; then
            load_env
            if [[ "$as_daemon" -eq 1 ]]; then
                local logfile="$SCRIPT_DIR/daemon-${name}.log"
                nohup "$path" "$@" </dev/null >> "$logfile" 2>&1 &
                echo "$PROGNAME: started $name as daemon (PID $!, log: $logfile)"
                exit 0
            else
                exec "$path" "$@"
            fi
        fi
    done < <(list_binaries)
    echo "$PROGNAME: no such binary: $name" >&2
    exit 1
}

DAEMON=0
BINARY_NAME=""
while [[ $# -gt 0 ]]; do
    case "${1:-}" in
        -d)
            DAEMON=1
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        --what)
            if [[ "${2:-}" != "binary" ]]; then
                echo "$PROGNAME: missing or invalid argument for --what" >&2
                usage >&2
                exit 1
            fi
            list_binaries
            exit 0
            ;;
        --binary)
            if [[ -z "${2:-}" ]]; then
                echo "$PROGNAME: option --binary requires an argument" >&2
                usage >&2
                exit 1
            fi
            BINARY_NAME="$2"
            shift 2
            break
            ;;
        *)
            if [[ -n "${1:-}" ]]; then
                if [[ "$1" =~ ^--b.*ary$ ]] && [[ "$1" != "--binary" ]]; then
                    echo "$PROGNAME: unknown option: $1 (did you mean --binary?)" >&2
                else
                    echo "$PROGNAME: unknown option: $1" >&2
                fi
            fi
            usage >&2
            exit 1
            ;;
    esac
done

if [[ -z "$BINARY_NAME" ]]; then
    usage >&2
    exit 1
fi

run_binary "$BINARY_NAME" "$DAEMON" "$@"

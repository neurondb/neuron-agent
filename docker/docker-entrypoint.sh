#!/bin/sh
# Docker entrypoint script for NeuronAgent
# Performs pre-start validation and initialization

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1" >&2
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1" >&2
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# Check if binary exists
if [ ! -f "/app/neuronagent" ]; then
    log_error "Binary /app/neuronagent not found!"
    exit 1
fi

if [ ! -x "/app/neuronagent" ]; then
    log_error "Binary /app/neuronagent is not executable!"
    exit 1
fi

log_info "Binary found and executable"

# Validate environment variables
if [ -z "${DB_HOST}" ]; then
    log_warn "DB_HOST not set, will use default or config file"
fi

if [ -z "${DB_DATABASE}" ] && [ -z "${DB_NAME}" ]; then
    log_warn "DB_NAME not set, will use default or config file"
fi

# Optional: Test database connectivity (requires psql or similar)
# Uncomment if you want to verify database connection before starting
# if command -v psql >/dev/null 2>&1; then
#     log_info "Testing database connectivity..."
#     CONN_STR="postgresql://${DB_USER:-neurondb}:${DB_PASSWORD:-neurondb}@${DB_HOST:-localhost}:${DB_PORT:-5432}/${DB_NAME:-neurondb}"
#     
#     if psql "${CONN_STR}" -c "SELECT 1;" >/dev/null 2>&1; then
#         log_info "Database connection successful"
#     else
#         log_warn "Database connection test failed (continuing anyway)"
#     fi
# fi

# Validate config file if specified
if [ -n "${CONFIG_PATH}" ] && [ -f "${CONFIG_PATH}" ]; then
    log_info "Config file found: ${CONFIG_PATH}"
    # Basic YAML validation (requires python3 and pyyaml)
    if command -v python3 >/dev/null 2>&1; then
        if python3 -c "import yaml; yaml.safe_load(open('${CONFIG_PATH}'))" >/dev/null 2>&1; then
            log_info "Config file is valid YAML"
        else
            log_warn "Config file validation failed (continuing anyway)"
        fi
    fi
elif [ -n "${CONFIG_PATH}" ]; then
    log_warn "Config file specified but not found: ${CONFIG_PATH}"
fi

# Check if initial schema needs to be run
INITIAL_SCHEMA="/app/initial_schema.sql"
if [ -f "$INITIAL_SCHEMA" ]; then
    log_info "Found initial_schema.sql file"
    
    # Check if NeuronAgent schema is already set up (idempotency check)
    if command -v psql >/dev/null 2>&1; then
        # Build connection string
        if [ -n "${DB_PASSWORD}" ]; then
            export PGPASSWORD="${DB_PASSWORD}"
        fi
        CONN_STR="postgresql://${DB_USER:-neurondb}:${DB_PASSWORD:-neurondb}@${DB_HOST:-localhost}:${DB_PORT:-5432}/${DB_NAME:-neurondb}"
        
        SCHEMA_EXISTS=$(psql "$CONN_STR" -tAc "SELECT 1 FROM information_schema.tables WHERE table_schema = 'neurondb_agent' AND table_name = 'agents'" 2>/dev/null || echo "0")
        
        if [ "$SCHEMA_EXISTS" = "1" ]; then
            log_info "NeuronAgent schema already exists, skipping initial schema setup"
        else
            log_info "NeuronAgent schema not found, initial schema setup may need to be run manually"
            log_info "To setup schema manually: psql \$CONN_STR -f $INITIAL_SCHEMA"
        fi
    fi
elif [ -d "/app/sql" ]; then
    MIGRATION_COUNT=$(find "/app/sql" -name "*.sql" | wc -l)
    log_info "Found ${MIGRATION_COUNT} migration file(s) in /app/sql (initial_schema.sql not found, using migration files)"
else
    log_warn "Initial schema file not found at $INITIAL_SCHEMA and sql directory not found"
fi

# Log startup information
log_info "Starting NeuronAgent server..."
log_info "  Database Host: ${DB_HOST:-localhost}"
log_info "  Database Port: ${DB_PORT:-5433}"
log_info "  Database Name: ${DB_NAME:-neurondb}"
log_info "  Database User: ${DB_USER:-neurondb}"
log_info "  Server Host: ${SERVER_HOST:-0.0.0.0}"
log_info "  Server Port: ${SERVER_PORT:-8080}"
log_info "  Log Level: ${LOG_LEVEL:-info}"
log_info "  Log Format: ${LOG_FORMAT:-json}"

# Execute the binary with all arguments
exec /app/neuronagent "$@"


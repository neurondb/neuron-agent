# NeuronAgent Comprehensive Test Suite

This directory contains a comprehensive test suite for NeuronAgent that tests **every single feature** of the system.

## Structure

The test suite is organized into logical modules:

- `test_api/` - API endpoint tests (65+ handlers)
- `test_runtime/` - Agent runtime core features
- `test_memory/` - Memory & knowledge features
- `test_tools/` - Tool system (20+ tools)
- `test_collaboration/` - Multi-agent features
- `test_workflow/` - Workflow engine
- `test_planning/` - Planning & task management
- `test_quality/` - Quality & evaluation
- `test_budget/` - Budget & cost management
- `test_hitl/` - Human-in-the-loop
- `test_versioning/` - Versioning & history
- `test_observability/` - Observability & monitoring
- `test_security/` - Security & safety
- `test_integrations/` - Integrations & connectors
- `test_storage/` - Storage & persistence
- `test_workers/` - Background workers
- `test_neurondb/` - NeuronDB integration
- `test_integration/` - Integration tests
- `test_database/` - Database tests

## Running Tests

### Prerequisites

1. Install test dependencies:
```bash
pip install -r tests/requirements.txt
```

2. Ensure NeuronAgent server is running:
```bash
# Start server
cd NeuronAgent
DB_USER=pge DB_PASSWORD="" go run cmd/agent-server/main.go
```

3. Set environment variables (optional):
```bash
export NEURONAGENT_BASE_URL=http://localhost:8080
export NEURONAGENT_API_KEY=your_api_key
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=neurondb
export DB_USER=pge
export DB_PASSWORD=""
```

### Run All Tests

```bash
# From NeuronAgent directory
pytest tests/

# Or use the test runner script
./tests/run_tests.sh all
```

### Run Specific Test Categories

```bash
# API tests only
pytest tests/test_api/ -v

# Tool tests only
pytest tests/test_tools/ -v

# NeuronDB integration tests
pytest tests/test_neurondb/ -v

# Fast tests (skip slow ones)
pytest tests/ -v -m "not slow"

# With coverage
pytest tests/ --cov=NeuronAgent --cov-report=html
```

### Run Specific Test File

```bash
pytest tests/test_api/test_agents.py -v
```

### Run Specific Test

```bash
pytest tests/test_api/test_agents.py::TestAgentCRUD::test_create_agent -v
```

## Test Markers

Tests are marked with categories for easy filtering:

- `@pytest.mark.unit` - Unit tests (fast, isolated)
- `@pytest.mark.integration` - Integration tests
- `@pytest.mark.api` - API endpoint tests
- `@pytest.mark.tool` - Tool execution tests
- `@pytest.mark.neurondb` - NeuronDB integration tests
- `@pytest.mark.slow` - Slow running tests
- `@pytest.mark.requires_db` - Requires database connection
- `@pytest.mark.requires_server` - Requires running server
- `@pytest.mark.requires_neurondb` - Requires NeuronDB extension
- `@pytest.mark.performance` - Performance benchmarks
- `@pytest.mark.security` - Security tests

## Test Coverage

The test suite aims for:
- **100% coverage** of all API endpoints
- **100% coverage** of all tools
- **100% coverage** of NeuronDB integration points
- **>90% overall code coverage**

## Fixtures

Common fixtures available in `conftest.py`:

- `api_client` - Authenticated API client
- `test_agent` - Create test agent (auto-cleanup)
- `test_session` - Create test session (auto-cleanup)
- `api_key` - API key for testing
- `db_connection` - Database connection
- `neurondb_available` - Check if NeuronDB extension is available
- `unique_name` - Generate unique names for test resources
- `test_data_generator` - Faker instance for test data

## Troubleshooting

### Server Not Available

If tests are skipped due to server not being available:
```bash
# Start the server
cd NeuronAgent
DB_USER=pge DB_PASSWORD="" go run cmd/agent-server/main.go
```

### Database Connection Issues

Ensure PostgreSQL is running and accessible:
```bash
psql -h localhost -p 5432 -U pge -d neurondb -c "SELECT 1;"
```

### Missing API Key

Generate an API key:
```bash
cd NeuronAgent
go run cmd/generate-key/main.go -org test -user test -rate 1000 -roles user,admin
```

### Import Errors

Ensure you're running tests from the NeuronAgent directory:
```bash
cd NeuronAgent
pytest tests/
```

## Contributing

When adding new tests:

1. Place tests in the appropriate category directory
2. Use appropriate markers (`@pytest.mark.*`)
3. Use fixtures from `conftest.py` for common setup
4. Clean up test resources (fixtures handle this automatically)
5. Add docstrings explaining what the test verifies


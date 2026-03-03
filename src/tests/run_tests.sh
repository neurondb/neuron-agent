#!/bin/bash
# Test runner script for NeuronAgent comprehensive test suite

set -e

cd "$(dirname "$0")/.."

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}NeuronAgent Comprehensive Test Suite${NC}"
echo "=========================================="
echo ""

# Check if pytest is installed
if ! command -v pytest &> /dev/null; then
    echo -e "${YELLOW}Installing test dependencies...${NC}"
    pip install -r tests/requirements.txt
fi

# Check if server is running
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${YELLOW}Warning: Server may not be running on http://localhost:8080${NC}"
    echo "Some tests may be skipped."
    echo ""
fi

# Run tests based on arguments
if [ "$1" == "all" ] || [ -z "$1" ]; then
    echo -e "${GREEN}Running all tests...${NC}"
    pytest tests/ -v
elif [ "$1" == "api" ]; then
    echo -e "${GREEN}Running API tests...${NC}"
    pytest tests/test_api/ -v -m api
elif [ "$1" == "tools" ]; then
    echo -e "${GREEN}Running tool tests...${NC}"
    pytest tests/test_tools/ -v -m tool
elif [ "$1" == "neurondb" ]; then
    echo -e "${GREEN}Running NeuronDB integration tests...${NC}"
    pytest tests/test_neurondb/ -v -m neurondb
elif [ "$1" == "coverage" ]; then
    echo -e "${GREEN}Running tests with coverage...${NC}"
    pytest tests/ --cov=NeuronAgent --cov-report=html --cov-report=term
elif [ "$1" == "fast" ]; then
    echo -e "${GREEN}Running fast tests only...${NC}"
    pytest tests/ -v -m "not slow"
else
    echo "Usage: $0 [all|api|tools|neurondb|coverage|fast]"
    exit 1
fi

echo ""
echo -e "${GREEN}Test run completed!${NC}"


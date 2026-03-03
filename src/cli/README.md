# NeuronAgent CLI

Comprehensive command-line tool for creating, managing, and testing AI agents with NeuronAgent.

## Features

- **Interactive Wizard** - Step-by-step agent creation
- **Template Support** - Pre-built agent templates
- **Workflow Management** - Create and validate workflows
- **Agent Testing** - Test agents interactively
- **Comprehensive Management** - List, show, update, delete, clone agents

## Installation

Build from source:

```bash
cd NeuronAgent/cli
go build -o neuronagent-cli
```

Or install globally:

```bash
go install github.com/neurondb/NeuronAgent/cli
```

## Quick Start

### Prerequisites

- NeuronAgent server running (default: http://localhost:8080)
- API key from NeuronAgent

Set environment variables:

```bash
export NEURONAGENT_URL=http://localhost:8080
export NEURONAGENT_API_KEY=your_api_key_here
```

### Interactive Creation

The easiest way to create an agent:

```bash
neuronagent-cli create --interactive
```

### Create from Template

```bash
# List available templates
neuronagent-cli template list

# Deploy a template
neuronagent-cli template deploy customer-support --name my-support-bot
```

### Quick Create

```bash
neuronagent-cli create \
  --name my-agent \
  --profile research-assistant \
  --tools sql,http \
  --model gpt-4
```

### Create from Config File

```bash
neuronagent-cli create --config agent-config.yaml
```

## Commands

### Create

Create a new agent using various methods:

```bash
# Interactive wizard
neuronagent-cli create --interactive

# Quick create with flags
neuronagent-cli create --name "my-agent" \
  --profile research-assistant \
  --tools sql,http \
  --model gpt-4

# From template
neuronagent-cli create --template customer-support

# From config file
neuronagent-cli create --config agent.yaml

# With workflow
neuronagent-cli create --name "data-pipeline" \
  --workflow workflow.yaml
```

### Template

Manage agent templates:

```bash
# List templates
neuronagent-cli template list

# Show template details
neuronagent-cli template show customer-support

# Search templates
neuronagent-cli template search "data"

# Deploy from template
neuronagent-cli template deploy customer-support \
  --name "my-support-bot"

# Save agent as template
neuronagent-cli template save <agent-id> \
  --name "my-custom-template"
```

### Workflow

Manage workflows:

```bash
# Validate workflow
neuronagent-cli workflow validate workflow.yaml

# List workflow templates
neuronagent-cli workflow templates

# Create workflow (when API supports it)
neuronagent-cli workflow create workflow.yaml
```

### Test

Test agents:

```bash
# Interactive test
neuronagent-cli test <agent-id>

# Single message test
neuronagent-cli test <agent-id> --message "Hello"

# Test with debug output
neuronagent-cli test <agent-id> --workflow --debug

# Validate config before creating
neuronagent-cli test --config agent.yaml --dry-run
```

### Management

Manage existing agents:

```bash
# List all agents
neuronagent-cli list

# Show agent details
neuronagent-cli show <agent-id>

# Update agent
neuronagent-cli update <agent-id> --config agent.yaml

# Delete agent
neuronagent-cli delete <agent-id>

# Clone agent
neuronagent-cli clone <agent-id> --name "new-agent"
```

## Configuration Files

### Agent Config (YAML)

Example `agent-config.yaml`:

```yaml
name: my-workflow-agent
description: A multi-step workflow agent
profile: workflow-agent

model:
  name: gpt-4
  temperature: 0.7
  max_tokens: 2000

tools:
  - sql
  - http
  - browser

config:
  temperature: 0.7
  max_tokens: 2000

memory:
  enabled: true
  hierarchical: true
  retention_days: 30

workflow:
  name: data-processing-workflow
  type: dag
  steps:
    - id: fetch_data
      name: Fetch Data
      type: tool
      tool: sql
      config:
        query: "SELECT * FROM customers LIMIT 100"
```

### Workflow Config (YAML)

Example `workflow.yaml`:

```yaml
name: data-processing-pipeline
description: Multi-step data processing workflow

steps:
  - id: extract
    name: Extract Data
    type: sql
    config:
      query: "SELECT * FROM source_table WHERE date > NOW() - INTERVAL '1 day'"
  
  - id: transform
    name: Transform Data
    type: code
    depends_on: [extract]
    config:
      language: python
      code: |
        import pandas as pd
        data = pd.DataFrame(input_data)
        return data.to_dict()
  
  - id: load
    name: Load to Database
    type: sql
    depends_on: [transform]
    config:
      query: "INSERT INTO target_table VALUES (...)"
```

## Available Templates

- **customer-support** - Multi-tier customer support agent
- **data-pipeline** - Data processing pipeline workflow
- **research-assistant** - Multi-source research assistant
- **document-qa** - RAG-based document Q&A agent
- **report-generator** - Automated report generation

## Examples

### Example 1: Create Customer Support Agent

```bash
neuronagent-cli template deploy customer-support \
  --name "acme-support"
```

### Example 2: Create Custom Agent

```bash
neuronagent-cli create --interactive
```

### Example 3: Test Agent

```bash
# List agents to find ID
neuronagent-cli list

# Test agent
neuronagent-cli test <agent-id>
```

### Example 4: Clone and Modify

```bash
# Clone existing agent
neuronagent-cli clone <agent-id> --name "my-copy"

# Update with new config
neuronagent-cli update <agent-id> --config updated-config.yaml
```

## Environment Variables

- `NEURONAGENT_URL` - NeuronAgent API URL (default: http://localhost:8080)
- `NEURONAGENT_API_KEY` - API key for authentication (required)

## Global Flags

- `--url` - Override API URL
- `--key` - Override API key
- `--format` - Output format (text, json)

## Troubleshooting

### Connection Issues

```bash
# Test connection
curl -H "Authorization: Bearer $NEURONAGENT_API_KEY" \
  $NEURONAGENT_URL/health
```

### Template Not Found

Templates are loaded from:
1. `cli/templates/` directory (relative to binary)
2. `templates/` directory (in current directory)

Ensure templates are in the correct location.

### Validation Errors

Workflow validation checks for:
- Required fields
- Valid step IDs
- Dependency cycles
- Valid step types

Use `--dry-run` to validate before creating.

## See Also

- [NeuronAgent Documentation](../../README.md)
- [API Documentation](../../docs/API.md)
- [Examples](../../examples/README.md)




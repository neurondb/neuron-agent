# NeuronAgent

<div align="center">

**AI agent runtime system providing REST API and WebSocket endpoints for building applications with long-term memory and tool execution**

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-blue.svg)](https://www.postgresql.org/)
[![Version](https://img.shields.io/badge/version-3.0.0--devel-blue.svg)](https://github.com/neurondb/neurondb)
[![License](https://img.shields.io/badge/License-Proprietary-red.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-neurondb.ai-brightgreen.svg)](https://www.neurondb.ai/docs/neuronagent)

</div>

## Table of Contents

<details>
<summary><strong>Expand full table of contents</strong></summary>

- [Overview](#overview)
  - [Key Capabilities](#key-capabilities)
- [Documentation](#documentation)
- [Features](#features)
- [Architecture](#architecture)
  - [System Architecture](#system-architecture)
  - [Agent Execution Flow](#agent-execution-flow)
  - [Workflow Engine (DAG)](#workflow-engine-dag)
  - [Memory Hierarchy](#memory-hierarchy)
- [Quick Start](#quick-start)
  - [Prerequisites](#prerequisites)
  - [Database Setup](#database-setup)
  - [Configuration](#configuration)
  - [Run Service](#run-service)
  - [Verify Installation](#verify-installation)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration-1)
- [Usage Examples](#usage-examples)
- [Documentation](#documentation-1)
- [System Requirements](#system-requirements)
- [Integration with NeuronDB](#integration-with-neurondb)
- [Security](#security)
- [Troubleshooting](#troubleshooting)
- [Support](#support)
- [License](#license)

</details>

---

## Overview

NeuronAgent integrates with NeuronDB PostgreSQL extension to provide agent runtime capabilities. Use it to build autonomous agent systems with persistent memory, tool execution, and streaming responses.

### Key Capabilities

- **Agents** — State machine, planning, tool execution, persistent memory
- **Memory** — Hierarchical (STM/MTM/LPM), HNSW vector search, promotion, feedback
- **Tools** — 16+ built-in tool types (SQL, HTTP, code, shell, browser, filesystem, memory, collaboration, etc.); NeuronDB tools (ML, vector, RAG) when NeuronDB client is configured. Custom tools from DB. See [FEATURES.md](FEATURES.md).
- **Workflows** — DAG workflows with agent/tool/HTTP/SQL/approval steps, retries, scheduling
- **Collaboration** — Workspaces, multi-agent, human-in-the-loop (approval gates, notifications)
- **Cost** — Per-agent and per-session budgets, tracking, alerts
- **API** — REST CRUD, WebSocket streaming, API key auth, RBAC

## Documentation

**[https://www.neurondb.ai/docs/neuronagent](https://www.neurondb.ai/docs/neuronagent)** — REST API reference, WebSocket guide, and deployment.

### Local documentation

- **[FEATURES.md](FEATURES.md)** — Accurate feature list (full / partial / not present)
- **[API Reference](docs/api-reference.md)** — REST API and configuration
- **[Troubleshooting](docs/troubleshooting.md)** — Common issues and fixes
- **[Release checklist](docs/release_checklist.md)** — Release and rollback

Optional references (when present): `docs/product_gap_report.txt`, `docs/config_env_schema.txt`, `docs/workflow_templates/README.md`.

## Features

<details>
<summary><strong>Complete Feature List</strong></summary>

| Feature | Description | Status |
|:--------|:------------|:-------|
| **Agent Runtime** | Complete state machine for autonomous task execution with persistent memory | Stable |
| **Multi-Agent Collaboration** | Agent-to-agent communication, task delegation, shared workspaces, and hierarchical agent structures | Stable |
| **Workflow Engine** | DAG-based workflow execution with agent, tool, HTTP, approval, and conditional steps | Stable |
| **Human-in-the-Loop (HITL)** | Approval gates, feedback loops, and human oversight in workflows with email/webhook notifications | Stable |
| **Hierarchical Memory** | Multi-level memory organization with HNSW-based vector search for better context retrieval | Stable |
| **Long-term Memory** | HNSW-based vector search for context retrieval with memory promotion | Stable |
| **Agentic RAG** | Intelligent retrieval where agent decides when and where to retrieve information | Stable |
| **Agent Memory** | Read/write memory with learning from interactions and personalization | Stable |
| **Memory Feedback** | User feedback system to improve memory quality over time | Stable |
| **Adaptive Memory** | Usage-based importance adjustment, consolidation, and compression | Stable |
| **Cross-Session Memory** | Share memories across sessions with automatic relevance detection | Stable |
| **Planning & Reflection** | LLM-based planning with task decomposition, agent self-reflection, and quality assessment | Stable |
| **Evaluation Framework** | Built-in evaluation system for agent performance with automated quality scoring | Stable |
| **Budget & Cost Management** | Real-time cost tracking, per-agent and per-session budget controls, and budget alerts | Stable |
| **Tool System** | 16+ built-in tool types (SQL, HTTP, code, shell, browser, filesystem, memory, collaboration, multimodal, web search). NeuronDB tools (ML, vector, RAG, etc.) when client configured. Custom tools from DB. See [FEATURES.md](FEATURES.md). | Stable |
| **REST API** | Full CRUD API for agents, sessions, messages, workflows, plans, budgets, and collaborations | Stable |
| **WebSocket Support** | Streaming agent responses in real-time with event streaming | Stable |
| **Authentication & Security** | API key-based authentication with bcrypt hashing, RBAC, fine-grained permissions, and audit logging | Stable |
| **Background Jobs** | PostgreSQL-based job queue with worker pool, async task execution, and memory promotion | Stable |
| **Observability** | Prometheus metrics, structured logging, distributed tracing, and debugging tools | Stable |
| **NeuronDB Integration** | Direct integration with NeuronDB embedding, LLM, vector search, and ML functions | Stable |
| **Virtual Filesystem** | Isolated filesystem for agents with secure file operations | Stable |
| **Versioning & History** | Version control for agents, execution replay, and state snapshots | Stable |

</details>

## Architecture

### System Architecture

```mermaid
graph TB
    subgraph API["API Layer"]
        REST[REST API<br/>Port 8080]
        WS[WebSocket<br/>Real-time Streaming]
        HEALTH[Health Check<br/>/health, /metrics]
    end
    
    subgraph CORE["Core Services"]
        STATE[Agent State Machine<br/>Task Execution]
        SESSION[Session Management<br/>Conversation Context]
        MEMORY[Memory Store<br/>HNSW Vector Search]
        TOOLS[Tool Registry<br/>18+ Tools]
        WORKFLOW[Workflow Engine<br/>DAG Execution]
    end
    
    subgraph WORKERS["Background Workers"]
        JOBQ[Job Queue<br/>PostgreSQL-based]
        PROMOTER[Memory Promoter<br/>Long-term Storage]
        VERIFIER[Verifier Worker<br/>Quality Checks]
    end
    
    subgraph DB["NeuronDB PostgreSQL"]
        VECTOR[Vector Search<br/>HNSW Indexes]
        EMBED[Embeddings<br/>Text/Image/Multimodal]
        LLM[LLM Integration<br/>OpenAI/Anthropic]
        ML[ML Functions<br/>25+ algorithm families]
    end
    
    REST --> STATE
    WS --> STATE
    STATE --> SESSION
    STATE --> MEMORY
    STATE --> TOOLS
    STATE --> WORKFLOW
    MEMORY --> VECTOR
    TOOLS --> DB
    WORKFLOW --> TOOLS
    JOBQ --> PROMOTER
    JOBQ --> VERIFIER
    PROMOTER --> MEMORY
    MEMORY --> VECTOR
    
    style API fill:#e3f2fd
    style CORE fill:#fff3e0
    style WORKERS fill:#f3e5f5
    style DB fill:#e8f5e9
```

### Agent Execution Flow

```mermaid
sequenceDiagram
    participant Client
    participant API as REST/WebSocket API
    participant Agent as Agent Runtime
    participant Memory as Memory Store
    participant Tools as Tool Registry
    participant DB as NeuronDB
    
    Client->>API: POST /api/v1/sessions/{id}/messages
    API->>Agent: Process message
    
    Agent->>Memory: Retrieve context (vector search)
    Memory->>DB: Query similar memories
    DB-->>Memory: Return relevant context
    Memory-->>Agent: Context + history
    
    Agent->>Agent: Generate response plan
    Agent->>Tools: Execute tool (if needed)
    Tools->>DB: Execute SQL/HTTP/Code
    DB-->>Tools: Tool results
    Tools-->>Agent: Tool output
    
    Agent->>Memory: Store new memory
    Memory->>DB: Store with embedding
    Agent->>API: Stream response
    API-->>Client: Real-time updates
```

### Workflow Engine (DAG)

```mermaid
graph LR
    START([Start]) --> STEP1[Agent Step<br/>LLM Processing]
    STEP1 --> COND{Conditional<br/>Check}
    COND -->|True| STEP2[Tool Step<br/>SQL Query]
    COND -->|False| STEP3[HTTP Step<br/>API Call]
    STEP2 --> APPROVAL{Approval<br/>Gate}
    STEP3 --> APPROVAL
    APPROVAL -->|Approved| STEP4[Agent Step<br/>Final Response]
    APPROVAL -->|Rejected| STEP1
    STEP4 --> END([End])
    
    style START fill:#c8e6c9
    style END fill:#ffcdd2
    style COND fill:#fff9c4
    style APPROVAL fill:#f8bbd0
```

### Memory Hierarchy

```mermaid
graph TD
    subgraph MEM["Memory System"]
        EPISODIC[Episodic Memory<br/>Recent conversations<br/>Short-term]
        SEMANTIC[Semantic Memory<br/>Vector embeddings<br/>HNSW search]
        WORKING[Working Memory<br/>Current session<br/>Active context]
    end
    
    WORKING -->|Promote| EPISODIC
    EPISODIC -->|Promote| SEMANTIC
    SEMANTIC -->|Retrieve| WORKING
    
    style EPISODIC fill:#ffebee
    style SEMANTIC fill:#e8f5e9
    style WORKING fill:#e3f2fd
```

> [!NOTE]
> Memory promotion follows a hierarchical structure: working memory (current session) → episodic memory (recent conversations) → semantic memory (long-term knowledge). HNSW indexes enable fast similarity search across all memory levels.

## Quick Start

### Run checklist

1. **Prerequisites:** PostgreSQL 16+ with NeuronDB extension, port 8080 free. Install NeuronDB first: [neurondb](https://github.com/neurondb/neurondb).
2. **Database:** Run migrations: `psql -d neurondb -f sql/neuron-agent.sql` or `./scripts/neuronagent-migrate.sh` (run all SQL in `sql/`).
3. **Config:** Set required env: `DB_HOST`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`. Optional: `DB_PORT`, `SERVER_PORT`, `CONFIG_PATH`. See table below.
4. **Run:** `make build` then `./bin/neuron-agent`, or `go run ./src/cmd/agent-server/main.go`. Optional: `./scripts/neuronagent-setup.sh` for DB setup; `./scripts/neuronagent-run.sh` to run.
5. **Verify:** `curl -s http://localhost:8080/health` → `{"status":"ok"}`. Generate an API key with `./scripts/neuronagent-generate-keys.sh` for `/api/v1/*`.

### Prerequisites

<details>
<summary><strong>Prerequisites Checklist</strong></summary>

- [ ] PostgreSQL 16 or later installed
- [ ] NeuronDB extension installed and enabled
- [ ] Go 1.24 or later (for building from source)
- [ ] Port 8080 available (configurable)
- [ ] API key generated (for authentication)

</details>

### Database Setup

**Option 1: Using Docker Compose (recommended for quick start)**

If you have a PostgreSQL instance with NeuronDB (e.g. from the [neurondb](https://github.com/neurondb/neurondb) repo Docker setup):

```bash
# From this repository root, run migrations against your database
psql "postgresql://neurondb:neurondb@localhost:5433/neurondb" -f sql/neuron-agent.sql
```

**Option 2: Native PostgreSQL Installation**

```bash
createdb neurondb
psql -d neurondb -c "CREATE EXTENSION neurondb;"

# Run migrations
psql -d neurondb -f sql/neuron-agent.sql
```

### Configuration

Set environment variables or create `config.yaml`:

**For Docker Compose setup (default):**
```bash
export DB_HOST=neurondb  # Service name in Docker network
export DB_PORT=5432       # Container port (not host port)
export DB_NAME=neurondb
export DB_USER=neurondb
export DB_PASSWORD=neurondb
export SERVER_PORT=8080
```

**For native PostgreSQL or connecting from host:**
```bash
export DB_HOST=localhost
export DB_PORT=5433       # Host port (Docker Compose default)
export DB_NAME=neurondb
export DB_USER=neurondb
export DB_PASSWORD=neurondb
export SERVER_PORT=8080
```

See [API Reference](docs/api-reference.md) for complete configuration options.

### Run Service

#### Automated Setup (recommended)

Use the setup script:

```bash
# From repository root
./scripts/neuronagent-setup.sh

# With system service enabled (if supported)
./scripts/neuronagent-setup.sh --enable-service
```

To generate API keys:

```bash
./scripts/neuronagent-generate-keys.sh
```

#### Manual Build and Run

From repository root:

```bash
go run ./src/cmd/agent-server/main.go
```

Or build and run:

```bash
make build
./bin/neuron-agent
```

#### Using Docker

**Option 1: Full stack (when running from neurondb repo)**

If you are in a setup that uses the neurondb repository’s Docker Compose (which can include neuronagent):

```bash
docker compose up -d neurondb
docker compose up -d neuronagent
docker compose ps neuronagent
docker compose logs -f neuronagent
```

**Option 2: NeuronAgent Docker image (no compose in this repo)**

Build and run the agent image. You must have a running PostgreSQL instance with NeuronDB (e.g. from the [neurondb](https://github.com/neurondb/neurondb) repo):

```bash
make docker-build
docker run --rm -e DB_HOST=host.docker.internal -e DB_PORT=5433 -e DB_NAME=neurondb -e DB_USER=neurondb -e DB_PASSWORD=neurondb -p 8080:8080 neuronagent:latest
```

On Linux use the actual DB host (e.g. `-e DB_HOST=172.17.0.1` or the service name if on the same Docker network). For a full stack including Postgres and NeuronSQL, use `make docker-up` (runs `docker-compose.neuronsql.yml`) or the neurondb repository.

#### Running as a Service

For systemd (Linux) or launchd (macOS), see your system documentation or the [neurondb installation services guide](https://github.com/neurondb/neurondb/blob/main/docs/getting-started/installation-services.md) for patterns.

### Verify Installation

Test health endpoint (no authentication required):

```bash
curl -s http://localhost:8080/health
```

**Expected output:**
```json
{"status":"ok"}
```

Test API with authentication:

```bash
# Replace YOUR_API_KEY with actual API key
curl -s -H "Authorization: Bearer YOUR_API_KEY" \
  http://localhost:8080/api/v1/agents | jq .
```

**Expected output:**
```json
[]
```

(Empty array if no agents created yet)

**Create your first agent:**
```bash
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-first-agent",
    "system_prompt": "You are a helpful assistant",
    "model_name": "gpt-4",
    "enabled_tools": [],
    "config": {}
  }' | jq .
```

**Expected output:**
```json
{
  "id": "agent_123",
  "name": "my-first-agent",
  "system_prompt": "You are a helpful assistant",
  "model_name": "gpt-4",
  "enabled_tools": [],
  "created_at": "2024-01-01T00:00:00Z"
}
```

> [!SUCCESS]
> **Agent created!** You can now create a session and start chatting with your agent.

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check endpoint |
| `/metrics` | GET | Prometheus metrics |
| `/api/v1/agents` | POST | Create new agent |
| `/api/v1/agents` | GET | List all agents |
| `/api/v1/agents/{id}` | GET | Get agent details |
| `/api/v1/agents/{id}` | PUT | Update agent |
| `/api/v1/agents/{id}` | DELETE | Delete agent |
| `/api/v1/sessions` | POST | Create new session |
| `/api/v1/sessions/{id}/messages` | POST | Send message to agent |
| `/ws` | WebSocket | Streaming agent responses |

See [API Reference](docs/api-reference.md) for complete API reference.

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | Database hostname |
| `DB_PORT` | `5432` | Database port |
| `DB_NAME` | `neurondb` | Database name |
| `DB_USER` | `neurondb` | Database username |
| `DB_PASSWORD` | `neurondb` | Database password |
| `DB_MAX_OPEN_CONNS` | `25` | Maximum open connections |
| `DB_MAX_IDLE_CONNS` | `5` | Maximum idle connections |
| `DB_CONN_MAX_LIFETIME` | `5m` | Connection max lifetime |
| `SERVER_HOST` | `0.0.0.0` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |
| `SERVER_READ_TIMEOUT` | `30s` | Read timeout |
| `SERVER_WRITE_TIMEOUT` | `30s` | Write timeout |
| `CONFIG_PATH` | - | Path to YAML config file (optional) |
| `AUTH_API_KEY_HEADER` | - | Custom header name for API key (optional) |
| `CORS_ALLOWED_ORIGINS` | - | Comma-separated origins for CORS |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `LOG_FORMAT` | `json` | Log format (json, text) |
| `MODULE_NEURONSQL_ENABLED` | - | Set to `true` to enable NeuronSQL module |

Other optional env: `SERVER_READ_TIMEOUT`, `SERVER_WRITE_TIMEOUT`, `WEBSOCKET_ALLOWED_ORIGINS`, `DISTRIBUTED_*`, `CACHE_*`, `TOOLS_TIMEOUT`, `WORKFLOW_MAX_DURATION`. See [docs/api-reference.md](docs/api-reference.md) and config in `src/internal/config`.

### Configuration File

Create `config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  name: neurondb
  user: neurondb
  password: neurondb
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

server:
  host: 0.0.0.0
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

logging:
  level: info
  format: json
```

Environment variables override configuration file values.

## Usage Examples

<details>
<summary><strong>Complete Usage Examples</strong></summary>

### Create Agent

```bash
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "research_agent",
    "system_prompt": "You are a research assistant that helps find and analyze information.",
    "model_name": "gpt-4",
    "profile": "research",
    "enabled_tools": ["sql", "http", "memory"],
    "config": {
      "temperature": 0.7,
      "max_tokens": 2000
    }
  }' | jq .
```

**Expected output:**
```json
{
  "id": "agent_research_001",
  "name": "research_agent",
  "system_prompt": "You are a research assistant...",
  "model_name": "gpt-4",
  "enabled_tools": ["sql", "http", "memory"],
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Create Session

```bash
curl -X POST http://localhost:8080/api/v1/sessions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "agent_research_001",
    "metadata": {
      "user_id": "user_123",
      "context": "research_project"
    }
  }' | jq .
```

**Expected output:**
```json
{
  "id": "session_abc123",
  "agent_id": "agent_research_001",
  "created_at": "2024-01-01T00:00:00Z",
  "metadata": {
    "user_id": "user_123",
    "context": "research_project"
  }
}
```

### Send Message

```bash
curl -X POST http://localhost:8080/api/v1/sessions/session_abc123/messages \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Find documents about machine learning",
    "metadata": {
      "priority": "high"
    }
  }' | jq .
```

**Expected output:**
```json
{
  "id": "msg_xyz789",
  "session_id": "session_abc123",
  "role": "user",
  "content": "Find documents about machine learning",
  "created_at": "2024-01-01T00:01:00Z"
}
```

### WebSocket Connection

Connect to WebSocket endpoint for streaming responses:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws?session_id=session_abc123');

ws.onopen = () => {
  console.log('WebSocket connected');
  // Send message
  ws.send(JSON.stringify({
    type: 'message',
    content: 'Find documents about machine learning'
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Agent response:', data);
  
  if (data.type === 'token') {
    process.stdout.write(data.content); // Stream tokens
  } else if (data.type === 'complete') {
    console.log('\nResponse complete');
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket closed');
};
```

### Python Example

```python
import requests
import json

API_BASE = "http://localhost:8080/api/v1"
API_KEY = "YOUR_API_KEY"
HEADERS = {
    "Authorization": f"Bearer {API_KEY}",
    "Content-Type": "application/json"
}

# Create agent
agent_data = {
    "name": "python_agent",
    "system_prompt": "You are a helpful Python assistant",
    "model_name": "gpt-4",
    "enabled_tools": ["sql", "code"]
}

response = requests.post(
    f"{API_BASE}/agents",
    headers=HEADERS,
    json=agent_data
)
agent = response.json()
print(f"Created agent: {agent['id']}")

# Create session
session_data = {"agent_id": agent['id']}
response = requests.post(
    f"{API_BASE}/sessions",
    headers=HEADERS,
    json=session_data
)
session = response.json()
print(f"Created session: {session['id']}")

# Send message
message_data = {
    "content": "Write a Python function to calculate fibonacci numbers"
}
response = requests.post(
    f"{API_BASE}/sessions/{session['id']}/messages",
    headers=HEADERS,
    json=message_data
)
message = response.json()
print(f"Message sent: {message['id']}")
```

</details>

## Documentation

| Document | Description |
|----------|-------------|
| [FEATURES.md](FEATURES.md) | Feature list (full / partial / not present) |
| [API Reference](docs/api-reference.md) | REST API and configuration |
| [Troubleshooting](docs/troubleshooting.md) | Common issues and solutions |
| [src/cli/README.md](src/cli/README.md) | Command-line client (`neuronagent-cli`): build from `src/cli`, flags `--url`, `--key`, `--format`; subcommands for create, template, workflow, test, list, show, update, delete, clone, neuronsql |

## System Requirements

| Component | Requirement |
|-----------|-------------|
| PostgreSQL | 16 or later |
| NeuronDB Extension | Installed and enabled ([install](https://github.com/neurondb/neurondb)) |
| Go | 1.24 or later (for building) |
| Network | Port 8080 available (configurable) |

Related: [NeuronDB](https://github.com/neurondb/neurondb) (extension), [NeuronMCP](https://github.com/neurondb/neuron-mcp) (MCP server).

## Integration with NeuronDB

NeuronAgent requires PostgreSQL 16+ with the NeuronDB extension. Install NeuronDB first: [neurondb repository](https://github.com/neurondb/neurondb) ([Simple Start](https://github.com/neurondb/neurondb/blob/main/docs/getting-started/simple-start.md)). Full-stack deployment (NeuronDB + NeuronAgent + NeuronMCP) is documented in each component’s repository.

## Security

- API key authentication required for all API endpoints
- Rate limiting configured per API key
- Database credentials stored securely via environment variables
- Supports TLS/SSL for encrypted connections
- Non-root user in Docker containers

See [API Reference](docs/api-reference.md) for security and configuration details.

## Troubleshooting

### Service Won't Start

Check database connection:

```bash
psql -h localhost -p 5432 -U neurondb -d neurondb -c "SELECT 1;"
```

Verify environment variables:

```bash
env | grep -E "DB_|SERVER_"
```

Check logs:

```bash
docker compose logs agent-server
```

### Database Connection Failed

Verify NeuronDB extension:

```sql
SELECT * FROM pg_extension WHERE extname = 'neurondb';
```

Check database permissions:

```sql
GRANT ALL PRIVILEGES ON DATABASE neurondb TO neurondb;
GRANT ALL ON SCHEMA neurondb_agent TO neurondb;
```

### API Not Responding

Test health endpoint:

```bash
curl http://localhost:8080/health
```

Verify API key:

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  http://localhost:8080/api/v1/agents
```

## Support

- **Documentation**: This README and the [docs](docs/) directory
- **GitHub Issues**: [Report issues](https://github.com/neurondb/neurondb/issues)
- **Email**: support@neurondb.ai

## License

See [LICENSE](LICENSE) for license information.

---

[Back to top](#neuronagent)

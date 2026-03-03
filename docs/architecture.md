# NeuronAgent Architecture

<div align="center">

**AI agent runtime system architecture and design**

[![Architecture](https://img.shields.io/badge/architecture-complete-brightgreen)](.)
[![Status](https://img.shields.io/badge/status-stable-blue)](.)

</div>

---

> [!NOTE]
> NeuronAgent integrates with the NeuronDB PostgreSQL extension. It provides agent capabilities including long-term memory, tool execution, planning, reflection, and advanced features.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    NeuronAgent Service                       │
├─────────────────────────────────────────────────────────────┤
│  REST API Layer        │  WebSocket        │  Health/Metrics │
├─────────────────────────────────────────────────────────────┤
│  Agent Runtime Engine  │  Session Manager  │  Tool Registry  │
├─────────────────────────────────────────────────────────────┤
│  Memory Manager        │  Planner          │  Reflector      │
├─────────────────────────────────────────────────────────────┤
│  Hierarchical Memory   │  Event Stream     │  VFS            │
├─────────────────────────────────────────────────────────────┤
│  Collaboration         │  Async Tasks      │  Sub-Agents     │
├─────────────────────────────────────────────────────────────┤
│  Verification Agent    │  Multimodal Proc  │  Browser Driver │
├─────────────────────────────────────────────────────────────┤
│  Background Workers    │  Job Queue        │  Notifications  │
├─────────────────────────────────────────────────────────────┤
│              NeuronDB PostgreSQL Extension                   │
│  (Vector Search │ Embeddings │ LLM │ ML │ RAG)              │
└─────────────────────────────────────────────────────────────┘
```

## Core Components

### Database Layer (`internal/db/`)

Manages all database interactions with prepared statements and connection pooling.

**Key Components:**
- **Models**: Go structs representing database entities (Agent, Session, Message, Tool, MemoryChunk, etc.)
- **Connection**: Connection pool management with configurable limits and health checks
- **Queries**: All SQL queries use prepared statements for security and performance
- **Transactions**: Transaction management for complex operations
- **Migrations**: 20 migration files tracking schema evolution

### Agent Runtime (`internal/agent/`)

Core execution engine orchestrating agent behavior.

**Runtime Components:**
- **Runtime**: Main execution engine with state machine
- **MemoryManager**: HNSW-based vector search for long-term memory
- **HierarchicalMemoryManager**: Three-tier memory architecture (working, episodic, semantic)
- **EventStreamManager**: Event logging with automatic summarization
- **LLMClient**: Integration with NeuronDB LLM functions
- **ContextLoader**: Context loading combining messages and memory
- **PromptBuilder**: Prompt construction with templating
- **Planner**: Task planning and decomposition
- **Reflector**: Self-reflection and improvement
- **VerificationAgent**: Quality assurance through configurable rules
- **VirtualFileSystem**: Hybrid storage (database + S3) with atomic operations
- **AsyncTaskExecutor**: Asynchronous task execution with notifications
- **SubAgentManager**: Sub-agent routing and delegation
- **TaskNotifier**: Alert and notification management
- **EnhancedMultimodalProcessor**: Image, audio, and code processing

### Tools System (`internal/tools/`)

Extensible tool registry with built-in and custom tools. See [features.md](features.md) for the full tool list.

### API Layer (`internal/api/`)

REST API and WebSocket endpoints. See [api.md](api.md) and [api-reference.md](api-reference.md).

### Authentication (`internal/auth/`)

Security and access control: APIKeyManager, Hasher, Validator, RateLimiter, PrincipalManager, RBAC, AuditLogger.

### Background Jobs (`internal/jobs/`)

PostgreSQL-based job queue system with MemoryPromoter, VerifierWorker, AsyncTaskWorker, and cleanup workers.

## Integration Points

### NeuronDB Extension
- Vector operations and search
- Embedding generation
- LLM function calls
- ML operations
- RAG pipelines

### External Services
- Email providers (SMTP)
- Webhook endpoints
- S3 storage
- Custom tool integrations

## Configuration

Configuration is managed through:
- Environment variables (preferred for production)
- Configuration file (YAML format)
- Default values for optional settings

## Monitoring and Observability

- **Metrics**: Prometheus metrics endpoint
- **Logging**: Structured JSON logging with request ID tracking
- **Tracing**: Distributed tracing support
- **Health Checks**: Database and component health status

---

See [overview.md](overview.md) and [api-reference.md](api-reference.md) for more details.

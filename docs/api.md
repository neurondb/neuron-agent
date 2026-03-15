# NeuronAgent API Documentation

> The canonical REST API reference is **[API Reference](api-reference.md)**. This document groups endpoints by area. Only endpoints **currently registered** on the server are listed; for request/response schemas and examples see the API reference and the [OpenAPI spec](../openapi/openapi.yaml).

## Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
- [Endpoints](#endpoints)
  - [Agents](#agents)
  - [Sessions](#sessions)
  - [Messages](#messages)
  - [Tools](#tools)
  - [Memory](#memory)
  - [Runs](#runs)
  - [Budgets](#budgets)
  - [Approvals](#approvals)
  - [Feedback](#feedback)
  - [Analytics](#analytics)
  - [Batch](#batch)
  - [Workflows](#workflows)
  - [Governance and Admin](#governance-and-admin)
  - [LLM SQL and RAG](#llm-sql-and-rag)
  - [WebSocket](#websocket)
- [Not Currently Exposed](#not-currently-exposed)
- [Error Handling and Rate Limiting](#error-handling-and-rate-limiting)

## Base URL

```
http://localhost:8080/api/v1
```

## OpenAPI Specification

See the [OpenAPI 3.0 spec](../openapi/openapi.yaml) for request/response schemas and to generate clients.

## Authentication

Generate keys with `./scripts/neuronagent-generate-keys.sh` or the `generate-key` binary.

```
Authorization: Bearer <api_key>
```

## Endpoints

### Agents

#### Create Agent
```
POST /api/v1/agents
```

Request body:
```json
{
  "name": "my-agent",
  "description": "A helpful agent",
  "system_prompt": "You are a helpful assistant.",
  "model_name": "gpt-4",
  "enabled_tools": ["sql", "http"],
  "config": {
    "temperature": 0.7,
    "max_tokens": 1000
  }
}
```

#### List Agents
```
GET /api/v1/agents
```

#### Get Agent
```
GET /api/v1/agents/{id}
```

#### Update Agent
```
PUT /api/v1/agents/{id}
```

#### Delete Agent
```
DELETE /api/v1/agents/{id}
```

#### Clone Agent
```
POST /api/v1/agents/{id}/clone
```

#### Generate Plan
```
POST /api/v1/agents/{id}/plan
```

#### Reflect on Response
```
POST /api/v1/agents/{id}/reflect
```

#### Delegate to Agent
```
POST /api/v1/agents/{id}/delegate
```

#### Get Agent Metrics
```
GET /api/v1/agents/{id}/metrics
```

#### Get Agent Costs
```
GET /api/v1/agents/{id}/costs
```

#### List Agent Versions
```
GET /api/v1/agents/{id}/versions
```
*Not currently registered on the server.*

#### Create Agent Version
*Not currently registered.*

#### Get Agent Version
*Not currently registered.*

#### Activate Agent Version
*Not currently registered.*

#### List Agent Relationships
*Not currently registered.*

#### Create Agent Relationship
*Not currently registered.*

#### Delete Agent Relationship
*Not currently registered.*

#### Batch Create Agents
```
POST /api/v1/batch/agents
```

#### Batch Delete Agents
```
POST /api/v1/batch/agents/delete
```

### Sessions

#### Create Session
```
POST /api/v1/sessions
```

Request body:
```json
{
  "agent_id": "uuid",
  "external_user_id": "user123",
  "metadata": {}
}
```

#### Get Session
```
GET /api/v1/sessions/{id}
```

#### Update Session
```
PUT /api/v1/sessions/{id}
```

#### Delete Session
```
DELETE /api/v1/sessions/{id}
```

#### List Sessions
```
GET /api/v1/agents/{agent_id}/sessions
```

### Messages

#### Send Message
```
POST /api/v1/sessions/{session_id}/messages
```

Request body:
```json
{
  "role": "user",
  "content": "Hello, how are you?",
  "stream": false
}
```

#### Get Messages
```
GET /api/v1/sessions/{session_id}/messages
```

#### Get Message
```
GET /api/v1/messages/{id}
```

#### Update Message
```
PUT /api/v1/messages/{id}
```

#### Delete Message
```
DELETE /api/v1/messages/{id}
```

### Tools

- `POST /api/v1/tools` — Create tool
- `GET /api/v1/tools` — List tools
- `GET /api/v1/tools/{id}` — Get tool
- `PUT /api/v1/tools/{id}` — Update tool
- `DELETE /api/v1/tools/{id}` — Delete tool
- `GET /api/v1/tools/{id}/analytics` — Tool analytics

#### Batch Delete Messages
```
POST /api/v1/batch/messages/delete
```

### Approvals (human-in-the-loop)

- `GET /api/v1/approvals` — List approval requests
- `GET /api/v1/approvals/{id}` — Get approval request
- `POST /api/v1/approvals/{id}/approve` — Approve
- `POST /api/v1/approvals/{id}/reject` — Reject

### Feedback

- `GET /api/v1/feedback` — List feedback
- `GET /api/v1/feedback/stats` — Feedback statistics

### WebSocket

#### Connect to WebSocket
```
WS /ws?session_id={session_id}&api_key={api_key}
```

Or use Authorization header:
```
WS /ws?session_id={session_id}
Headers: Authorization: Bearer {api_key}
```

**Features:**
- API key authentication (query parameter or header)
- Ping/pong keepalive (60s timeout)
- Message queue for concurrent requests
- Graceful error handling

**Message Format:**
```json
{
  "content": "Your message here"
}
```

**Response Format:**
```json
{
  "type": "chunk",
  "content": "Response chunk..."
}
```

```json
{
  "type": "response",
  "content": "Full response",
  "complete": true,
  "tokens_used": 150,
  "tool_calls": [],
  "tool_results": []
}
```

**Error Format:**
```json
{
  "type": "error",
  "error": "Error message"
}
```

Send messages:
```json
{
  "content": "Hello"
}
```

Receive responses:
```json
{
  "type": "response",
  "content": "Hello! How can I help you?",
  "complete": true
}
```

## Evaluation Framework

*These endpoints are not currently registered on the server.*

- `POST /api/v1/eval/tasks`
- `GET /api/v1/eval/tasks`
- `POST /api/v1/eval/runs`
- `POST /api/v1/eval/runs/{run_id}/execute`
- `GET /api/v1/eval/runs/{run_id}/results`

## Execution Snapshots and Replay

*These endpoints are not currently registered on the server.*

- `POST /api/v1/sessions/{session_id}/snapshots`
- `GET /api/v1/sessions/{session_id}/snapshots`, `GET /api/v1/agents/{agent_id}/snapshots`
- `POST /api/v1/snapshots/{id}/replay`
- `DELETE /api/v1/snapshots/{id}`

## Workflow Schedules

### Create/Update Workflow Schedule
```
POST /api/v1/workflows/{workflow_id}/schedule
```

Request body:
```json
{
  "cron_expression": "0 0 * * *",
  "timezone": "UTC",
  "enabled": true
}
```

### Get Workflow Schedule
```
GET /api/v1/workflows/{workflow_id}/schedule
```

### List Workflow Schedules
```
GET /api/v1/workflows/schedules
```

### Delete Workflow Schedule
```
DELETE /api/v1/workflows/{workflow_id}/schedule
```

## Agent Specializations

*These endpoints are not currently registered on the server.*

- `POST /api/v1/agents/{agent_id}/specialization`
- `GET /api/v1/agents/{agent_id}/specialization`
- `GET /api/v1/specializations`
- `PUT /api/v1/agents/{agent_id}/specialization`
- `DELETE /api/v1/agents/{agent_id}/specialization`

## Memory Management

### Submit Memory Feedback
```
POST /api/v1/agents/{id}/memory/feedback
```
Submit user feedback on memory retrieval. (Not `POST /api/v1/memory/{memory_id}/feedback`.)

Request body:
```json
{
  "agent_id": "uuid",
  "session_id": "uuid (optional)",
  "memory_tier": "chunk|stm|mtm|lpm",
  "feedback_type": "positive|negative|neutral|correction",
  "feedback_text": "Optional feedback text",
  "query": "Query that led to this memory retrieval (optional)",
  "relevance_score": 0.85,
  "metadata": {}
}
```

Response:
```json
{
  "feedback_id": "uuid",
  "memory_id": "uuid",
  "agent_id": "uuid",
  "feedback_type": "positive",
  "status": "recorded",
  "message": "Feedback recorded and memory quality updated",
  "duration_ms": 45,
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Get Retrieval Statistics
```
GET /api/v1/analytics/retrieval-stats
```
Retrieval statistics (not under `/api/v1/agents/{id}/retrieval-stats`).

### Consolidate Memory
```
POST /api/v1/agents/{id}/memory/consolidate
```

Consolidate similar memories to reduce duplication and improve storage efficiency.

Request body:
```json
{
  "tier": "stm|mtm|lpm",
  "similarity_threshold": 0.9
}
```

Response:
```json
{
  "agent_id": "uuid",
  "tier": "mtm",
  "similarity_threshold": 0.9,
  "consolidated_count": 15,
  "status": "completed",
  "duration_ms": 2500,
  "completed_at": "2024-01-01T00:00:00Z"
}
```

### Get Memory Quality
```
GET /api/v1/agents/{id}/memory/quality?memory_id={uuid}&tier={tier}
```

Get quality metrics for a specific memory or update all memory quality scores.

Query parameters:
- `memory_id` (optional): Specific memory ID
- `tier` (optional): Memory tier (stm, mtm, lpm) - required if memory_id provided

For complete API documentation including all endpoints, request/response schemas, and examples, see the [OpenAPI specification](../openapi/openapi.yaml).

---

<div align="center">

[⬆ Back to Top](#neuronagent-api-documentation)

</div>

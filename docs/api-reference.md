# NeuronAgent API Reference

## Overview

NeuronAgent provides a comprehensive REST API for managing AI agents, sessions, messages, tools, and advanced features. All API endpoints are versioned under `/api/v1` and require API key authentication.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All API requests require an API key in the `Authorization` header:

```
Authorization: Bearer YOUR_API_KEY
```

API keys can be generated using the `generate-key` command or through your NeuronAgent administration interface.

## API Endpoints

### Agents

#### Create Agent
`POST /api/v1/agents`

Creates a new agent with specified configuration.

**Request Body:**
```json
{
  "name": "research_agent",
  "description": "Agent for research tasks",
  "system_prompt": "You are a helpful research assistant.",
  "model_name": "gpt-4",
  "memory_table": "research_memory",
  "enabled_tools": ["sql", "http", "browser"],
  "config": {
    "temperature": 0.7,
    "max_tokens": 2000
  }
}
```

**Example (curl):**
```bash
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "research_agent",
    "description": "Agent for research tasks",
    "system_prompt": "You are a helpful research assistant.",
    "model_name": "gpt-4",
    "memory_table": "research_memory",
    "enabled_tools": ["sql", "http", "browser"],
    "config": {
      "temperature": 0.7,
      "max_tokens": 2000
    }
  }'
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "name": "research_agent",
  "description": "Agent for research tasks",
  "system_prompt": "You are a helpful research assistant.",
  "model_name": "gpt-4",
  "memory_table": "research_memory",
  "enabled_tools": ["sql", "http", "browser"],
  "config": {},
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

#### List Agents
`GET /api/v1/agents?search=query`

Lists all agents, optionally filtered by search query.

#### Get Agent
`GET /api/v1/agents/{id}`

Retrieves detailed information about a specific agent.

#### Update Agent
`PUT /api/v1/agents/{id}`

Updates an existing agent configuration.

#### Delete Agent
`DELETE /api/v1/agents/{id}`

Deletes an agent and all associated data.

#### Clone Agent
`POST /api/v1/agents/{id}/clone`

Creates a copy of an existing agent with a new ID.

#### Generate Plan
`POST /api/v1/agents/{id}/plan`

Generates an execution plan for a given task.

#### Reflect on Response
`POST /api/v1/agents/{id}/reflect`

Submits agent response for reflection and improvement.

#### Delegate to Agent
`POST /api/v1/agents/{id}/delegate`

Delegates a task to another agent or sub-agent.

#### Get Agent Metrics
`GET /api/v1/agents/{id}/metrics`

Retrieves performance metrics for an agent.

#### Get Agent Costs
`GET /api/v1/agents/{id}/costs`

Retrieves cost tracking information for an agent.

### Sessions

#### Create Session
`POST /api/v1/sessions`

Creates a new conversation session with an agent.

**Request Body:**
```json
{
  "agent_id": "agent_research_001",
  "external_user_id": "user123",
  "metadata": {
    "project": "research",
    "source": "web"
  }
}
```

**Response:** `201 Created`
```json
{
  "id": "session_abc123",
  "agent_id": "agent_research_001",
  "external_user_id": "user123",
  "metadata": {},
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

#### Get Session
`GET /api/v1/sessions/{id}`

Retrieves session details.

#### Update Session
`PUT /api/v1/sessions/{id}`

Updates session metadata.

#### Delete Session
`DELETE /api/v1/sessions/{id}`

Deletes a session and all associated messages.

#### List Sessions
`GET /api/v1/agents/{agent_id}/sessions`

Lists all sessions for an agent. Supports `limit` and `offset` for pagination.

### Messages

#### Send Message
`POST /api/v1/sessions/{session_id}/messages`

Sends a message to an agent in a session. The agent processes the message (loads context, may call tools, updates memory) and returns a response.

**Request Body:**
```json
{
  "content": "Find documents about machine learning",
  "role": "user",
  "stream": false,
  "metadata": {
    "priority": "high"
  }
}
```

- `content` (required): Message text.
- `role`: Usually `"user"`; can be omitted (defaults to user).
- `stream`: If `true`, use WebSocket or streaming response where supported.
- `metadata`: Optional key-value metadata.

**Response:** `201 Created` (or streaming)
```json
{
  "id": "msg_xyz789",
  "session_id": "session_abc123",
  "role": "assistant",
  "content": "I found the following documents about machine learning...",
  "metadata": {
    "tool_calls": [],
    "tokens_used": 150
  },
  "created_at": "2025-01-01T00:01:00Z"
}
```

For streaming, the response may be chunked or delivered via WebSocket (see [WebSocket](#websocket) below).

#### Get Messages
`GET /api/v1/sessions/{session_id}/messages`

Retrieves message history for a session. Supports `limit` and `offset`.

#### Get Message
`GET /api/v1/messages/{id}`

Retrieves a specific message.

#### Update Message
`PUT /api/v1/messages/{id}`

Updates message content or metadata.

#### Delete Message
`DELETE /api/v1/messages/{id}`

Deletes a message.

### Tools

#### List Tools
`GET /api/v1/tools`

Lists all available tools.

#### Create Tool
`POST /api/v1/tools`

Registers a new custom tool.

#### Get Tool
`GET /api/v1/tools/{name}`

Retrieves tool details.

#### Update Tool
`PUT /api/v1/tools/{name}`

Updates tool configuration.

#### Delete Tool
`DELETE /api/v1/tools/{name}`

Deletes a tool.

#### Get Tool Analytics
`GET /api/v1/tools/{name}/analytics`

Retrieves usage analytics for a tool.

### Memory

#### List Memory Chunks
`GET /api/v1/agents/{id}/memory`

Lists memory chunks for an agent.

#### Search Memory
`POST /api/v1/agents/{id}/memory/search`

Searches agent memory using vector similarity.

#### Get Memory Chunk
`GET /api/v1/memory/{chunk_id}`

Retrieves a specific memory chunk.

#### Delete Memory Chunk
`DELETE /api/v1/memory/{chunk_id}`

Deletes a memory chunk.

#### Summarize Memory
`POST /api/v1/memory/{id}/summarize`

Generates a summary of memory chunks.

### Plans, Reflections, Budgets, Webhooks, Human-in-the-Loop, Collaboration Workspaces, Async Tasks, Alert Preferences, Batch Operations, Analytics, and more are documented in [api.md](api.md).

### WebSocket

#### WebSocket Connection
`GET /ws?session_id={session_id}`

Establishes a WebSocket connection for streaming agent responses.

**Authentication:** Pass the API key via query parameter `api_key={key}` or via the `Authorization: Bearer {key}` header when establishing the connection.

**Query parameters:**
- `session_id` (required): ID of an existing session. The session must belong to an agent and exist before connecting.

**Message format (client → server):**
```json
{
  "content": "Your message here"
}
```

**Message format (server → client):** Streamed events may include:
- `type: "chunk"` — Incremental content (e.g. token stream).
- `type: "response"` — Full or final response with optional `tool_calls`, `tool_results`, `tokens_used`, `complete: true`.
- `type: "error"` — Error message.

**Features:** Ping/pong keepalive, message queue for concurrent requests, graceful error handling. See [api.md](api.md) for full WebSocket message schema.

**Example (JavaScript):**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws?session_id=' + sessionId + '&api_key=' + apiKey);
ws.onmessage = (e) => {
  const data = JSON.parse(e.data);
  if (data.type === 'chunk') process.stdout.write(data.content);
  if (data.type === 'response' && data.complete) console.log('\nDone.');
};
ws.send(JSON.stringify({ content: 'Hello, agent!' }));
```

## Error Responses

All errors follow a consistent format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "request_id": "uuid",
    "details": {}
  }
}
```

**HTTP Status Codes:**
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid API key
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service temporarily unavailable

## Rate Limiting

API requests are rate-limited per API key. When rate limits are exceeded, the API returns `429 Too Many Requests` with a `Retry-After` header.

## Pagination

List endpoints support pagination using query parameters:

- `limit` - Maximum number of items to return (default: 50, max: 1000)
- `offset` - Number of items to skip (default: 0)

## OpenAPI Specification

A complete OpenAPI 3.0 specification is available at:

```
http://localhost:8080/openapi.yaml
```

Use this specification to generate client libraries or explore the API interactively.

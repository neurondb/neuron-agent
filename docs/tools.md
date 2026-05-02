# Tools

Tools are rows in the database (`handler_type`, configuration JSON) executed via `internal/tools/registry.go`.

## Built-in handler types (baseline registry)

- `sql`, `http`, `code`, `shell`, `browser`, `visualization`, `mcp`

## MCP tools

`Registry.SyncFromMCP(ctx, baseURL)` calls JSON-RPC `tools/list` on an HTTP MCP server and upserts tools with `handler_type = mcp` and `handler_config` `{ "mcp_server_url", "tool_name" }`. Runtime execution uses JSON-RPC `tools/call`.

## NeuronDB-augmented handlers

When the DB-backed NeuronDB client is wired (`NewRegistryWithNeuronDB`): `ml`, `vector`, `rag`, `analytics`, `hybrid_search`, `reranking`.

## NeuronSQL module

Prefixed tools such as `neuronsql.schema_snapshot`, `neuronsql.explain_json`, etc., when the NeuronSQL module is enabled.

## Permissions

Tool execution respects timeouts, optional RBAC, circuit breakers, and audit logging — see code paths in `registry.go`.

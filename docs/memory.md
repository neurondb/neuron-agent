# Memory

NeuronAgent stores recall data in PostgreSQL:

- **Flat chunks** — `memory_chunks` with embeddings for similarity search (`<=>` / cosine-style ordering depending on schema).
- **Tiers** — STM / MTM / LPM tables with promotion workers moving content based on importance and access patterns.

## APIs (overview)

- Search: `POST /api/v1/agents/{id}/memory/search`
- List: `GET /api/v1/agents/{id}/memory`
- Forget / consolidate / conflicts — see `/api/v1/agents/{id}/memory/*` routes in the OpenAPI spec (`/docs`).

For diagrams and deeper behavior, see hosted docs and [`architecture.md`](architecture.md).

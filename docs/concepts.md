# Core concepts

- **Agent** — configuration + model + tools; state is stored in PostgreSQL.
- **Session** — a conversation thread for one agent; messages and tool calls are tied to a session.
- **Memory** — vector-backed chunks and hierarchical tiers (STM/MTM/LPM) for recall and promotion.
- **Tools** — registered handlers (`sql`, `http`, `rag`, NeuronSQL namespaced tools, …) invoked by the runtime with permissions and timeouts.
- **Workflow** — DAG of steps (`agent`, `tool`, `sql`, `http`, `approval`, `conditional`) with executions persisted in the database.
- **NeuronDB** — PostgreSQL extension providing embeddings, RAG helpers, hybrid search, etc.

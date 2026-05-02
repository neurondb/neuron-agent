# RAG (Retrieval-Augmented Generation)

RAG uses NeuronDB functions inside PostgreSQL for chunking, embedding, retrieval, reranking, and answer generation where configured.

- **Ingest** — `POST /api/v1/rag/ingest` (registered when the NeuronDB RAG client is available).
- **Tools** — `rag` handler is registered when using `NewRegistryWithNeuronDB` (default server path).

Requires NeuronDB extension and compatible embedding dimensions; validate with your DB image.

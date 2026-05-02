# RAG with documents

Use `POST /api/v1/rag/ingest` when NeuronDB RAG is available, then query via agent sessions or tools.

See [`tests/fixtures/sample-doc.txt`](../../tests/fixtures/sample-doc.txt) for a small text fixture.

Example ingest (requires valid agent context and API key):

```bash
curl -X POST "$NEURON_AGENT_URL/api/v1/rag/ingest" \
  -H "Authorization: Bearer $NEURON_AGENT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"title":"Fixture","content":"See tests/fixtures/sample-doc.txt"}'
```

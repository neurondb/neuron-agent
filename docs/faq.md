# FAQ

**Why PostgreSQL?**  
ACID semantics, operational familiarity, backup/replication, and NeuronDB as an extension.

**Where is the OpenAPI spec?**  
Bundled at **`GET /docs/openapi.yaml`** and interactive **`GET /docs`** (Redoc).

**Why does `/readyz` return 503?**  
Database unreachable — check `DB_*` env vars, network, and migrations.

**NeuronDB missing in Docker?**  
Stock `postgres` images do not include NeuronDB; vector/RAG features need an extension-enabled image or custom build.

**How do I generate API keys?**  
Use `src/cmd/generate-key` or `scripts/neuronagent-generate-keys.sh` against a migrated database.

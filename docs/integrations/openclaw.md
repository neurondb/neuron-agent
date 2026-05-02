# OpenClaw bridge

NeuronAgent exposes a channel-friendly surface under **`/claw/v1`**:

- `GET /claw/v1/health`
- `POST /claw/v1/tools/list`
- `POST /claw/v1/tools/run`

**OpenClaw** routes user-facing channels; **NeuronAgent** provides durable memory, RAG, SQL tooling, workflows, and audit.

Wire OpenClaw to this base URL with an API key that has appropriate roles. See [openclaw-bridge-plan.md](../roadmap/openclaw-bridge-plan.md) for the full permission matrix.

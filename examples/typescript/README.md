# NeuronSQL minimal example (TypeScript)

Run NeuronAgent with docker-compose, then call the NeuronSQL generate and Claw list endpoints.

## Prerequisites

- Docker and docker-compose
- Node 18+

## 1. Start stack

From the repository root:

```bash
cd docker && docker-compose -f docker-compose.neuronsql.yml up -d
```

Set `NEURONAGENT_URL` and `NEURONAGENT_API_KEY` (create key via generate-key CLI or DB).

## 2. Install and run

```bash
npm install
export NEURONAGENT_URL=http://localhost:8080
export NEURONAGENT_API_KEY=your-api-key
npx ts-node neuronsql_minimal.ts
```

Or with node after compiling: `tsc && node neuronsql_minimal.js`

## What the example does

- Calls `POST /api/v1/neuronsql/generate` with a question and DSN.
- Calls `POST /claw/v1/tools/list` to list neuronsql.* tools.

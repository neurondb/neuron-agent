# NeuronSQL minimal example (Python)

Run NeuronAgent with docker-compose, then call the NeuronSQL generate and Claw list endpoints.

## Prerequisites

- Docker and docker-compose
- Python 3.9+

## 1. Start stack

From the repository root:

```bash
cd docker && docker-compose -f docker-compose.neuronsql.yml up -d
```

Wait for Postgres to be healthy, then apply schema (if needed) and create an API key. Example with psql:

```bash
export PGPASSWORD=neurondb
psql -h localhost -p 5433 -U neurondb -d neurondb -f ../sql/neuron-agent.sql  # if first run
# Create an API key via the generate-key CLI or INSERT into neurondb_agent.api_keys
```

Set `NEURONAGENT_URL` and `NEURONAGENT_API_KEY` (see below).

## 2. Install and run

```bash
pip install requests
export NEURONAGENT_URL=http://localhost:8080
export NEURONAGENT_API_KEY=your-api-key
python neuronsql_minimal.py
```

## What the example does

- `neuronsql_minimal.py` calls `POST /api/v1/neuronsql/generate` with a question and DSN, then prints the generated SQL and validation result.
- Optionally calls `POST /claw/v1/tools/list` to list neuronsql.* tools.

/**
 * Minimal NeuronSQL example: generate SQL and list Claw tools.
 * Run against docker-compose. Set NEURONAGENT_URL and NEURONAGENT_API_KEY.
 */
const BASE = process.env.NEURONAGENT_URL || "http://localhost:8080";
const API_KEY = process.env.NEURONAGENT_API_KEY || "";
const DSN = process.env.NEURONAGENT_DSN || "postgres://neurondb:neurondb@localhost:5433/neurondb";

async function request(path: string, method = "GET", data?: object): Promise<any> {
  const res = await fetch(`${BASE}${path}`, {
    method,
    headers: {
      Authorization: `Bearer ${API_KEY}`,
      "Content-Type": "application/json",
    },
    body: data ? JSON.stringify(data) : undefined,
  });
  if (!res.ok) throw new Error(`${res.status} ${await res.text()}`);
  return res.json();
}

async function main() {
  if (!API_KEY) {
    console.log("Set NEURONAGENT_API_KEY");
    return;
  }
  console.log("1. NeuronSQL generate");
  try {
    const out = await request("/api/v1/neuronsql/generate", "POST", {
      db_dsn: DSN,
      question: "How many users are there?",
    });
    console.log("   SQL:", out.sql || "(none)");
    console.log("   Valid:", out.validation_report?.valid);
  } catch (e: any) {
    console.log("   Error:", e.message);
  }
  console.log("2. Claw tools list");
  try {
    const out = await request("/claw/v1/tools/list", "POST", {});
    console.log("   Tools:", out.tools || []);
  } catch (e: any) {
    console.log("   Error:", e.message);
  }
}

main();

#!/usr/bin/env python3
"""Minimal NeuronSQL example: generate SQL and list Claw tools. Run against docker-compose."""
import os
import json
import urllib.request
import urllib.error

BASE = os.environ.get("NEURONAGENT_URL", "http://localhost:8080")
API_KEY = os.environ.get("NEURONAGENT_API_KEY", "")
DSN = os.environ.get("NEURONAGENT_DSN", "postgres://neurondb:neurondb@localhost:5433/neurondb")


def request(path, method="GET", data=None):
    req = urllib.request.Request(
        f"{BASE}{path}",
        method=method,
        headers={"Authorization": f"Bearer {API_KEY}", "Content-Type": "application/json"},
    )
    if data is not None:
        req.data = json.dumps(data).encode()
    with urllib.request.urlopen(req) as resp:
        return json.load(resp)


def main():
    if not API_KEY:
        print("Set NEURONAGENT_API_KEY (e.g. from generate-key CLI)")
        return
    print("1. NeuronSQL generate")
    try:
        out = request("/api/v1/neuronsql/generate", method="POST", data={
            "db_dsn": DSN,
            "question": "How many users are there?",
        })
        print("   SQL:", out.get("sql", "(none)"))
        print("   Valid:", out.get("validation_report", {}).get("valid"))
    except urllib.error.HTTPError as e:
        print("   Error:", e.code, e.read().decode())
    print("2. Claw tools list")
    try:
        out = request("/claw/v1/tools/list", method="POST", data={})
        print("   Tools:", out.get("tools", []))
    except urllib.error.HTTPError as e:
        print("   Error:", e.code, e.read().decode())


if __name__ == "__main__":
    main()

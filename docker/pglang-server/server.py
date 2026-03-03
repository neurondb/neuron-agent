#!/usr/bin/env python3
"""
PGLang HTTP server: POST /v1/completions with {prompt, max_tokens, model} -> {text, usage}.
Stub implementation for demo; replace with memorable inference for production.
"""
import os
import sys
from contextlib import asynccontextmanager

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

# Optional: add memorable to path for real inference
MEMORABLE_PATH = os.environ.get("MEMORABLE_PATH", "")
if MEMORABLE_PATH:
    sys.path.insert(0, MEMORABLE_PATH)

STUB_RESPONSE = os.environ.get("PGLANG_STUB_RESPONSE", "SELECT 1;")


class CompletionsRequest(BaseModel):
    prompt: str
    max_tokens: int = 256
    model: str = "scratch"


class CompletionsResponse(BaseModel):
    text: str
    usage: dict | None = None


def run_stub(prompt: str, max_tokens: int, model: str) -> tuple[str, dict]:
    """Stub: return fixed SQL. Replace with run_scratch/run_qlora when memorable is available."""
    return STUB_RESPONSE, {"prompt_tokens": 0, "completion_tokens": 1, "total_tokens": 1}


def run_memorable(prompt: str, max_tokens: int, model: str) -> tuple[str, dict]:
    """Call memorable inference if available."""
    try:
        from inference import run_scratch, run_qlora
        from pathlib import Path
        project_root = Path(MEMORABLE_PATH or ".").resolve()
        model_dir = project_root / "models" / ("postgres-llm-scratch" if model == "scratch" else "postgres-llm-qlora")
        if model == "scratch":
            text = run_scratch(prompt, model_dir, max_tokens)
        else:
            text = run_qlora(prompt, model_dir, "codellama/CodeLlama-7b-Instruct-hf", max_tokens)
        return text, {"prompt_tokens": 0, "completion_tokens": len(text.split()), "total_tokens": 0}
    except Exception as e:
        raise RuntimeError(f"Memorable inference failed: {e}") from e


@asynccontextmanager
async def lifespan(app: FastAPI):
    yield
    # shutdown


app = FastAPI(title="PGLang", version="0.1.0", lifespan=lifespan)


@app.post("/v1/completions", response_model=CompletionsResponse)
async def completions(req: CompletionsRequest):
    try:
        if MEMORABLE_PATH:
            text, usage = run_memorable(req.prompt, req.max_tokens, req.model)
        else:
            text, usage = run_stub(req.prompt, req.max_tokens, req.model)
        return CompletionsResponse(text=text, usage=usage)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/health")
async def health():
    return {"status": "ok"}


if __name__ == "__main__":
    import uvicorn
    port = int(os.environ.get("PORT", "9090"))
    uvicorn.run(app, host="0.0.0.0", port=port)

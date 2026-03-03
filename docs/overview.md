# 🤖 NeuronAgent

<div align="center">

**AI agent runtime system with REST API and WebSocket endpoints**

[![Status](https://img.shields.io/badge/status-stable-brightgreen)](.)
[![API](https://img.shields.io/badge/API-REST%20%7C%20WebSocket-blue)](.)
[![Tools](https://img.shields.io/badge/tools-16+-green)](.)

</div>

---

> [!TIP]
> NeuronAgent provides a complete platform for building autonomous AI agents. It includes persistent memory, tool execution, and multi-agent collaboration.

---

## 📋 What It Is

NeuronAgent is an AI agent runtime system providing REST API and WebSocket endpoints for building autonomous agent applications.

| Feature | Description | Status |
|---------|-------------|--------|
| **Agent Runtime** | Complete state machine for autonomous task execution with persistent memory | ✅ Stable |
| **REST API** | Full CRUD API for agents, sessions, messages, and advanced features | ✅ Stable |
| **WebSocket Support** | Real-time streaming agent responses | ✅ Stable |
| **Tool System** | Extensible tool registry with 16+ built-in tools (extensible via custom registration) | ✅ Stable |
| **Multi-Agent Collaboration** | Agent-to-agent communication and task delegation | ✅ Stable |
| **Workflow Engine** | DAG-based workflow execution with human-in-the-loop support | ✅ Stable |
| **Memory Management** | HNSW-based vector search for long-term memory with hierarchical organization | ✅ Stable |
| **Integration** | Direct integration with NeuronDB for embeddings, LLM, and vector operations | ✅ Stable |

## 🎯 Key Features & Modules

See [features.md](features.md) and [api.md](api.md) for full details. For architecture, see [architecture.md](architecture.md).

---

## 📚 Documentation

| Resource | Location | Description |
|----------|----------|-------------|
| **API Reference** | [api.md](api.md) | Complete API documentation |
| **Architecture** | [architecture.md](architecture.md) | Architecture details |
| **API Reference (detailed)** | [api-reference.md](api-reference.md) | Full API reference |
| **OpenAPI Spec** | `src/openapi/openapi.yaml` | OpenAPI 3.0 specification |

---

## 🐳 Docker

| Service | Description |
|---------|-------------|
| **neuronagent** | Main service (CPU) |
| **neuronagent-cuda** | NVIDIA GPU variant |
| **neuronagent-rocm** | AMD GPU variant |
| **neuronagent-metal** | Apple Silicon GPU variant |

---

## 🚀 Quick Start

```bash
# Check health endpoint
curl -sS http://localhost:8080/health

# Expected output: {"status":"ok"}
```

See [README.md](../README.md) for complete setup instructions.

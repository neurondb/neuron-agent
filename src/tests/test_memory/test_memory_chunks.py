"""Tests for Memory Chunks - embedding storage and retrieval."""
import pytest

@pytest.mark.requires_server
@pytest.mark.requires_neurondb
class TestMemoryChunks:
    def test_memory_chunks(self, api_client, test_agent):
        """Test memory chunk operations."""
        response = api_client.get(f"/api/v1/agents/{test_agent['id']}/memory")
        assert isinstance(response, list)


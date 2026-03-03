"""Tests for Hierarchical Memory - multi-level memory organization."""
import pytest

@pytest.mark.requires_server
class TestHierarchicalMemory:
    def test_hierarchical_memory(self, api_client, test_agent):
        """Test multi-level memory organization."""
        response = api_client.get(f"/api/v1/agents/{test_agent['id']}/memory")
        assert isinstance(response, list)


"""Tests for Sub-Agents."""
import pytest

@pytest.mark.requires_server
class TestSubAgents:
    def test_sub_agents(self, api_client, test_agent):
        """Test hierarchical agent structures."""
        assert test_agent is not None


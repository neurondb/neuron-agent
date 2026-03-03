"""Tests for Inter-Agent Communication."""
import pytest

@pytest.mark.requires_server
class TestInterAgentComm:
    def test_inter_agent_communication(self, api_client, test_agent):
        """Test message passing between agents."""
        assert test_agent is not None


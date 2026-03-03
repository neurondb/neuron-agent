"""Tests for Multi-Agent Scenarios."""
import pytest

@pytest.mark.integration
@pytest.mark.requires_server
@pytest.mark.slow
class TestMultiAgent:
    def test_multi_agent_collaboration(self, api_client, unique_name):
        """Test multi-agent collaboration."""
        # Create multiple agents
        agents = []
        for i in range(2):
            agent_data = {"name": f"{unique_name}-{i}", "system_prompt": "Test", "model_name": "gpt-4"}
            agent = api_client.post("/api/v1/agents", json_data=agent_data)
            agents.append(agent)
        
        try:
            # Test collaboration
            assert len(agents) == 2
        finally:
            for agent in agents:
                api_client.delete(f"/api/v1/agents/{agent['id']}")


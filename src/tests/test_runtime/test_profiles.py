"""Tests for Agent Profiles - profile-based configuration."""
import pytest

@pytest.mark.requires_server
class TestProfiles:
    def test_agent_profiles(self, api_client, unique_name):
        """Test agent profile configuration."""
        agent_data = {"name": unique_name, "profile": "research", "system_prompt": "Test", "model_name": "gpt-4"}
        try:
            agent = api_client.post("/api/v1/agents", json_data=agent_data)
            assert agent is not None
        finally:
            if "id" in agent:
                api_client.delete(f"/api/v1/agents/{agent['id']}")


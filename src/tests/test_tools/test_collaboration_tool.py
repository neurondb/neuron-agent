"""Tests for Collaboration Tool - agent-to-agent communication."""
import pytest

@pytest.mark.tool
@pytest.mark.requires_server
class TestCollaborationTool:
    def test_agent_collaboration(self, api_client, test_agent, unique_name):
        """Test agent-to-agent communication."""
        agent2 = api_client.post("/api/v1/agents", json_data={"name": unique_name, "system_prompt": "Helper", "model_name": "gpt-4"})
        try:
            message_data = {"content": f"Ask agent {agent2['id']} to help with a task", "role": "user"}
            response = api_client.post(f"/api/v1/sessions/{test_agent['id']}/messages", json_data=message_data)
            assert response is not None
        finally:
            api_client.delete(f"/api/v1/agents/{agent2['id']}")


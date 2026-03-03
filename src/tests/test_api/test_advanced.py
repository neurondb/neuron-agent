"""
Comprehensive tests for Advanced API endpoints.

Tests advanced operations like cloning, delegation, planning, etc.
"""

import pytest
import uuid
from typing import Dict, Any


@pytest.mark.api
@pytest.mark.requires_server
class TestAdvancedOperations:
    """Test advanced agent operations."""
    
    def test_clone_agent(self, api_client, test_agent):
        """Test cloning an agent."""
        response = api_client.post(f"/api/v1/agents/{test_agent['id']}/clone")
        assert "id" in response
        assert response["id"] != test_agent["id"]
        
        # Cleanup
        api_client.delete(f"/api/v1/agents/{response['id']}")
    
    def test_delegate_to_agent(self, api_client, test_agent):
        """Test delegating task to another agent."""
        # Create another agent for delegation
        agent2_data = {
            "name": f"delegate-target-{uuid.uuid4().hex[:8]}",
            "system_prompt": "Specialized agent",
            "model_name": "gpt-4"
        }
        agent2 = api_client.post("/api/v1/agents", json_data=agent2_data)
        
        try:
            delegate_data = {
                "target_agent_id": agent2["id"],
                "task": "Perform specialized task"
            }
            response = api_client.post(
                f"/api/v1/agents/{test_agent['id']}/delegate",
                json_data=delegate_data
            )
            assert isinstance(response, dict)
        finally:
            api_client.delete(f"/api/v1/agents/{agent2['id']}")
    
    def test_generate_plan(self, api_client, test_agent):
        """Test generating a plan for an agent."""
        plan_data = {
            "goal": "Complete a complex task",
            "context": "Test context"
        }
        
        response = api_client.post(
            f"/api/v1/agents/{test_agent['id']}/plan",
            json_data=plan_data
        )
        assert isinstance(response, dict)
        # Should have plan structure
    
    def test_reflect_on_response(self, api_client, test_session):
        """Test reflecting on agent response."""
        # First send a message
        message_data = {
            "content": "Test message for reflection",
            "role": "user"
        }
        msg_result = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        
        # Get message ID if available
        message_id = msg_result.get("id") or (msg_result.get("message", {}).get("id") if isinstance(msg_result, dict) else None)
        if message_id:
            reflection_data = {
                "message_id": message_id,
                "criteria": ["accuracy", "relevance"]
            }
            response = api_client.post(
                f"/api/v1/sessions/{test_session['id']}/reflect",
                json_data=reflection_data
            )
            assert isinstance(response, dict)


@pytest.mark.api
@pytest.mark.requires_server
class TestAnalytics:
    """Test analytics endpoints."""
    
    def test_get_analytics_overview(self, api_client):
        """Test getting analytics overview."""
        response = api_client.get("/api/v1/analytics/overview")
        assert isinstance(response, dict)
        # Should have analytics data


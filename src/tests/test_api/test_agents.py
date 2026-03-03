"""
Comprehensive tests for Agent API endpoints.

Tests all agent CRUD operations, profiles, router, and advanced features.
"""

import pytest
import uuid
from typing import Dict, Any


@pytest.mark.api
@pytest.mark.requires_server
class TestAgentCRUD:
    """Test agent CRUD operations."""
    
    def test_create_agent(self, api_client, unique_name):
        """Test creating a new agent."""
        agent_data = {
            "name": unique_name,
            "description": "Test agent",
            "system_prompt": "You are a helpful assistant.",
            "model_name": "gpt-4",
            "enabled_tools": ["sql", "http"],
            "config": {
                "temperature": 0.7,
                "max_tokens": 1000
            }
        }
        
        response = api_client.post("/api/v1/agents", json_data=agent_data)
        
        assert "id" in response
        assert response["name"] == unique_name
        assert response["system_prompt"] == agent_data["system_prompt"]
        assert response["model_name"] == agent_data["model_name"]
        assert "sql" in response.get("enabled_tools", [])
        assert_valid_uuid(response["id"])
        
        # Cleanup
        api_client.delete(f"/api/v1/agents/{response['id']}")
    
    def test_get_agent(self, api_client, test_agent):
        """Test retrieving an agent by ID."""
        response = api_client.get(f"/api/v1/agents/{test_agent['id']}")
        
        assert response["id"] == test_agent["id"]
        assert response["name"] == test_agent["name"]
        assert "system_prompt" in response
        assert "model_name" in response
    
    def test_list_agents(self, api_client):
        """Test listing all agents."""
        response = api_client.get("/api/v1/agents")
        
        assert isinstance(response, list)
        # Should have at least our test agents
    
    def test_list_agents_with_search(self, api_client, unique_name):
        """Test listing agents with search filter."""
        # Create an agent with unique name
        agent_data = {
            "name": f"{unique_name}-search",
            "system_prompt": "Test",
            "model_name": "gpt-4"
        }
        agent = api_client.post("/api/v1/agents", json_data=agent_data)
        
        try:
            response = api_client.get("/api/v1/agents", params={"search": unique_name})
            assert isinstance(response, list)
            # Should find our agent
            found = any(a["id"] == agent["id"] for a in response)
            assert found, "Created agent should be in search results"
        finally:
            api_client.delete(f"/api/v1/agents/{agent['id']}")
    
    def test_update_agent(self, api_client, test_agent):
        """Test updating an agent."""
        update_data = {
            "name": test_agent["name"],
            "system_prompt": "Updated system prompt",
            "model_name": "gpt-4",
            "description": "Updated description",
            "enabled_tools": ["sql", "http", "code"]
        }
        
        response = api_client.put(f"/api/v1/agents/{test_agent['id']}", json_data=update_data)
        
        assert response["system_prompt"] == update_data["system_prompt"]
        assert response["description"] == update_data["description"]
        assert len(response.get("enabled_tools", [])) == 3
    
    def test_delete_agent(self, api_client, unique_name):
        """Test deleting an agent."""
        # Create agent
        agent_data = {
            "name": unique_name,
            "system_prompt": "Test",
            "model_name": "gpt-4"
        }
        agent = api_client.post("/api/v1/agents", json_data=agent_data)
        agent_id = agent["id"]
        
        # Delete agent
        api_client.delete(f"/api/v1/agents/{agent_id}")
        
        # Verify deletion
        with pytest.raises(Exception):  # Should raise 404
            api_client.get(f"/api/v1/agents/{agent_id}")


@pytest.mark.api
@pytest.mark.requires_server
class TestAgentValidation:
    """Test agent validation and error handling."""
    
    def test_create_agent_missing_required_fields(self, api_client):
        """Test creating agent with missing required fields."""
        # Missing system_prompt
        agent_data = {
            "name": f"test-{uuid.uuid4().hex[:8]}",
            "model_name": "gpt-4"
        }
        
        with pytest.raises(Exception):  # Should raise 400
            api_client.post("/api/v1/agents", json_data=agent_data)
    
    def test_create_agent_invalid_uuid(self, api_client):
        """Test operations with invalid UUID."""
        with pytest.raises(Exception):  # Should raise 400
            api_client.get("/api/v1/agents/invalid-uuid")
    
    def test_get_nonexistent_agent(self, api_client):
        """Test retrieving non-existent agent."""
        fake_id = str(uuid.uuid4())
        with pytest.raises(Exception):  # Should raise 404
            api_client.get(f"/api/v1/agents/{fake_id}")
    
    def test_create_agent_duplicate_name(self, api_client, test_agent):
        """Test creating agent with duplicate name."""
        agent_data = {
            "name": test_agent["name"],
            "system_prompt": "Test",
            "model_name": "gpt-4"
        }
        
        with pytest.raises(Exception):  # Should raise 400 or 409
            api_client.post("/api/v1/agents", json_data=agent_data)


@pytest.mark.api
@pytest.mark.requires_server
class TestAgentAdvanced:
    """Test advanced agent operations."""
    
    def test_clone_agent(self, api_client, test_agent):
        """Test cloning an agent."""
        response = api_client.post(f"/api/v1/agents/{test_agent['id']}/clone")
        
        assert "id" in response
        assert response["id"] != test_agent["id"]
        assert "_clone" in response["name"] or response["name"] != test_agent["name"]
        assert response["system_prompt"] == test_agent["system_prompt"]
        
        # Cleanup
        api_client.delete(f"/api/v1/agents/{response['id']}")
    
    def test_get_agent_metrics(self, api_client, test_agent):
        """Test getting agent metrics."""
        response = api_client.get(f"/api/v1/agents/{test_agent['id']}/metrics")
        
        assert isinstance(response, dict)
        # Metrics should have relevant fields
    
    def test_get_agent_costs(self, api_client, test_agent):
        """Test getting agent cost tracking."""
        response = api_client.get(f"/api/v1/agents/{test_agent['id']}/costs")
        
        assert isinstance(response, dict)
        # Should have cost-related fields


@pytest.mark.api
@pytest.mark.requires_server
class TestAgentProfiles:
    """Test agent profiles functionality."""
    
    def test_create_agent_with_profile(self, api_client, unique_name):
        """Test creating agent with profile."""
        agent_data = {
            "name": unique_name,
            "profile": "research",
            "system_prompt": "Test",
            "model_name": "gpt-4"
        }
        
        response = None
        try:
            response = api_client.post("/api/v1/agents", json_data=agent_data)
            assert "id" in response
            # Profile should be applied
        finally:
            if response and "id" in response:
                api_client.delete(f"/api/v1/agents/{response['id']}")


# Helper function from conftest
def assert_valid_uuid(value: str):
    """Assert that value is a valid UUID."""
    try:
        uuid.UUID(value)
    except ValueError:
        pytest.fail(f"Invalid UUID: {value}")


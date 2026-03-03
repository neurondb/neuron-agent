"""
Comprehensive tests for Tool Management API endpoints.

Tests tool registration, management, and tool-related operations.
"""

import pytest
import uuid
from typing import Dict, Any


@pytest.mark.api
@pytest.mark.requires_server
class TestToolCRUD:
    """Test tool CRUD operations."""
    
    def test_create_tool(self, api_client, unique_name):
        """Test creating a custom tool."""
        tool_data = {
            "name": unique_name,
            "description": "Test tool",
            "handler_type": "http",
            "config": {
                "url": "https://api.example.com",
                "method": "GET"
            }
        }
        
        try:
            response = api_client.post("/api/v1/tools", json_data=tool_data)
            assert "id" in response
            assert response["name"] == unique_name
            assert response["handler_type"] == "http"
        finally:
            if "id" in response:
                api_client.delete(f"/api/v1/tools/{response['id']}")
    
    def test_list_tools(self, api_client):
        """Test listing all tools."""
        response = api_client.get("/api/v1/tools")
        assert isinstance(response, list)
    
    def test_get_tool(self, api_client, unique_name):
        """Test retrieving a tool by ID."""
        # Create tool first
        tool_data = {
            "name": unique_name,
            "description": "Test tool",
            "handler_type": "http"
        }
        tool = api_client.post("/api/v1/tools", json_data=tool_data)
        
        try:
            response = api_client.get(f"/api/v1/tools/{tool['id']}")
            assert response["id"] == tool["id"]
            assert response["name"] == unique_name
        finally:
            api_client.delete(f"/api/v1/tools/{tool['id']}")
    
    def test_update_tool(self, api_client, unique_name):
        """Test updating a tool."""
        tool_data = {
            "name": unique_name,
            "description": "Test tool",
            "handler_type": "http"
        }
        tool = api_client.post("/api/v1/tools", json_data=tool_data)
        
        try:
            update_data = {
                "name": unique_name,
                "description": "Updated description",
                "handler_type": "http"
            }
            response = api_client.put(f"/api/v1/tools/{tool['id']}", json_data=update_data)
            assert response["description"] == "Updated description"
        finally:
            api_client.delete(f"/api/v1/tools/{tool['id']}")
    
    def test_delete_tool(self, api_client, unique_name):
        """Test deleting a tool."""
        tool_data = {
            "name": unique_name,
            "description": "Test tool",
            "handler_type": "http"
        }
        tool = api_client.post("/api/v1/tools", json_data=tool_data)
        tool_id = tool["id"]
        
        api_client.delete(f"/api/v1/tools/{tool_id}")
        
        # Verify deletion
        with pytest.raises(Exception):
            api_client.get(f"/api/v1/tools/{tool_id}")


@pytest.mark.api
@pytest.mark.requires_server
class TestToolAnalytics:
    """Test tool analytics endpoints."""
    
    def test_get_tool_analytics(self, api_client):
        """Test getting tool analytics."""
        tools = api_client.get("/api/v1/tools")
        if tools and len(tools) > 0:
            tool_id = tools[0]["id"]
            response = api_client.get(f"/api/v1/tools/{tool_id}/analytics")
            assert isinstance(response, dict)


"""
Comprehensive tests for Session API endpoints.

Tests session management, caching, and lifecycle operations.
"""

import pytest
import uuid
from typing import Dict, Any


@pytest.mark.api
@pytest.mark.requires_server
class TestSessionCRUD:
    """Test session CRUD operations."""
    
    def test_create_session(self, api_client, test_agent, unique_name):
        """Test creating a new session."""
        session_data = {
            "agent_id": test_agent["id"],
            "external_user_id": unique_name,
            "metadata": {"test": True, "source": "pytest"}
        }
        
        response = api_client.post("/api/v1/sessions", json_data=session_data)
        
        assert "id" in response
        assert response["agent_id"] == test_agent["id"]
        assert response["external_user_id"] == unique_name
        assert_valid_uuid(response["id"])
        
        # Cleanup
        api_client.delete(f"/api/v1/sessions/{response['id']}")
    
    def test_get_session(self, api_client, test_session):
        """Test retrieving a session by ID."""
        response = api_client.get(f"/api/v1/sessions/{test_session['id']}")
        
        assert response["id"] == test_session["id"]
        assert response["agent_id"] == test_session["agent_id"]
        assert "created_at" in response
    
    def test_list_sessions(self, api_client, test_agent):
        """Test listing sessions."""
        response = api_client.get(f"/api/v1/agents/{test_agent['id']}/sessions")
        
        assert isinstance(response, list)
    
    def test_list_sessions_all(self, api_client):
        """Test listing all sessions."""
        response = api_client.get("/api/v1/sessions")
        
        assert isinstance(response, list)
    
    def test_update_session(self, api_client, test_session):
        """Test updating a session."""
        update_data = {
            "metadata": {"updated": True, "test": True}
        }
        
        response = api_client.put(f"/api/v1/sessions/{test_session['id']}", json_data=update_data)
        
        assert response["metadata"].get("updated") is True
    
    def test_delete_session(self, api_client, test_agent, unique_name):
        """Test deleting a session."""
        # Create session
        session_data = {
            "agent_id": test_agent["id"],
            "external_user_id": unique_name
        }
        session = api_client.post("/api/v1/sessions", json_data=session_data)
        session_id = session["id"]
        
        # Delete session
        api_client.delete(f"/api/v1/sessions/{session_id}")
        
        # Verify deletion
        with pytest.raises(Exception):  # Should raise 404
            api_client.get(f"/api/v1/sessions/{session_id}")


@pytest.mark.api
@pytest.mark.requires_server
class TestSessionValidation:
    """Test session validation and error handling."""
    
    def test_create_session_invalid_agent_id(self, api_client):
        """Test creating session with invalid agent ID."""
        fake_id = str(uuid.uuid4())
        session_data = {
            "agent_id": fake_id,
            "external_user_id": "test-user"
        }
        
        with pytest.raises(Exception):  # Should raise 404 or 400
            api_client.post("/api/v1/sessions", json_data=session_data)
    
    def test_get_nonexistent_session(self, api_client):
        """Test retrieving non-existent session."""
        fake_id = str(uuid.uuid4())
        with pytest.raises(Exception):  # Should raise 404
            api_client.get(f"/api/v1/sessions/{fake_id}")


@pytest.mark.api
@pytest.mark.requires_server
class TestSessionMetadata:
    """Test session metadata operations."""
    
    def test_session_with_metadata(self, api_client, test_agent, unique_name):
        """Test creating session with metadata."""
        metadata = {
            "user_id": "user123",
            "session_type": "test",
            "custom_field": "custom_value"
        }
        
        session_data = {
            "agent_id": test_agent["id"],
            "external_user_id": unique_name,
            "metadata": metadata
        }
        
        try:
            session = api_client.post("/api/v1/sessions", json_data=session_data)
            assert session["metadata"] == metadata
        finally:
            api_client.delete(f"/api/v1/sessions/{session['id']}")


# Helper function
def assert_valid_uuid(value: str):
    """Assert that value is a valid UUID."""
    try:
        uuid.UUID(value)
    except ValueError:
        pytest.fail(f"Invalid UUID: {value}")


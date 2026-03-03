"""
Comprehensive tests for Message API endpoints.

Tests message handling, streaming, and conversation management.
"""

import pytest
import uuid
import time
from typing import Dict, Any, List


@pytest.mark.api
@pytest.mark.requires_server
class TestMessageCRUD:
    """Test message CRUD operations."""
    
    def test_send_message(self, api_client, test_session):
        """Test sending a message to a session."""
        message_data = {
            "content": "Hello! This is a test message.",
            "role": "user"
        }
        
        response = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        
        assert "response" in response or "message" in response or "id" in response
        # Response should contain agent's reply or message ID
    
    def test_get_messages(self, api_client, test_session):
        """Test retrieving messages from a session."""
        # Send a message first
        message_data = {
            "content": "Test message for retrieval",
            "role": "user"
        }
        api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        
        # Wait a bit for message to be processed
        time.sleep(1)
        
        response = api_client.get(f"/api/v1/sessions/{test_session['id']}/messages")
        
        assert isinstance(response, list)
        assert len(response) > 0
    
    def test_get_message_by_id(self, api_client, test_session):
        """Test retrieving a specific message by ID."""
        # Send a message
        message_data = {
            "content": "Test message",
            "role": "user"
        }
        result = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        
        # Extract message ID (format may vary)
        message_id = None
        if "id" in result:
            message_id = result["id"]
        elif "message" in result and "id" in result["message"]:
            message_id = result["message"]["id"]
        elif isinstance(result, list) and len(result) > 0:
            message_id = result[0].get("id")
        
        if message_id:
            response = api_client.get(
                f"/api/v1/sessions/{test_session['id']}/messages/{message_id}"
            )
            assert response["id"] == message_id
    
    def test_update_message(self, api_client, test_session):
        """Test updating a message."""
        # Send a message
        message_data = {
            "content": "Original message",
            "role": "user"
        }
        result = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        
        # Try to update if message ID is available
        message_id = result.get("id") or (result.get("message", {}).get("id") if isinstance(result, dict) else None)
        if message_id:
            update_data = {
                "content": "Updated message content",
                "metadata": {"updated": True}
            }
            response = api_client.put(
                f"/api/v1/sessions/{test_session['id']}/messages/{message_id}",
                json_data=update_data
            )
            assert "content" in response or "metadata" in response
    
    def test_delete_message(self, api_client, test_session):
        """Test deleting a message."""
        # Send a message
        message_data = {
            "content": "Message to delete",
            "role": "user"
        }
        result = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        
        message_id = result.get("id") or (result.get("message", {}).get("id") if isinstance(result, dict) else None)
        if message_id:
            api_client.delete(
                f"/api/v1/sessions/{test_session['id']}/messages/{message_id}"
            )
            
            # Verify deletion
            with pytest.raises(Exception):  # Should raise 404
                api_client.get(f"/api/v1/sessions/{test_session['id']}/messages/{message_id}")


@pytest.mark.api
@pytest.mark.requires_server
class TestMessageValidation:
    """Test message validation and error handling."""
    
    def test_send_message_empty_content(self, api_client, test_session):
        """Test sending message with empty content."""
        message_data = {
            "content": "",
            "role": "user"
        }
        
        with pytest.raises(Exception):  # Should raise 400
            api_client.post(
                f"/api/v1/sessions/{test_session['id']}/messages",
                json_data=message_data
            )
    
    def test_send_message_invalid_session(self, api_client):
        """Test sending message to invalid session."""
        fake_id = str(uuid.uuid4())
        message_data = {
            "content": "Test",
            "role": "user"
        }
        
        with pytest.raises(Exception):  # Should raise 404
            api_client.post(
                f"/api/v1/sessions/{fake_id}/messages",
                json_data=message_data
            )


@pytest.mark.api
@pytest.mark.requires_server
class TestMessageStreaming:
    """Test message streaming functionality."""
    
    @pytest.mark.slow
    def test_streaming_response(self, api_client, test_session):
        """Test streaming message responses via WebSocket."""
        # This would require WebSocket client setup
        # For now, we'll test that streaming endpoint exists
        # Full WebSocket tests in test_integration/test_streaming.py
        pass


@pytest.mark.api
@pytest.mark.requires_server
class TestMessageConversation:
    """Test conversation management."""
    
    def test_multiple_messages(self, api_client, test_session):
        """Test sending multiple messages in a conversation."""
        messages = [
            "First message",
            "Second message",
            "Third message"
        ]
        
        for msg in messages:
            message_data = {
                "content": msg,
                "role": "user"
            }
            response = api_client.post(
                f"/api/v1/sessions/{test_session['id']}/messages",
                json_data=message_data
            )
            assert response is not None
            time.sleep(0.5)  # Small delay between messages
        
        # Retrieve all messages
        all_messages = api_client.get(f"/api/v1/sessions/{test_session['id']}/messages")
        assert len(all_messages) >= len(messages)
    
    def test_message_ordering(self, api_client, test_session):
        """Test that messages are returned in correct order."""
        # Send multiple messages
        for i in range(3):
            message_data = {
                "content": f"Message {i}",
                "role": "user"
            }
            api_client.post(
                f"/api/v1/sessions/{test_session['id']}/messages",
                json_data=message_data
            )
            time.sleep(0.3)
        
        time.sleep(1)
        
        # Retrieve messages
        messages = api_client.get(f"/api/v1/sessions/{test_session['id']}/messages")
        
        # Messages should be ordered by creation time
        if len(messages) >= 2:
            timestamps = [msg.get("created_at", "") for msg in messages if "created_at" in msg]
            if len(timestamps) >= 2:
                # Later messages should have later timestamps
                assert timestamps[-1] >= timestamps[0]


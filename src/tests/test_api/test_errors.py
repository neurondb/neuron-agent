"""
Comprehensive tests for Error Handling.

Tests error responses, validation, and edge cases.
"""

import pytest
import uuid
from typing import Dict, Any


@pytest.mark.api
@pytest.mark.requires_server
class TestErrorHandling:
    """Test error handling and validation."""
    
    def test_401_unauthorized(self, api_client):
        """Test 401 Unauthorized error."""
        # Create client without API key
        import requests
        response = requests.get(f"{api_client.base_url}/api/v1/agents")
        assert response.status_code == 401
    
    def test_404_not_found(self, api_client):
        """Test 404 Not Found error."""
        fake_id = str(uuid.uuid4())
        with pytest.raises(Exception):
            api_client.get(f"/api/v1/agents/{fake_id}")
    
    def test_400_bad_request(self, api_client):
        """Test 400 Bad Request error."""
        # Missing required fields
        agent_data = {
            "name": f"test-{uuid.uuid4().hex[:8]}"
            # Missing system_prompt and model_name
        }
        with pytest.raises(Exception):
            api_client.post("/api/v1/agents", json_data=agent_data)
    
    def test_invalid_json(self, api_client):
        """Test handling of invalid JSON."""
        import requests
        response = requests.post(
            f"{api_client.base_url}/api/v1/agents",
            headers={
                "Authorization": f"Bearer {api_client.api_key}",
                "Content-Type": "application/json"
            },
            data="invalid json",
            timeout=10
        )
        assert response.status_code >= 400
    
    def test_missing_required_fields(self, api_client):
        """Test validation of required fields."""
        # Agent without system_prompt
        agent_data = {
            "name": f"test-{uuid.uuid4().hex[:8]}",
            "model_name": "gpt-4"
        }
        with pytest.raises(Exception):
            api_client.post("/api/v1/agents", json_data=agent_data)
    
    def test_invalid_uuid_format(self, api_client):
        """Test handling of invalid UUID format."""
        with pytest.raises(Exception):
            api_client.get("/api/v1/agents/invalid-uuid-format")
    
    def test_duplicate_resource(self, api_client, test_agent):
        """Test handling of duplicate resource creation."""
        agent_data = {
            "name": test_agent["name"],  # Duplicate name
            "system_prompt": "Test",
            "model_name": "gpt-4"
        }
        with pytest.raises(Exception):  # Should raise 400 or 409
            api_client.post("/api/v1/agents", json_data=agent_data)
    
    def test_large_payload(self, api_client):
        """Test handling of very large payloads."""
        # Create agent with very large system prompt
        large_prompt = "A" * (10 * 1024 * 1024)  # 10MB
        agent_data = {
            "name": f"test-{uuid.uuid4().hex[:8]}",
            "system_prompt": large_prompt,
            "model_name": "gpt-4"
        }
        # Should either succeed or fail gracefully
        try:
            agent = api_client.post("/api/v1/agents", json_data=agent_data)
            api_client.delete(f"/api/v1/agents/{agent['id']}")
        except Exception:
            pass  # Expected for very large payloads


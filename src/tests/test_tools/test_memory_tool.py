"""Tests for Memory Tool - direct memory manipulation."""
import pytest

@pytest.mark.tool
@pytest.mark.requires_server
class TestMemoryTool:
    def test_memory_manipulation(self, api_client, test_session):
        """Test direct memory manipulation."""
        message_data = {"content": "Store this in memory: Important information", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


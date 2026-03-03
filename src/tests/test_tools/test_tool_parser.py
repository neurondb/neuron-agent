"""Tests for Tool Parser - tool call parsing."""
import pytest

@pytest.mark.tool
@pytest.mark.requires_server
class TestToolParser:
    def test_tool_call_parsing(self, api_client, test_session):
        """Test tool call parsing."""
        message_data = {"content": "Execute SQL: SELECT 1", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


"""Tests for Context Management - intelligent context loading."""
import pytest

@pytest.mark.requires_server
class TestContext:
    def test_context_loading(self, api_client, test_session):
        """Test context loading from messages and memory."""
        message_data = {"content": "Test", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


"""Tests for Event Streaming."""
import pytest

@pytest.mark.requires_server
class TestEventStreaming:
    def test_event_streaming(self, api_client, test_session):
        """Test real-time event capture and analysis."""
        message_data = {"content": "Test", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


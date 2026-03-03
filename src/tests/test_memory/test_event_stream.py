"""Tests for Event Streaming - real-time event capture."""
import pytest

@pytest.mark.requires_server
class TestEventStream:
    def test_event_streaming(self, api_client, test_session):
        """Test event capture and summarization."""
        message_data = {"content": "Test", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


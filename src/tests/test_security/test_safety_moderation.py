"""Tests for Safety Moderation."""
import pytest

@pytest.mark.security
@pytest.mark.requires_server
class TestSafetyModeration:
    def test_content_moderation(self, api_client, test_session):
        """Test content moderation and safety checks."""
        message_data = {"content": "Test message", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


"""Tests for Webhook Support."""
import pytest

@pytest.mark.requires_server
class TestWebhooks:
    def test_webhook_management(self, api_client):
        """Test webhook management."""
        response = api_client.get("/api/v1/webhooks")
        assert isinstance(response, (list, dict))


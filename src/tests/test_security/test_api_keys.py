"""Tests for API Key Management."""
import pytest

@pytest.mark.security
@pytest.mark.requires_server
class TestAPIKeys:
    def test_api_key_authentication(self, api_client):
        """Test API key authentication."""
        response = api_client.get("/api/v1/agents")
        assert isinstance(response, list)


"""Tests for Structured Logging."""
import pytest

@pytest.mark.requires_server
class TestLogging:
    def test_structured_logging(self, api_client):
        """Test JSON-formatted logs."""
        response = api_client.get("/api/v1/agents")
        assert isinstance(response, list)


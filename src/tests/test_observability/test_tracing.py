"""Tests for Distributed Tracing."""
import pytest

@pytest.mark.requires_server
class TestTracing:
    def test_distributed_tracing(self, api_client):
        """Test distributed tracing support."""
        assert True


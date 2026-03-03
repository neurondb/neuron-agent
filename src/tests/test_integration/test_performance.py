"""Tests for Performance Benchmarks."""
import pytest
import time

@pytest.mark.performance
@pytest.mark.requires_server
@pytest.mark.slow
class TestPerformance:
    def test_api_response_time(self, api_client):
        """Test API response time."""
        start = time.time()
        response = api_client.get("/api/v1/agents")
        elapsed = time.time() - start
        assert elapsed < 5.0, "API should respond within 5 seconds"
        assert isinstance(response, list)


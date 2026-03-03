"""Tests for Async Tasks."""
import pytest

@pytest.mark.requires_server
@pytest.mark.slow
class TestAsyncTasks:
    def test_async_tasks(self, api_client):
        """Test background task execution."""
        response = api_client.get("/api/v1/async-tasks")
        assert isinstance(response, (list, dict))


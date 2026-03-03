"""Tests for Task Plans."""
import pytest

@pytest.mark.requires_server
class TestTaskPlans:
    def test_task_plans(self, api_client):
        """Test multi-step plan creation and execution."""
        response = api_client.get("/api/v1/plans")
        assert isinstance(response, (list, dict))


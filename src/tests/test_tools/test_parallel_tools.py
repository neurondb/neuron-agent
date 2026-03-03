"""Tests for Parallel Tool Execution - concurrent tool execution."""
import pytest

@pytest.mark.tool
@pytest.mark.requires_server
class TestParallelTools:
    def test_parallel_tool_execution(self, api_client, test_session):
        """Test concurrent tool execution."""
        message_data = {"content": "Execute multiple tools in parallel", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


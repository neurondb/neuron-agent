"""Tests for Workspace Management."""
import pytest

@pytest.mark.requires_server
class TestWorkspace:
    def test_workspace_management(self, api_client):
        """Test shared workspaces for collaborative agents."""
        response = api_client.get("/api/v1/workspaces")
        assert isinstance(response, (list, dict))


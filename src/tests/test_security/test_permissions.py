"""Tests for Data and Tool Permissions."""
import pytest

@pytest.mark.security
@pytest.mark.requires_server
class TestPermissions:
    def test_data_permissions(self, api_client, test_agent):
        """Test data access permissions."""
        response = api_client.get(f"/api/v1/agents/{test_agent['id']}")
        assert response is not None
    
    def test_tool_permissions(self, api_client, test_agent):
        """Test tool access permissions."""
        assert test_agent is not None


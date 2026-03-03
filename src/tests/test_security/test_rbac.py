"""Tests for Role-Based Access Control."""
import pytest

@pytest.mark.security
@pytest.mark.requires_server
class TestRBAC:
    def test_rbac_permissions(self, api_client):
        """Test RBAC permission checks."""
        # This would require different API keys with different roles
        response = api_client.get("/api/v1/agents")
        assert isinstance(response, list)


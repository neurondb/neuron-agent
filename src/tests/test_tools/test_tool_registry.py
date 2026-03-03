"""Tests for Tool Registry - extensible tool registration."""
import pytest

@pytest.mark.tool
@pytest.mark.requires_server
class TestToolRegistry:
    def test_tool_registry(self, api_client):
        """Test tool registry operations."""
        response = api_client.get("/api/v1/tools")
        assert isinstance(response, list)


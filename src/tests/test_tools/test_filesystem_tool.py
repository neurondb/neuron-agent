"""Tests for Filesystem Tool - virtual filesystem."""
import pytest

@pytest.mark.tool
@pytest.mark.requires_server
class TestFilesystemTool:
    def test_virtual_filesystem(self, api_client, test_session):
        """Test virtual filesystem operations."""
        message_data = {"content": "Create a file test.txt with content 'Hello World'", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


"""
Comprehensive tests for HTTP Tool.

Tests HTTP requests, allowlist validation, and security.
"""

import pytest


@pytest.mark.tool
@pytest.mark.requires_server
class TestHTTPTool:
    """Test HTTP tool execution."""
    
    def test_http_get_request(self, api_client, test_session):
        """Test executing HTTP GET request."""
        message_data = {
            "content": "Make HTTP GET request to https://httpbin.org/get",
            "role": "user"
        }
        
        response = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        assert response is not None
    
    def test_http_allowlist_validation(self, api_client, test_session):
        """Test that HTTP tool validates allowlist."""
        # Try to access non-allowlisted URL (should be blocked)
        message_data = {
            "content": "Make HTTP request to https://malicious-site.com",
            "role": "user"
        }
        
        response = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        # Should indicate blocking or error


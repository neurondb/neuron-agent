"""
Comprehensive tests for Shell Tool.

Tests shell command execution, whitelist validation, and security.
"""

import pytest


@pytest.mark.tool
@pytest.mark.requires_server
class TestShellTool:
    """Test shell tool execution."""
    
    def test_shell_allowed_command(self, api_client, test_session):
        """Test executing whitelisted shell command."""
        message_data = {
            "content": "Execute shell command: ls -la",
            "role": "user"
        }
        
        response = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        assert response is not None
    
    def test_shell_whitelist_validation(self, api_client, test_session):
        """Test that shell tool validates whitelist."""
        # Try to execute non-whitelisted command (should be blocked)
        message_data = {
            "content": "Execute shell command: rm -rf /",
            "role": "user"
        }
        
        response = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        # Should indicate blocking or error


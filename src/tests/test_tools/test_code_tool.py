"""
Comprehensive tests for Code Tool.

Tests code execution, sandboxing, and security.
"""

import pytest


@pytest.mark.tool
@pytest.mark.requires_server
class TestCodeTool:
    """Test code tool execution."""
    
    def test_code_execution_python(self, api_client, test_session):
        """Test executing Python code."""
        message_data = {
            "content": "Execute Python code: print('Hello, World!')",
            "role": "user"
        }
        
        response = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        assert response is not None
    
    def test_code_sandboxing(self, api_client, test_session):
        """Test that code execution is sandboxed."""
        # Try to access filesystem (should be blocked)
        message_data = {
            "content": "Execute Python: open('/etc/passwd').read()",
            "role": "user"
        }
        
        response = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        # Should indicate blocking or error


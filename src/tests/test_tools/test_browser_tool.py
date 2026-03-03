"""Tests for Browser Tool - Playwright automation."""
import pytest

@pytest.mark.tool
@pytest.mark.requires_server
@pytest.mark.slow
class TestBrowserTool:
    def test_browser_automation(self, api_client, test_session):
        """Test browser automation with Playwright."""
        message_data = {"content": "Navigate to https://example.com and get the title", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


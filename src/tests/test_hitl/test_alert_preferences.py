"""Tests for Alert Preferences."""
import pytest

@pytest.mark.requires_server
class TestAlertPreferences:
    def test_alert_preferences(self, api_client):
        """Test configurable alert preferences."""
        response = api_client.get("/api/v1/alert-preferences")
        assert isinstance(response, (list, dict))


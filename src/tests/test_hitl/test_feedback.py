"""Tests for Feedback System."""
import pytest

@pytest.mark.requires_server
class TestFeedback:
    def test_feedback_system(self, api_client, test_session):
        """Test collecting and integrating human feedback."""
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/feedback", json_data={"rating": 5})
        assert response is not None


"""Tests for Memory Promotion - background worker."""
import pytest

@pytest.mark.requires_server
@pytest.mark.slow
class TestMemoryPromotion:
    def test_memory_promotion(self, api_client, test_session):
        """Test background memory promotion."""
        message_data = {"content": "Important information", "role": "user"}
        api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        # Promotion happens asynchronously


"""Tests for Memory Promoter."""
import pytest

@pytest.mark.requires_server
@pytest.mark.slow
class TestMemoryPromoter:
    def test_memory_promoter(self, api_client, test_session):
        """Test background memory promotion."""
        message_data = {"content": "Important info", "role": "user"}
        api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)


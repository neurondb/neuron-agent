"""Tests for Multimodal Tool - image and multimedia processing."""
import pytest

@pytest.mark.tool
@pytest.mark.requires_server
class TestMultimodalTool:
    def test_multimodal_processing(self, api_client, test_session):
        """Test image and multimedia processing."""
        message_data = {"content": "Process this image: [image data]", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


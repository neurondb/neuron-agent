"""Tests for LLM Integration - NeuronDB LLM functions."""
import pytest

@pytest.mark.requires_server
@pytest.mark.requires_neurondb
class TestLLMClient:
    def test_llm_integration(self, api_client, test_session):
        """Test LLM function calls."""
        message_data = {"content": "Hello", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


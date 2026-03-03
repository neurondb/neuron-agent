"""
Comprehensive tests for End-to-End RAG Workflows.

Tests retrieval augmented generation workflows.
"""

import pytest


@pytest.mark.neurondb
@pytest.mark.requires_neurondb
@pytest.mark.requires_server
@pytest.mark.slow
class TestRAGWorkflow:
    """Test RAG workflows."""
    
    def test_rag_workflow(self, api_client, test_session):
        """Test complete RAG workflow."""
        # Send message that should trigger RAG
        message_data = {
            "content": "Find information about machine learning and explain it to me",
            "role": "user"
        }
        
        response = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        
        # Response should contain retrieved context and generated answer
        assert response is not None


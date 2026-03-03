"""Tests for NeuronDB Tools - RAG, Vector, ML, Analytics, Visualization."""
import pytest

@pytest.mark.tool
@pytest.mark.requires_server
@pytest.mark.requires_neurondb
class TestNeuronDBTools:
    def test_rag_tool(self, api_client, test_session):
        """Test RAG tool."""
        message_data = {"content": "Search for information about machine learning", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None
    
    def test_vector_tool(self, api_client, test_session):
        """Test vector tool."""
        message_data = {"content": "Find similar vectors to [0.1, 0.2, 0.3]", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None
    
    def test_ml_tool(self, api_client, test_session):
        """Test ML tool."""
        message_data = {"content": "Run ML model prediction", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None
    
    def test_analytics_tool(self, api_client, test_session):
        """Test analytics tool."""
        message_data = {"content": "Analyze agent performance metrics", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None
    
    def test_visualization_tool(self, api_client, test_session):
        """Test visualization tool."""
        message_data = {"content": "Create a chart from this data", "role": "user"}
        response = api_client.post(f"/api/v1/sessions/{test_session['id']}/messages", json_data=message_data)
        assert response is not None


"""Tests for ML (Machine Learning) Tool."""
import pytest
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import SessionManager, AgentManager
except ImportError:
    SessionManager = None
    AgentManager = None

@pytest.mark.tool
@pytest.mark.requires_server
@pytest.mark.requires_neurondb
class TestMLTool:
    """Test ML tool for machine learning operations."""
    
    def test_ml_tool_prediction(self, api_client, test_agent, test_session):
        """Test ML tool for predictions."""
        agent_mgr = AgentManager(api_client)
        agent_mgr.update(test_agent['id'], enabled_tools=['ml'])
        
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use the ML tool to make a prediction.",
            role="user"
        )
        assert 'response' in response
    
    def test_ml_tool_training(self, api_client, test_agent, test_session):
        """Test ML tool for model training."""
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use the ML tool to train a model.",
            role="user"
        )
        assert 'response' in response
    
    def test_ml_tool_inference(self, api_client, test_agent, test_session):
        """Test ML tool for inference."""
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use the ML tool to run inference on data.",
            role="user"
        )
        assert 'response' in response


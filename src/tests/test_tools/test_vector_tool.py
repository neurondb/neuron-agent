"""Tests for Vector Tool."""
import pytest
import sys
import os
try:
    import numpy as np
except ImportError:
    np = None
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import SessionManager, AgentManager
except ImportError:
    SessionManager = None
    AgentManager = None

@pytest.mark.tool
@pytest.mark.requires_server
@pytest.mark.requires_neurondb
class TestVectorTool:
    """Test vector tool operations."""
    
    def test_vector_tool_similarity(self, api_client, test_agent, test_session):
        """Test vector similarity calculation."""
        if SessionManager is None or AgentManager is None:
            pytest.skip("Client managers not available")
        session_mgr = SessionManager(api_client)
        agent_mgr = AgentManager(api_client)
        agent_mgr.update(test_agent['id'], enabled_tools=['vector'])
        
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use the vector tool to calculate similarity between two vectors.",
            role="user"
        )
        assert 'response' in response
    
    def test_vector_tool_embedding(self, api_client, test_agent, test_session):
        """Test vector embedding generation."""
        if SessionManager is None:
            pytest.skip("SessionManager not available")
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use the vector tool to generate an embedding for 'test text'.",
            role="user"
        )
        assert 'response' in response
    
    def test_vector_tool_search(self, api_client, test_agent, test_session):
        """Test vector similarity search."""
        if SessionManager is None:
            pytest.skip("SessionManager not available")
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use the vector tool to search for similar vectors.",
            role="user"
        )
        assert 'response' in response
    
    def test_vector_tool_operations(self, api_client, test_agent, test_session):
        """Test vector mathematical operations."""
        if SessionManager is None:
            pytest.skip("SessionManager not available")
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use the vector tool to perform vector addition or subtraction.",
            role="user"
        )
        assert 'response' in response

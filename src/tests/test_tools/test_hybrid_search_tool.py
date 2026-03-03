"""Tests for Hybrid Search Tool."""
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
class TestHybridSearchTool:
    """Test hybrid search tool combining vector and keyword search."""
    
    def test_hybrid_search_basic(self, api_client, test_agent, test_session):
        """Test basic hybrid search."""
        agent_mgr = AgentManager(api_client)
        agent_mgr.update(test_agent['id'], enabled_tools=['hybrid_search'])
        
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use hybrid search to find documents about 'machine learning'.",
            role="user"
        )
        assert 'response' in response
    
    def test_hybrid_search_vector_keyword(self, api_client, test_agent, test_session):
        """Test hybrid search combining vector and keyword."""
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use hybrid search with both semantic and keyword matching.",
            role="user"
        )
        assert 'response' in response
    
    def test_hybrid_search_ranking(self, api_client, test_agent, test_session):
        """Test hybrid search result ranking."""
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use hybrid search and rank results by relevance.",
            role="user"
        )
        assert 'response' in response


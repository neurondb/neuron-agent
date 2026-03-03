"""Tests for Reranking Tool."""
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
class TestRerankingTool:
    """Test reranking tool for search results."""
    
    def test_reranking_tool_basic(self, api_client, test_agent, test_session):
        """Test basic reranking functionality."""
        if AgentManager is None:
            pytest.skip("AgentManager not available")
        agent_mgr = AgentManager(api_client)
        agent_mgr.update(test_agent['id'], enabled_tools=['reranking'])
        
        if SessionManager is None:
            pytest.skip("SessionManager not available")
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use the reranking tool to rerank search results.",
            role="user"
        )
        assert 'response' in response
    
    def test_reranking_tool_cross_encoder(self, api_client, test_agent, test_session):
        """Test reranking with cross-encoder model."""
        if SessionManager is None:
            pytest.skip("SessionManager not available")
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use reranking with cross-encoder to improve result order.",
            role="user"
        )
        assert 'response' in response
    
    def test_reranking_tool_relevance(self, api_client, test_agent, test_session):
        """Test reranking by relevance score."""
        if SessionManager is None:
            pytest.skip("SessionManager not available")
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Rerank results by relevance to the query.",
            role="user"
        )
        assert 'response' in response

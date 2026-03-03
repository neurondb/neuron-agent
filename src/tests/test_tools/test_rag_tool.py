"""Tests for RAG (Retrieval Augmented Generation) Tool."""
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
class TestRAGTool:
    """Test RAG tool execution."""
    
    def test_rag_tool_basic_retrieval(self, api_client, test_agent, test_session):
        """Test basic RAG retrieval."""
        agent_mgr = AgentManager(api_client)
        agent_mgr.update(test_agent['id'], enabled_tools=['rag'])
        
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use the RAG tool to find information about vector databases.",
            role="user"
        )
        assert 'response' in response
        assert isinstance(response.get('response'), str)
    
    def test_rag_tool_with_context(self, api_client, test_agent, test_session):
        """Test RAG tool with context."""
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Search for information about machine learning using RAG.",
            role="user"
        )
        assert 'response' in response
    
    def test_rag_tool_no_results(self, api_client, test_agent, test_session):
        """Test RAG tool when no results found."""
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use RAG to find information about nonexistent topic xyzabc123.",
            role="user"
        )
        assert 'response' in response
    
    def test_rag_tool_hybrid_search(self, api_client, test_agent, test_session):
        """Test RAG tool with hybrid search."""
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use RAG with hybrid search to find relevant documents.",
            role="user"
        )
        assert 'response' in response


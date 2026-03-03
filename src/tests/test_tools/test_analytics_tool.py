"""Tests for Analytics Tool."""
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
class TestAnalyticsTool:
    """Test analytics tool for data analysis."""
    
    def test_analytics_tool_basic(self, api_client, test_agent, test_session):
        """Test basic analytics operations."""
        agent_mgr = AgentManager(api_client)
        agent_mgr.update(test_agent['id'], enabled_tools=['analytics'])
        
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use the analytics tool to analyze data.",
            role="user"
        )
        assert 'response' in response
    
    def test_analytics_tool_aggregation(self, api_client, test_agent, test_session):
        """Test analytics aggregation functions."""
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use analytics tool to aggregate data (sum, avg, count).",
            role="user"
        )
        assert 'response' in response
    
    def test_analytics_tool_statistics(self, api_client, test_agent, test_session):
        """Test analytics statistical functions."""
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use analytics tool to calculate statistics.",
            role="user"
        )
        assert 'response' in response


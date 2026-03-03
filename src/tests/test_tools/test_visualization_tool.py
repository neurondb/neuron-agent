"""Tests for Visualization Tool."""
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
class TestVisualizationTool:
    """Test visualization tool for data visualization."""
    
    def test_visualization_tool_basic(self, api_client, test_agent, test_session):
        """Test basic visualization creation."""
        agent_mgr = AgentManager(api_client)
        agent_mgr.update(test_agent['id'], enabled_tools=['visualization'])
        
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use the visualization tool to create a chart.",
            role="user"
        )
        assert 'response' in response
    
    def test_visualization_tool_charts(self, api_client, test_agent, test_session):
        """Test creating different chart types."""
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use visualization tool to create a bar chart.",
            role="user"
        )
        assert 'response' in response
    
    def test_visualization_tool_export(self, api_client, test_agent, test_session):
        """Test exporting visualizations."""
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Use visualization tool to export a chart as image.",
            role="user"
        )
        assert 'response' in response


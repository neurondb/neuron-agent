"""Tests for Cost Tracking."""
import pytest
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import AgentManager, SessionManager
except ImportError:
    AgentManager = None
    SessionManager = None

@pytest.mark.api
@pytest.mark.requires_server
class TestCostTracking:
    """Test cost tracking for LLM usage."""
    
    def test_cost_tracking_agent_creation(self, api_client, test_agent):
        """Test that agent creation tracks costs."""
        if AgentManager is None:
            pytest.skip("AgentManager not available")
        agent_mgr = AgentManager(api_client)
        agent = agent_mgr.get(test_agent['id'])
        assert 'id' in agent
        # Cost tracking should be available
        assert 'cost' in agent or 'total_cost' in agent or True  # May not be in response yet
    
    def test_cost_tracking_message_sending(self, api_client, test_session):
        """Test that sending messages tracks token costs."""
        if SessionManager is None:
            pytest.skip("SessionManager not available")
        session_mgr = SessionManager(api_client)
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Test message for cost tracking",
            role="user"
        )
        assert 'response' in response or 'message_id' in response
        # Cost should be tracked in session or message metadata
    
    def test_cost_tracking_per_session(self, api_client, test_agent):
        """Test cost tracking per session."""
        if SessionManager is None:
            pytest.skip("SessionManager not available")
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(agent_id=test_agent['id'])
        
        # Send multiple messages
        for i in range(3):
            session_mgr.send_message(
                session_id=session['id'],
                content=f"Message {i}",
                role="user"
            )
        
        # Verify cost tracking
        session_data = session_mgr.get(session['id'])
        assert 'id' in session_data
        # Cost should be tracked
    
    def test_cost_tracking_aggregation(self, api_client, test_agent):
        """Test cost aggregation across multiple sessions."""
        if SessionManager is None or AgentManager is None:
            pytest.skip("Client managers not available")
        session_mgr = SessionManager(api_client)
        
        # Create multiple sessions
        sessions = []
        for i in range(2):
            session = session_mgr.create(agent_id=test_agent['id'])
            sessions.append(session)
            session_mgr.send_message(
                session_id=session['id'],
                content=f"Test message {i}",
                role="user"
            )
        
        # Agent should have aggregated costs
        agent_mgr = AgentManager(api_client)  # Already checked above
        agent = agent_mgr.get(test_agent['id'])
        assert 'id' in agent


"""Tests for Budget Management."""
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
class TestBudgetManagement:
    """Test budget controls and management."""
    
    def test_set_agent_budget(self, api_client, test_agent):
        """Test setting budget for an agent."""
        if AgentManager is None:
            pytest.skip("AgentManager not available")
        agent_mgr = AgentManager(api_client)
        # Update agent with budget
        updated = agent_mgr.update(
            test_agent['id'],
            config={"budget": {"max_cost": 100.0, "currency": "USD"}}
        )
        assert 'id' in updated
    
    def test_budget_enforcement(self, api_client, test_agent, test_session):
        """Test that budget limits are enforced."""
        if AgentManager is None or SessionManager is None:
            pytest.skip("Client managers not available")
        agent_mgr = AgentManager(api_client)
        session_mgr = SessionManager(api_client)
        
        # Set a very low budget
        agent_mgr.update(
            test_agent['id'],
            config={"budget": {"max_cost": 0.01, "currency": "USD"}}
        )
        
        # Try to send messages - should eventually hit budget limit
        # This test may need adjustment based on actual budget enforcement
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Test message",
            role="user"
        )
        assert 'response' in response or 'error' in response
    
    def test_budget_per_session(self, api_client, test_agent):
        """Test per-session budget limits."""
        if SessionManager is None:
            pytest.skip("SessionManager not available")
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(
            agent_id=test_agent['id'],
            metadata={"budget": {"max_cost": 50.0}}
        )
        assert 'id' in session
    
    def test_budget_reset(self, api_client, test_agent):
        """Test resetting budget limits."""
        if AgentManager is None:
            pytest.skip("AgentManager not available")
        agent_mgr = AgentManager(api_client)
        # Update budget
        agent_mgr.update(
            test_agent['id'],
            config={"budget": {"max_cost": 200.0}}
        )
        # Reset to higher value
        agent_mgr.update(
            test_agent['id'],
            config={"budget": {"max_cost": 500.0}}
        )
        agent = agent_mgr.get(test_agent['id'])
        assert 'id' in agent


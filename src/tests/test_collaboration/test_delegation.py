"""
Comprehensive tests for Agent Delegation.

Tests delegating tasks to specialized agents.
"""
import pytest
import uuid
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import AgentManager, SessionManager
except ImportError:
    AgentManager = None
    SessionManager = None

@pytest.mark.collaboration
@pytest.mark.requires_server
class TestDelegation:
    """Test agent delegation functionality."""
    
    def test_delegate_to_specialized_agent(self, api_client, unique_name):
        """Test delegating task to a specialized agent."""
        agent_mgr = AgentManager(api_client)
        
        # Create specialized agent
        specialized = agent_mgr.create(
            name=f"{unique_name}-specialized",
            system_prompt="You are a SQL expert.",
            model_name="gpt-4",
            enabled_tools=['sql']
        )
        
        # Create delegating agent
        delegator = agent_mgr.create(
            name=f"{unique_name}-delegator",
            system_prompt="You delegate tasks to experts.",
            model_name="gpt-4",
            enabled_tools=['collaboration']
        )
        
        try:
            # Delegate task
            session_mgr = SessionManager(api_client)
            session = session_mgr.create(agent_id=delegator['id'])
            
            response = session_mgr.send_message(
                session_id=session['id'],
                content=f"Delegate this SQL query task to agent {specialized['id']}: SELECT 1",
                role="user"
            )
            
            assert 'response' in response or 'message_id' in response
        finally:
            agent_mgr.delete(specialized['id'])
            agent_mgr.delete(delegator['id'])
    
    def test_delegation_with_context(self, api_client, unique_name):
        """Test delegation with context passing."""
        agent_mgr = AgentManager(api_client)
        
        agent1 = agent_mgr.create(
            name=f"{unique_name}-agent1",
            system_prompt="Agent 1",
            model_name="gpt-4"
        )
        
        agent2 = agent_mgr.create(
            name=f"{unique_name}-agent2",
            system_prompt="Agent 2",
            model_name="gpt-4"
        )
        
        try:
            session_mgr = SessionManager(api_client)
            session = session_mgr.create(agent_id=agent1['id'])
            
            # Delegate with context
            response = session_mgr.send_message(
                session_id=session['id'],
                content=f"Delegate to {agent2['id']} with context: previous work was done",
                role="user"
            )
            
            assert 'response' in response
        finally:
            agent_mgr.delete(agent1['id'])
            agent_mgr.delete(agent2['id'])
    
    def test_delegation_error_handling(self, api_client, test_agent):
        """Test delegation error handling."""
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(agent_id=test_agent['id'])
        
        # Try to delegate to non-existent agent
        fake_id = str(uuid.uuid4())
        response = session_mgr.send_message(
            session_id=session['id'],
            content=f"Delegate to agent {fake_id}",
            role="user"
        )
        
        # Should handle error gracefully
        assert 'response' in response or 'error' in str(response).lower()

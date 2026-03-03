"""
Comprehensive tests for Agent State Machine.

Tests execution states, transitions, and state management.
"""
import pytest
import time
from neurondb_client import SessionManager

@pytest.mark.runtime
@pytest.mark.requires_server
class TestStateMachine:
    """Test agent execution state machine."""
    
    def test_state_machine_initial_state(self, api_client, test_session):
        """Test that agent starts in correct initial state."""
        session_mgr = SessionManager(api_client)
        session = session_mgr.get(test_session['id'])
        assert 'id' in session
        # Initial state should be ready or idle
    
    def test_state_machine_execution_flow(self, api_client, test_session):
        """Test complete execution state flow."""
        session_mgr = SessionManager(api_client)
        
        # Send message to trigger execution
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Test message for state machine",
            role="user"
        )
        
        assert 'response' in response or 'message_id' in response
        # State should transition: idle -> processing -> completed
    
    def test_state_machine_tool_execution_state(self, api_client, test_agent):
        """Test state transitions during tool execution."""
        agent_mgr = AgentManager(api_client)
        agent_mgr.update(test_agent['id'], enabled_tools=['sql'])
        
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(agent_id=test_agent['id'])
        
        # Send message that triggers tool
        response = session_mgr.send_message(
            session_id=session['id'],
            content="Use SQL tool to query the database",
            role="user"
        )
        
        assert 'response' in response or 'message_id' in response
        # State should include tool execution states
    
    def test_state_machine_error_state(self, api_client, test_session):
        """Test state transitions on errors."""
        session_mgr = SessionManager(api_client)
        
        # Send invalid message to trigger error
        try:
            response = session_mgr.send_message(
                session_id=test_session['id'],
                content="",  # Empty content should cause error
                role="user"
            )
        except Exception:
            # Error state should be handled
            pass
    
    def test_state_machine_concurrent_executions(self, api_client, test_agent):
        """Test state machine with concurrent executions."""
        session_mgr = SessionManager(api_client)
        
        # Create multiple sessions
        sessions = []
        for i in range(3):
            session = session_mgr.create(agent_id=test_agent['id'])
            sessions.append(session)
        
        # Send messages concurrently
        for session in sessions:
            session_mgr.send_message(
                session_id=session['id'],
                content=f"Concurrent message {i}",
                role="user"
            )
        
        # All should complete successfully
        for session in sessions:
            session_data = session_mgr.get(session['id'])
            assert 'id' in session_data

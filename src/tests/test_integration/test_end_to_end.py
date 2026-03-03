"""
Comprehensive tests for Complete Workflows.

Tests complete end-to-end scenarios: Agent → Session → Messages → Memory → Tools.
"""
import pytest
import time
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import AgentManager, SessionManager
except ImportError:
    AgentManager = None
    SessionManager = None

@pytest.mark.integration
@pytest.mark.requires_server
@pytest.mark.slow
class TestEndToEnd:
    """Test complete end-to-end workflows."""
    
    def test_complete_workflow(self, api_client, unique_name):
        """Test complete workflow: Agent → Session → Messages → Memory → Tools."""
        agent_mgr = AgentManager(api_client)
        session_mgr = SessionManager(api_client)
        
        # Create agent
        agent = agent_mgr.create(
            name=unique_name,
            system_prompt="You are a helpful assistant.",
            model_name="gpt-4",
            enabled_tools=['sql', 'http']
        )
        
        try:
            # Create session
            session = session_mgr.create(agent_id=agent['id'])
            
            try:
                # Send messages
                response1 = session_mgr.send_message(
                    session_id=session['id'],
                    content="Hello! My name is Alice.",
                    role="user"
                )
                assert 'response' in response1 or 'message_id' in response1
                
                time.sleep(1)
                
                # Send follow-up message
                response2 = session_mgr.send_message(
                    session_id=session['id'],
                    content="What is my name?",
                    role="user"
                )
                assert 'response' in response2
                
                # Verify memory/context
                messages = session_mgr.get_messages(session['id'])
                assert len(messages) >= 2
                
            finally:
                session_mgr.delete(session['id'])
        finally:
            agent_mgr.delete(agent['id'])
    
    def test_workflow_with_tool_usage(self, api_client, unique_name):
        """Test complete workflow with tool usage."""
        agent_mgr = AgentManager(api_client)
        session_mgr = SessionManager(api_client)
        
        agent = agent_mgr.create(
            name=unique_name,
            system_prompt="You can use SQL tool to query the database.",
            model_name="gpt-4",
            enabled_tools=['sql']
        )
        
        try:
            session = session_mgr.create(agent_id=agent['id'])
            
            try:
                # Send message that should trigger tool
                response = session_mgr.send_message(
                    session_id=session['id'],
                    content="Use SQL tool to show tables in neurondb_agent schema",
                    role="user"
                )
                
                assert 'response' in response
                # Response should indicate tool usage
            finally:
                session_mgr.delete(session['id'])
        finally:
            agent_mgr.delete(agent['id'])
    
    def test_workflow_multi_turn_conversation(self, api_client, unique_name):
        """Test multi-turn conversation workflow."""
        agent_mgr = AgentManager(api_client)
        session_mgr = SessionManager(api_client)
        
        agent = agent_mgr.create(
            name=unique_name,
            system_prompt="You maintain context across conversations.",
            model_name="gpt-4"
        )
        
        try:
            session = session_mgr.create(agent_id=agent['id'])
            
            try:
                # Multiple turns
                turns = [
                    "I like Python programming.",
                    "What programming language did I mention?",
                    "Can you write a simple Python function?",
                    "What was the first thing I said?"
                ]
                
                for turn in turns:
                    response = session_mgr.send_message(
                        session_id=session['id'],
                        content=turn,
                        role="user"
                    )
                    assert 'response' in response
                    time.sleep(0.5)
                
                # Agent should remember context
                messages = session_mgr.get_messages(session['id'])
                assert len(messages) >= len(turns) * 2  # User + assistant messages
                
            finally:
                session_mgr.delete(session['id'])
        finally:
            agent_mgr.delete(agent['id'])
    
    def test_workflow_error_recovery(self, api_client, unique_name):
        """Test workflow error recovery."""
        agent_mgr = AgentManager(api_client)
        session_mgr = SessionManager(api_client)
        
        agent = agent_mgr.create(
            name=unique_name,
            system_prompt="Test agent",
            model_name="gpt-4"
        )
        
        try:
            session = session_mgr.create(agent_id=agent['id'])
            
            try:
                # Send invalid message
                try:
                    response = session_mgr.send_message(
                        session_id=session['id'],
                        content="",  # Empty content
                        role="user"
                    )
                except Exception:
                    # Error expected
                    pass
                
                # Should still be able to send valid messages
                response = session_mgr.send_message(
                    session_id=session['id'],
                    content="Valid message after error",
                    role="user"
                )
                assert 'response' in response
                
            finally:
                session_mgr.delete(session['id'])
        finally:
            agent_mgr.delete(agent['id'])

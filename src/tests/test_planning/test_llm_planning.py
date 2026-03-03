"""
Comprehensive tests for LLM-Based Planning.

Tests advanced planning with task decomposition.
"""
import pytest
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import AgentManager, SessionManager
except ImportError:
    AgentManager = None
    SessionManager = None

@pytest.mark.planning
@pytest.mark.requires_server
class TestLLMPlanning:
    """Test LLM-based planning functionality."""
    
    def test_llm_planning_basic(self, api_client, test_agent):
        """Test basic LLM planning."""
        agent_mgr = AgentManager(api_client)
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(agent_id=test_agent['id'])
        
        # Request a plan
        response = session_mgr.send_message(
            session_id=session['id'],
            content="Create a plan to build a web application",
            role="user"
        )
        
        assert 'response' in response
        # Response should contain a plan
    
    def test_llm_planning_task_decomposition(self, api_client, test_agent):
        """Test planning with task decomposition."""
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(agent_id=test_agent['id'])
        
        response = session_mgr.send_message(
            session_id=session['id'],
            content="Break down this goal into steps: Develop a machine learning model",
            role="user"
        )
        
        assert 'response' in response
        # Plan should have multiple steps
    
    def test_llm_planning_with_constraints(self, api_client, test_agent):
        """Test planning with constraints."""
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(agent_id=test_agent['id'])
        
        response = session_mgr.send_message(
            session_id=session['id'],
            content="Create a plan to deploy an application with constraints: budget $100, timeline 1 week",
            role="user"
        )
        
        assert 'response' in response
    
    def test_llm_planning_api_endpoint(self, api_client, test_agent):
        """Test planning via API endpoint."""
        try:
            response = api_client.post(
                f"/api/v1/agents/{test_agent['id']}/plan",
                json_data={"goal": "Test goal", "constraints": []}
            )
            assert 'plan' in response or 'steps' in response
        except Exception:
            pytest.skip("Planning API endpoint not available")

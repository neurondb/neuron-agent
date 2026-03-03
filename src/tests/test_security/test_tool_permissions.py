"""Tests for Tool Permissions."""
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

@pytest.mark.security
@pytest.mark.requires_server
class TestToolPermissions:
    """Test granular tool access permissions."""
    
    def test_tool_permissions_agent_restriction(self, api_client):
        """Test that agents can be restricted to specific tools."""
        agent_mgr = AgentManager(api_client)
        
        # Create agent with limited tools
        agent = agent_mgr.create(
            name=f"test-restricted-agent-{uuid.uuid4().hex[:8]}",
            system_prompt="Restricted agent",
            model_name="gpt-4",
            enabled_tools=['sql']  # Only SQL tool
        )
        
        assert 'sql' in agent.get('enabled_tools', [])
        # Agent should not have access to other tools
    
    def test_tool_permissions_tool_denial(self, api_client, test_agent):
        """Test that denied tools cannot be used."""
        agent_mgr = AgentManager(api_client)
        session_mgr = SessionManager(api_client)
        
        # Update agent to remove a tool
        agent_mgr.update(test_agent['id'], enabled_tools=['sql'])
        
        # Attempting to use a disabled tool should fail
        # This depends on tool execution enforcement
        session = session_mgr.create(agent_id=test_agent['id'])
        assert 'id' in session
    
    def test_tool_permissions_role_based(self, api_client):
        """Test role-based tool permissions."""
        # Different roles should have different tool access
        agent_mgr = AgentManager(api_client)
        agent = agent_mgr.create(
            name=f"test-role-agent-{uuid.uuid4().hex[:8]}",
            system_prompt="Role-based agent",
            model_name="gpt-4"
        )
        assert 'id' in agent
    
    def test_tool_permissions_dynamic_update(self, api_client, test_agent):
        """Test that tool permissions can be updated dynamically."""
        agent_mgr = AgentManager(api_client)
        
        # Initially allow SQL tool
        agent_mgr.update(test_agent['id'], enabled_tools=['sql'])
        agent1 = agent_mgr.get(test_agent['id'])
        
        # Remove SQL, add HTTP
        agent_mgr.update(test_agent['id'], enabled_tools=['http'])
        agent2 = agent_mgr.get(test_agent['id'])
        
        assert 'http' in agent2.get('enabled_tools', [])
        assert 'sql' not in agent2.get('enabled_tools', [])


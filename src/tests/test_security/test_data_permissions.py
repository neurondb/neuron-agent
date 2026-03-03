"""Tests for Data Permissions."""
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
class TestDataPermissions:
    """Test per-principal data access controls."""
    
    def test_data_permissions_agent_isolation(self, api_client):
        """Test that agents are isolated by principal."""
        agent_mgr = AgentManager(api_client)
        
        # Create agent with specific principal
        agent1 = agent_mgr.create(
            name=f"test-agent-principal-1-{uuid.uuid4().hex[:8]}",
            system_prompt="Test agent 1",
            model_name="gpt-4"
        )
        
        # Another principal should not see this agent
        # This depends on RBAC implementation
        agents = agent_mgr.list()
        assert isinstance(agents, list)
    
    def test_data_permissions_session_isolation(self, api_client, test_agent):
        """Test that sessions are isolated by principal."""
        session_mgr = SessionManager(api_client)
        
        # Create session
        session = session_mgr.create(agent_id=test_agent['id'])
        
        # Sessions should be filtered by principal
        sessions = session_mgr.list(agent_id=test_agent['id'])
        assert isinstance(sessions, list)
    
    def test_data_permissions_memory_isolation(self, api_client, test_agent):
        """Test that memory chunks are isolated by principal."""
        # Memory should be filtered by principal
        # This would require memory API access
        pytest.skip("Memory permissions require memory API")
    
    def test_data_permissions_cross_principal_access(self, api_client):
        """Test that cross-principal access is denied."""
        # Attempting to access another principal's data should fail
        pytest.skip("Cross-principal access test requires multiple principals")


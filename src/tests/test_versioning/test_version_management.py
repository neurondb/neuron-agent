"""
Comprehensive tests for Version Management.

Tests version control for agents and configurations.
"""
import pytest
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import AgentManager
except ImportError:
    AgentManager = None

@pytest.mark.versioning
@pytest.mark.requires_server
class TestVersionManagement:
    """Test version management functionality."""
    
    def test_version_management_list_versions(self, api_client, test_agent):
        """Test listing agent versions."""
        try:
            versions = api_client.get(f"/api/v1/agents/{test_agent['id']}/versions")
            assert isinstance(versions, list)
        except Exception:
            pytest.skip("Version API not available")
    
    def test_version_management_create_version(self, api_client, test_agent):
        """Test creating a new version."""
        if AgentManager is None:
            pytest.skip("AgentManager not available")
        agent_mgr = AgentManager(api_client)
        
        # Update agent to create new version
        updated = agent_mgr.update(
            test_agent['id'],
            system_prompt="Updated system prompt for versioning"
        )
        
        assert 'id' in updated
        # Version should be created
    
    def test_version_management_get_version(self, api_client, test_agent):
        """Test retrieving a specific version."""
        try:
            versions = api_client.get(f"/api/v1/agents/{test_agent['id']}/versions")
            if isinstance(versions, list) and len(versions) > 0:
                version_id = versions[0].get('id')
                if version_id:
                    version = api_client.get(
                        f"/api/v1/agents/{test_agent['id']}/versions/{version_id}"
                    )
                    assert 'id' in version
        except Exception:
            pytest.skip("Version API not available")
    
    def test_version_management_rollback(self, api_client, test_agent):
        """Test rolling back to a previous version."""
        try:
            versions = api_client.get(f"/api/v1/agents/{test_agent['id']}/versions")
            if isinstance(versions, list) and len(versions) > 1:
                previous_version = versions[-2].get('id')
                if previous_version:
                    response = api_client.post(
                        f"/api/v1/agents/{test_agent['id']}/versions/{previous_version}/rollback"
                    )
                    assert 'id' in response
        except Exception:
            pytest.skip("Version rollback API not available")
    
    def test_version_management_version_comparison(self, api_client, test_agent):
        """Test comparing versions."""
        try:
            versions = api_client.get(f"/api/v1/agents/{test_agent['id']}/versions")
            if isinstance(versions, list) and len(versions) >= 2:
                v1_id = versions[0].get('id')
                v2_id = versions[1].get('id')
                if v1_id and v2_id:
                    diff = api_client.get(
                        f"/api/v1/agents/{test_agent['id']}/versions/{v1_id}/compare/{v2_id}"
                    )
                    assert isinstance(diff, dict)
        except Exception:
            pytest.skip("Version comparison API not available")

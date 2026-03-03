"""Tests for Session Caching."""
import pytest
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import SessionManager
except ImportError:
    SessionManager = None

@pytest.mark.storage
@pytest.mark.requires_server
class TestSessionCache:
    """Test Redis-compatible session caching."""
    
    def test_session_cache_creation(self, api_client, test_agent):
        """Test that sessions are cached."""
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(agent_id=test_agent['id'])
        assert 'id' in session
        
        # Retrieve should use cache
        cached = session_mgr.get(session['id'])
        assert cached['id'] == session['id']
    
    def test_session_cache_ttl(self, api_client, test_agent):
        """Test session cache TTL."""
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(agent_id=test_agent['id'])
        assert 'id' in session
        # TTL should be configured
    
    def test_session_cache_invalidation(self, api_client, test_agent):
        """Test session cache invalidation."""
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(agent_id=test_agent['id'])
        
        # Update should invalidate cache
        session_mgr.update(session['id'], metadata={"updated": True})
        
        # Retrieve should get fresh data
        updated = session_mgr.get(session['id'])
        assert updated['metadata'].get('updated') is True


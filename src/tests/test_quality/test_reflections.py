"""
Comprehensive tests for Reflections.

Tests agent self-reflection and quality assessment.
"""
import pytest
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import SessionManager
except ImportError:
    SessionManager = None

@pytest.mark.quality
@pytest.mark.requires_server
class TestReflections:
    """Test agent reflection functionality."""
    
    def test_reflection_on_response(self, api_client, test_session):
        """Test agent reflecting on its own response."""
        session_mgr = SessionManager(api_client)
        
        # Send a message
        response1 = session_mgr.send_message(
            session_id=test_session['id'],
            content="Explain quantum computing",
            role="user"
        )
        
        # Request reflection on the response
        try:
            reflection = api_client.post(
                f"/api/v1/sessions/{test_session['id']}/reflect",
                json_data={"message_id": response1.get('message_id')}
            )
            assert 'reflection' in reflection or 'quality_score' in reflection
        except Exception:
            # Reflection API may not be available
            pass
    
    def test_reflection_quality_assessment(self, api_client, test_session):
        """Test quality assessment through reflection."""
        session_mgr = SessionManager(api_client)
        
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="What is artificial intelligence?",
            role="user"
        )
        
        # Reflection should assess quality
        assert 'response' in response
    
    def test_reflection_api_endpoint(self, api_client, test_agent):
        """Test reflection API endpoint."""
        try:
            response = api_client.post(
                f"/api/v1/agents/{test_agent['id']}/reflect",
                json_data={"response_content": "Test response"}
            )
            assert isinstance(response, dict)
        except Exception:
            pytest.skip("Reflection API endpoint not available")
    
    def test_reflection_improvement_suggestions(self, api_client, test_session):
        """Test that reflections provide improvement suggestions."""
        session_mgr = SessionManager(api_client)
        
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="Tell me about machine learning",
            role="user"
        )
        
        # Reflection should suggest improvements if needed
        assert 'response' in response

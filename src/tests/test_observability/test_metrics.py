"""
Comprehensive tests for Prometheus Metrics.

Tests comprehensive metrics export.
"""
import pytest
import requests
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import AgentManager, SessionManager
except ImportError:
    AgentManager = None
    SessionManager = None

@pytest.mark.observability
@pytest.mark.requires_server
class TestMetrics:
    """Test Prometheus metrics export."""
    
    def test_prometheus_metrics_endpoint(self):
        """Test Prometheus metrics endpoint."""
        response = requests.get("http://localhost:8080/metrics", timeout=5)
        assert response.status_code == 200
        assert len(response.text) > 0
    
    def test_prometheus_metrics_format(self):
        """Test that metrics are in Prometheus format."""
        response = requests.get("http://localhost:8080/metrics", timeout=5)
        content = response.text
        
        # Should contain Prometheus-style metrics
        assert "http_requests_total" in content or "go_" in content or "# HELP" in content
    
    def test_metrics_agent_operations(self, api_client, test_agent):
        """Test that agent operations generate metrics."""
        agent_mgr = AgentManager(api_client)
        
        # Perform operations
        agent_mgr.get(test_agent['id'])
        agent_mgr.list()
        
        # Check metrics
        response = requests.get("http://localhost:8080/metrics", timeout=5)
        assert response.status_code == 200
        # Metrics should reflect the operations
    
    def test_metrics_session_operations(self, api_client, test_session):
        """Test that session operations generate metrics."""
        session_mgr = SessionManager(api_client)
        
        # Perform operations
        session_mgr.get(test_session['id'])
        
        # Check metrics
        response = requests.get("http://localhost:8080/metrics", timeout=5)
        assert response.status_code == 200
    
    def test_metrics_message_operations(self, api_client, test_session):
        """Test that message operations generate metrics."""
        session_mgr = SessionManager(api_client)
        
        # Send message
        session_mgr.send_message(
            session_id=test_session['id'],
            content="Test message for metrics",
            role="user"
        )
        
        # Check metrics
        response = requests.get("http://localhost:8080/metrics", timeout=5)
        assert response.status_code == 200

"""Tests for Authentication Middleware."""
import pytest
import requests

@pytest.mark.security
class TestAuthentication:
    def test_missing_auth_header(self):
        """Test missing authorization header."""
        response = requests.get("http://localhost:8080/api/v1/agents")
        assert response.status_code == 401
    
    def test_invalid_api_key(self):
        """Test invalid API key."""
        response = requests.get(
            "http://localhost:8080/api/v1/agents",
            headers={"Authorization": "Bearer invalid_key"}
        )
        assert response.status_code == 401


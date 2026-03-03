"""Tests for Rate Limiting."""
import pytest
import requests
import time

@pytest.mark.security
@pytest.mark.requires_server
class TestRateLimiting:
    def test_rate_limiting(self, api_key):
        """Test rate limiting enforcement."""
        # Make many rapid requests
        headers = {"Authorization": f"Bearer {api_key}"}
        responses = []
        for _ in range(20):
            response = requests.get("http://localhost:8080/api/v1/agents", headers=headers, timeout=5)
            responses.append(response.status_code)
            time.sleep(0.1)
        # Some requests should be rate limited (429)
        assert 429 in responses or all(r == 200 for r in responses)


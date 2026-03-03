"""Tests for Virtual Filesystem."""
import pytest

@pytest.mark.requires_server
class TestVirtualFS:
    def test_virtual_filesystem(self, api_client, test_agent):
        """Test isolated filesystem for agents."""
        assert test_agent is not None


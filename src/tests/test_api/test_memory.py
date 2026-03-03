"""
Comprehensive tests for Memory API endpoints.

Tests memory operations, search, and management.
"""

import pytest
import uuid
from typing import Dict, Any, List


@pytest.mark.api
@pytest.mark.requires_server
class TestMemoryCRUD:
    """Test memory chunk CRUD operations."""
    
    def test_list_memory_chunks(self, api_client, test_agent):
        """Test listing memory chunks for an agent."""
        response = api_client.get(f"/api/v1/agents/{test_agent['id']}/memory")
        assert isinstance(response, list)
    
    def test_get_memory_chunk(self, api_client, test_agent):
        """Test retrieving a specific memory chunk."""
        # First get list of chunks
        chunks = api_client.get(f"/api/v1/agents/{test_agent['id']}/memory")
        
        if chunks and len(chunks) > 0:
            chunk_id = chunks[0]["id"]
            response = api_client.get(f"/api/v1/memory/{chunk_id}")
            assert response["id"] == chunk_id
            assert "content" in response
    
    def test_delete_memory_chunk(self, api_client, test_agent):
        """Test deleting a memory chunk."""
        chunks = api_client.get(f"/api/v1/agents/{test_agent['id']}/memory")
        
        if chunks and len(chunks) > 0:
            chunk_id = chunks[0]["id"]
            api_client.delete(f"/api/v1/memory/{chunk_id}")
            
            # Verify deletion
            with pytest.raises(Exception):
                api_client.get(f"/api/v1/memory/{chunk_id}")


@pytest.mark.api
@pytest.mark.requires_server
class TestMemorySearch:
    """Test memory search functionality."""
    
    def test_search_memory(self, api_client, test_agent):
        """Test searching memory chunks."""
        search_data = {
            "query": "test query",
            "limit": 10
        }
        
        response = api_client.post(
            f"/api/v1/agents/{test_agent['id']}/memory/search",
            json_data=search_data
        )
        
        assert isinstance(response, list) or isinstance(response, dict)


@pytest.mark.api
@pytest.mark.requires_server
class TestMemorySummarization:
    """Test memory summarization."""
    
    def test_summarize_memory(self, api_client, test_agent):
        """Test summarizing agent memory."""
        response = api_client.post(f"/api/v1/agents/{test_agent['id']}/memory/summarize")
        assert isinstance(response, dict)
        # Should have summary content


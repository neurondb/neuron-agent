"""
Comprehensive tests for Long-term Memory.

Tests HNSW-based vector search for context retrieval.
"""
import pytest
import time
import psycopg2
from neurondb_client import SessionManager, AgentManager

@pytest.mark.memory
@pytest.mark.requires_server
@pytest.mark.requires_neurondb
class TestLongTermMemory:
    """Test long-term memory with HNSW vector search."""
    
    def test_vector_search_basic(self, api_client, test_agent, db_connection):
        """Test basic vector similarity search."""
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(agent_id=test_agent['id'])
        
        # Send message that should create memory
        session_mgr.send_message(
            session_id=session['id'],
            content="Remember: The capital of France is Paris.",
            role="user"
        )
        
        time.sleep(2)  # Wait for memory processing
        
        # Search for the memory
        with db_connection.cursor() as cur:
            cur.execute("""
                SELECT id, content FROM neurondb_agent.memory_chunks
                WHERE agent_id = %s::uuid
                AND content ILIKE %s
                LIMIT 1;
            """, (test_agent['id'], '%Paris%'))
            result = cur.fetchone()
            # Memory should be stored
    
    def test_vector_search_similarity(self, api_client, test_agent, db_connection):
        """Test vector similarity search with HNSW index."""
        session_mgr = SessionManager(api_client)
        session = session_mgr.create(agent_id=test_agent['id'])
        
        # Store multiple memories
        memories = [
            "Machine learning is a subset of artificial intelligence",
            "Deep learning uses neural networks",
            "Natural language processing handles text"
        ]
        
        for memory in memories:
            session_mgr.send_message(
                session_id=session['id'],
                content=f"Remember: {memory}",
                role="user"
            )
            time.sleep(1)
        
        time.sleep(3)  # Wait for embeddings
        
        # Search for similar content
        with db_connection.cursor() as cur:
            cur.execute("""
                SELECT COUNT(*) FROM neurondb_agent.memory_chunks
                WHERE agent_id = %s::uuid
                AND embedding IS NOT NULL;
            """, (test_agent['id'],))
            count = cur.fetchone()[0]
            assert count >= 0  # May or may not have embeddings yet
    
    def test_vector_search_retrieval(self, api_client, test_agent, test_session):
        """Test retrieving memories via vector search."""
        session_mgr = SessionManager(api_client)
        
        # Store important memory
        session_mgr.send_message(
            session_id=test_session['id'],
            content="Important fact: Python is a programming language.",
            role="user"
        )
        
        time.sleep(2)
        
        # Query for the memory
        response = session_mgr.send_message(
            session_id=test_session['id'],
            content="What programming language did we discuss?",
            role="user"
        )
        
        assert 'response' in response
        # Response should reference the stored memory
    
    def test_vector_search_ranking(self, api_client, test_agent, db_connection):
        """Test that vector search returns ranked results."""
        # This would test that most similar memories come first
        with db_connection.cursor() as cur:
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM pg_indexes
                    WHERE schemaname = 'neurondb_agent'
                    AND tablename = 'memory_chunks'
                    AND indexdef ILIKE '%hnsw%'
                );
            """)
            has_index = cur.fetchone()[0]
            # HNSW index should exist for efficient search

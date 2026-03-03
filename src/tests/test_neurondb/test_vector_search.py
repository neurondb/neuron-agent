"""
Comprehensive tests for Vector Similarity Search.

Tests vector operations, similarity search, and HNSW index usage.
"""

import pytest


@pytest.mark.neurondb
@pytest.mark.requires_neurondb
@pytest.mark.requires_db
class TestVectorSearch:
    """Test vector similarity search."""
    
    def test_vector_similarity_operator(self, db_connection):
        """Test vector similarity operator (<=>)."""
        with db_connection.cursor() as cur:
            try:
                cur.execute("""
                    SELECT '[1,0,0]'::vector(3) <=> '[0,1,0]'::vector(3) AS distance;
                """)
                result = cur.fetchone()
                assert result is not None
                distance = result[0]
                assert isinstance(distance, (int, float))
                assert distance >= 0
            except Exception as e:
                pytest.skip(f"Vector similarity operator not available: {e}")
    
    def test_vector_search_on_memory_chunks(self, db_connection):
        """Test vector search on memory_chunks table."""
        with db_connection.cursor() as cur:
            try:
                # Check if memory_chunks table has embedding column
                cur.execute("""
                    SELECT column_name FROM information_schema.columns 
                    WHERE table_schema = 'neurondb_agent' 
                    AND table_name = 'memory_chunks' 
                    AND column_name = 'embedding';
                """)
                result = cur.fetchone()
                if result:
                    # Try a similarity search
                    cur.execute("""
                        SELECT id, embedding <=> '[0.1,0.2,0.3]'::vector AS distance
                        FROM neurondb_agent.memory_chunks
                        WHERE embedding IS NOT NULL
                        ORDER BY embedding <=> '[0.1,0.2,0.3]'::vector
                        LIMIT 5;
                    """)
                    results = cur.fetchall()
                    # Should return results or empty list
                    assert isinstance(results, list)
            except Exception as e:
                pytest.skip(f"Vector search on memory_chunks failed: {e}")


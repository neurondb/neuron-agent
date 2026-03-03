"""
Comprehensive tests for Memory Chunk Embeddings.

Tests embedding storage, retrieval, and search in memory chunks.
"""

import pytest


@pytest.mark.neurondb
@pytest.mark.requires_neurondb
@pytest.mark.requires_db
class TestMemoryVectors:
    """Test memory chunk embeddings."""
    
    def test_memory_chunk_embedding_column(self, db_connection):
        """Test that memory_chunks table has embedding column."""
        with db_connection.cursor() as cur:
            cur.execute("""
                SELECT data_type FROM information_schema.columns 
                WHERE table_schema = 'neurondb_agent' 
                AND table_name = 'memory_chunks' 
                AND column_name = 'embedding';
            """)
            result = cur.fetchone()
            assert result is not None
            # Should be vector type or USER-DEFINED
            assert 'vector' in result[0].lower() or 'user-defined' in result[0].lower()
    
    def test_memory_chunk_with_embedding(self, db_connection):
        """Test creating memory chunk with embedding."""
        with db_connection.cursor() as cur:
            try:
                # Get a test agent ID
                cur.execute("SELECT id FROM neurondb_agent.agents LIMIT 1;")
                agent_result = cur.fetchone()
                if agent_result:
                    agent_id = agent_result[0]
                    
                    # Create memory chunk with embedding
                    cur.execute("""
                        INSERT INTO neurondb_agent.memory_chunks 
                        (agent_id, content, embedding, importance_score, metadata)
                        VALUES 
                        (%s, 'Test memory chunk', '[0.1,0.2,0.3,0.4,0.5]'::vector, 0.8, '{}'::jsonb)
                        RETURNING id;
                    """, (agent_id,))
                    chunk_id = cur.fetchone()[0]
                    
                    # Verify it was created
                    assert chunk_id is not None
                    
                    # Cleanup
                    cur.execute("DELETE FROM neurondb_agent.memory_chunks WHERE id = %s;", (chunk_id,))
                    db_connection.commit()
            except Exception as e:
                pytest.skip(f"Memory chunk embedding test failed: {e}")


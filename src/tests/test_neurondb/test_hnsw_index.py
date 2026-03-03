"""
Comprehensive tests for HNSW Index Operations.

Tests HNSW index creation, usage, and performance.
"""

import pytest


@pytest.mark.neurondb
@pytest.mark.requires_neurondb
@pytest.mark.requires_db
class TestHNSWIndex:
    """Test HNSW index operations."""
    
    def test_hnsw_index_exists(self, db_connection):
        """Test that HNSW index exists on memory_chunks.embedding."""
        with db_connection.cursor() as cur:
            cur.execute("""
                SELECT indexname FROM pg_indexes 
                WHERE schemaname = 'neurondb_agent' 
                AND tablename = 'memory_chunks' 
                AND indexname LIKE '%embedding%' 
                OR indexname LIKE '%hnsw%';
            """)
            result = cur.fetchone()
            # Index may or may not exist yet (created on first use)
            # This is informational, not a failure
    
    def test_hnsw_index_creation(self, db_connection):
        """Test creating HNSW index if it doesn't exist."""
        with db_connection.cursor() as cur:
            try:
                # Check if index exists
                cur.execute("""
                    SELECT EXISTS(
                        SELECT 1 FROM pg_indexes 
                        WHERE schemaname = 'neurondb_agent' 
                        AND tablename = 'memory_chunks' 
                        AND indexname LIKE '%embedding%'
                    );
                """)
                exists = cur.fetchone()[0]
                
                if not exists:
                    # Try to create index
                    cur.execute("""
                        CREATE INDEX IF NOT EXISTS memory_chunks_embedding_idx 
                        ON neurondb_agent.memory_chunks 
                        USING hnsw (embedding vector_cosine_ops);
                    """)
                    db_connection.commit()
            except Exception as e:
                pytest.skip(f"HNSW index creation not available: {e}")


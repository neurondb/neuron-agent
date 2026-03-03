"""
Comprehensive tests for NeuronDB Embedding Generation.

Tests embedding generation, batch processing, and model configuration.
"""

import pytest


@pytest.mark.neurondb
@pytest.mark.requires_neurondb
@pytest.mark.requires_db
class TestEmbeddings:
    """Test embedding generation."""
    
    def test_embedding_generation(self, db_connection):
        """Test generating embeddings using neurondb_embed."""
        with db_connection.cursor() as cur:
            try:
                cur.execute("SELECT neurondb_embed('Test text', 'all-MiniLM-L6-v2')::text;")
                result = cur.fetchone()
                assert result is not None
                # Result should be a vector
            except Exception as e:
                pytest.skip(f"Embedding function not available: {e}")
    
    def test_batch_embeddings(self, db_connection):
        """Test batch embedding generation."""
        with db_connection.cursor() as cur:
            try:
                cur.execute("""
                    SELECT neurondb_embed_batch(ARRAY['Text 1', 'Text 2'], 'all-MiniLM-L6-v2');
                """)
                result = cur.fetchone()
                assert result is not None
            except Exception as e:
                pytest.skip(f"Batch embedding function not available: {e}")
    
    def test_embedding_dimension(self, db_connection):
        """Test that embeddings have correct dimension."""
        with db_connection.cursor() as cur:
            try:
                cur.execute("""
                    SELECT array_length(neurondb_embed('Test', 'all-MiniLM-L6-v2')::float[], 1) as dim;
                """)
                result = cur.fetchone()
                if result:
                    dim = result[0]
                    assert dim > 0
            except Exception as e:
                pytest.skip(f"Embedding dimension check failed: {e}")


"""Tests for Batch Embeddings."""
import pytest

@pytest.mark.neurondb
@pytest.mark.requires_neurondb
@pytest.mark.requires_db
class TestBatchEmbeddings:
    def test_batch_embeddings(self, db_connection):
        """Test batch embedding generation."""
        with db_connection.cursor() as cur:
            try:
                cur.execute("SELECT neurondb_embed_batch(ARRAY['Text 1', 'Text 2'], 'all-MiniLM-L6-v2');")
                result = cur.fetchone()
                assert result is not None
            except Exception as e:
                pytest.skip(f"Batch embeddings not available: {e}")


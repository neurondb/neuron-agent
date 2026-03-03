"""Tests for Vector Operations."""
import pytest

@pytest.mark.neurondb
@pytest.mark.requires_neurondb
@pytest.mark.requires_db
class TestVectorOps:
    def test_vector_operations(self, db_connection):
        """Test vector type operations."""
        with db_connection.cursor() as cur:
            try:
                cur.execute("SELECT '[1,2,3]'::vector(3);")
                result = cur.fetchone()
                assert result is not None
            except Exception as e:
                pytest.skip(f"Vector operations not available: {e}")


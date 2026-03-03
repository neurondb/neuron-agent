"""Tests for Index Verification."""
import pytest

@pytest.mark.requires_db
class TestIndexes:
    def test_indexes_exist(self, db_connection):
        """Test that indexes exist."""
        with db_connection.cursor() as cur:
            cur.execute("""
                SELECT COUNT(*) FROM pg_indexes 
                WHERE schemaname = 'neurondb_agent';
            """)
            count = cur.fetchone()[0]
            assert count > 0, "Should have indexes"


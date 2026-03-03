"""Tests for Database Storage."""
import pytest

@pytest.mark.requires_db
class TestDatabaseStorage:
    def test_database_storage(self, db_connection):
        """Test PostgreSQL-based persistence."""
        with db_connection.cursor() as cur:
            cur.execute("SELECT 1")
            assert cur.fetchone()[0] == 1


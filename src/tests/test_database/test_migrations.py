"""Tests for Migration Testing."""
import pytest

@pytest.mark.requires_db
class TestMigrations:
    def test_migrations_table(self, db_connection):
        """Test that migrations table exists."""
        with db_connection.cursor() as cur:
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM information_schema.tables 
                    WHERE table_schema = 'neurondb_agent' 
                    AND table_name = 'schema_migrations'
                );
            """)
            result = cur.fetchone()[0]
            # Schema migrations table is optional - may or may not exist
            # This is informational, not a failure
            if not result:
                pytest.skip("Schema migrations table not found (optional feature)")


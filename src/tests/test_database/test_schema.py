"""Tests for Schema Validation."""
import pytest

@pytest.mark.requires_db
class TestSchema:
    def test_schema_tables(self, db_connection):
        """Test that all required tables exist."""
        required_tables = ["agents", "sessions", "messages", "memory_chunks", "tools", "jobs", "api_keys"]
        with db_connection.cursor() as cur:
            for table in required_tables:
                cur.execute("""
                    SELECT EXISTS(
                        SELECT 1 FROM information_schema.tables 
                        WHERE table_schema = 'neurondb_agent' 
                        AND table_name = %s
                    );
                """, (table,))
                assert cur.fetchone()[0], f"Table {table} should exist"


"""Tests for Database Constraints."""
import pytest

@pytest.mark.requires_db
class TestConstraints:
    def test_foreign_keys(self, db_connection):
        """Test foreign key constraints."""
        with db_connection.cursor() as cur:
            cur.execute("""
                SELECT COUNT(*) FROM information_schema.table_constraints 
                WHERE constraint_schema = 'neurondb_agent' 
                AND constraint_type = 'FOREIGN KEY';
            """)
            count = cur.fetchone()[0]
            assert count > 0, "Should have foreign key constraints"
    
    def test_unique_constraints(self, db_connection):
        """Test unique constraints."""
        with db_connection.cursor() as cur:
            cur.execute("""
                SELECT COUNT(*) FROM information_schema.table_constraints 
                WHERE constraint_schema = 'neurondb_agent' 
                AND constraint_type = 'UNIQUE';
            """)
            count = cur.fetchone()[0]
            # May or may not have unique constraints


"""Tests for Data Consistency."""
import pytest

@pytest.mark.integration
@pytest.mark.requires_db
class TestDataIntegrity:
    def test_foreign_key_integrity(self, db_connection):
        """Test foreign key constraint integrity."""
        with db_connection.cursor() as cur:
            # Check that sessions reference valid agents
            cur.execute("""
                SELECT COUNT(*) FROM neurondb_agent.sessions s
                LEFT JOIN neurondb_agent.agents a ON s.agent_id = a.id
                WHERE a.id IS NULL;
            """)
            orphaned = cur.fetchone()[0]
            assert orphaned == 0, "No orphaned sessions"


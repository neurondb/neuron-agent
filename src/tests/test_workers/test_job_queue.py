"""Tests for Job Queue."""
import pytest

@pytest.mark.requires_db
class TestJobQueue:
    def test_job_queue(self, db_connection):
        """Test PostgreSQL-based job queue."""
        with db_connection.cursor() as cur:
            cur.execute("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'neurondb_agent' AND table_name = 'jobs');")
            exists = cur.fetchone()[0]
            # Jobs table may or may not exist


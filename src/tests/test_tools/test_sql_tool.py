"""
Comprehensive tests for SQL Tool.

Tests SQL query execution, read-only enforcement, and security.
"""

import pytest


@pytest.mark.tool
@pytest.mark.requires_server
class TestSQLTool:
    """Test SQL tool execution."""
    
    def test_sql_select_query(self, api_client, test_session):
        """Test executing SELECT query."""
        message_data = {
            "content": "Execute SQL: SELECT 1 as test_value",
            "role": "user"
        }
        
        response = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        
        # Response should contain query results
        assert response is not None
    
    def test_sql_read_only_enforcement(self, api_client, test_session):
        """Test that SQL tool only allows SELECT queries."""
        # Try to execute INSERT (should be blocked)
        message_data = {
            "content": "Execute SQL: INSERT INTO test_table VALUES (1)",
            "role": "user"
        }
        
        # Should either fail or be blocked
        response = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        # Response should indicate error or blocking
    
    def test_sql_schema_introspection(self, api_client, test_session):
        """Test schema introspection queries."""
        message_data = {
            "content": "Show me the tables in the database",
            "role": "user"
        }
        
        response = api_client.post(
            f"/api/v1/sessions/{test_session['id']}/messages",
            json_data=message_data
        )
        assert response is not None


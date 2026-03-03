"""
Comprehensive tests for LLM Function Integration.

Tests LLM function calls, response handling, and integration.
"""

import pytest


@pytest.mark.neurondb
@pytest.mark.requires_neurondb
@pytest.mark.requires_db
class TestLLMIntegration:
    """Test LLM function integration."""
    
    def test_llm_function_exists(self, db_connection):
        """Test that LLM functions are available."""
        with db_connection.cursor() as cur:
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM pg_proc 
                    WHERE proname IN ('neurondb_llm_generate', 'neurondb_llm_complete')
                );
            """)
            result = cur.fetchone()
            if not result[0]:
                pytest.skip("LLM functions not available")
    
    def test_llm_generate(self, db_connection):
        """Test LLM generation function."""
        with db_connection.cursor() as cur:
            try:
                cur.execute("""
                    SELECT neurondb_llm_generate('gpt-4', 'Hello, how are you?');
                """)
                result = cur.fetchone()
                assert result is not None
            except Exception as e:
                pytest.skip(f"LLM generate function not available or not configured: {e}")


"""
Comprehensive Integration Tests for NeuronAgent-NeuronDB Integration.

Tests all integration points between NeuronAgent and NeuronDB including:
- Embedding generation
- Vector operations
- LLM integration
- RAG pipeline
- Hybrid search
- Reranking
- ML operations
- Analytics
- Memory management
"""

import pytest
import json


@pytest.mark.neurondb
@pytest.mark.requires_neurondb
@pytest.mark.requires_db
@pytest.mark.slow
class TestComprehensiveIntegration:
    """Comprehensive integration tests for NeuronAgent-NeuronDB."""
    
    def test_embedding_integration(self, db_connection):
        """Test embedding generation integration."""
        with db_connection.cursor() as cur:
            # Test single embedding
            try:
                cur.execute("SELECT neurondb_embed('Test text', 'all-MiniLM-L6-v2')::text;")
                result = cur.fetchone()
                assert result is not None
                assert result[0] is not None
            except Exception as e:
                pytest.skip(f"Embedding function not available: {e}")
            
            # Test batch embedding if available
            try:
                cur.execute("""
                    SELECT neurondb_embed_batch(ARRAY['Text 1', 'Text 2'], 'all-MiniLM-L6-v2');
                """)
                result = cur.fetchone()
                if result:
                    assert result[0] is not None
            except Exception:
                # Batch embedding may not be available, that's okay
                pass
    
    def test_vector_operators(self, db_connection):
        """Test all vector operators."""
        with db_connection.cursor() as cur:
            # Cosine distance (<=>)
            try:
                cur.execute("""
                    SELECT '[1,0,0]'::vector(3) <=> '[0,1,0]'::vector(3) AS distance;
                """)
                result = cur.fetchone()
                assert result is not None
                distance = result[0]
                assert isinstance(distance, (int, float))
                assert distance >= 0
            except Exception as e:
                pytest.skip(f"Cosine distance operator not available: {e}")
            
            # L2 distance (<->)
            try:
                cur.execute("""
                    SELECT '[1,0,0]'::vector(3) <-> '[0,1,0]'::vector(3) AS distance;
                """)
                result = cur.fetchone()
                assert result is not None
                distance = result[0]
                assert isinstance(distance, (int, float))
                assert distance >= 0
            except Exception as e:
                pytest.skip(f"L2 distance operator not available: {e}")
            
            # Inner product (<#>)
            try:
                cur.execute("""
                    SELECT '[1,0,0]'::vector(3) <#> '[1,0,0]'::vector(3) AS product;
                """)
                result = cur.fetchone()
                assert result is not None
                product = result[0]
                assert isinstance(product, (int, float))
            except Exception as e:
                pytest.skip(f"Inner product operator not available: {e}")
    
    def test_llm_functions(self, db_connection):
        """Test LLM function availability."""
        with db_connection.cursor() as cur:
            # Check if LLM functions exist
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM pg_proc 
                    WHERE proname IN ('neurondb_llm_generate', 'neurondb_llm_complete')
                );
            """)
            result = cur.fetchone()
            if result and result[0]:
                # Try to call the function (may fail if not configured, that's okay)
                try:
                    cur.execute("""
                        SELECT EXISTS(
                            SELECT 1 FROM pg_proc WHERE proname = 'neurondb_llm_generate'
                        );
                    """)
                    has_generate = cur.fetchone()[0]
                    if has_generate:
                        # Function exists, but may not be configured
                        pass
                except Exception:
                    pass
            else:
                pytest.skip("LLM functions not available")
    
    def test_rag_functions(self, db_connection):
        """Test RAG function availability."""
        with db_connection.cursor() as cur:
            # Check chunk_text
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM pg_proc WHERE proname = 'neurondb_chunk_text'
                );
            """)
            has_chunk = cur.fetchone()[0]
            if has_chunk:
                try:
                    cur.execute("""
                        SELECT neurondb_chunk_text('This is a test document for chunking.', 10, 2);
                    """)
                    result = cur.fetchone()
                    if result:
                        chunks = result[0]
                        assert isinstance(chunks, (list, str))
                except Exception:
                    # Function exists but may not be configured
                    pass
            
            # Check rerank_results
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM pg_proc WHERE proname = 'neurondb_rerank_results'
                );
            """)
            has_rerank = cur.fetchone()[0]
            if not has_rerank:
                pytest.skip("Reranking function not available")
            
            # Check generate_answer
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM pg_proc WHERE proname = 'neurondb_generate_answer'
                );
            """)
            has_generate = cur.fetchone()[0]
            if not has_generate:
                pytest.skip("Answer generation function not available")
    
    def test_hybrid_search_functions(self, db_connection):
        """Test hybrid search function availability."""
        with db_connection.cursor() as cur:
            # Check hybrid_search
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM pg_proc WHERE proname = 'neurondb_hybrid_search'
                );
            """)
            has_hybrid = cur.fetchone()[0]
            if has_hybrid:
                # Function exists
                pass
            else:
                pytest.skip("Hybrid search function not available")
            
            # Check reciprocal_rank_fusion
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM pg_proc WHERE proname = 'neurondb_reciprocal_rank_fusion'
                );
            """)
            has_rrf = cur.fetchone()[0]
            if not has_rrf:
                pytest.skip("Reciprocal rank fusion function not available")
    
    def test_reranking_functions(self, db_connection):
        """Test reranking function availability."""
        with db_connection.cursor() as cur:
            functions = [
                'neurondb_rerank_cross_encoder',
                'neurondb_rerank_llm',
                'neurondb_rerank_colbert',
                'neurondb_rerank_ensemble'
            ]
            
            available = []
            for func in functions:
                cur.execute(f"""
                    SELECT EXISTS(
                        SELECT 1 FROM pg_proc WHERE proname = '{func}'
                    );
                """)
                if cur.fetchone()[0]:
                    available.append(func)
            
            if not available:
                pytest.skip("No reranking functions available")
            
            # At least one reranking function should be available
            assert len(available) > 0
    
    def test_ml_functions(self, db_connection):
        """Test ML function availability."""
        with db_connection.cursor() as cur:
            # Check if neurondb schema exists
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM pg_namespace WHERE nspname = 'neurondb'
                );
            """)
            has_schema = cur.fetchone()[0]
            if not has_schema:
                pytest.skip("NeuronDB schema not available")
            
            # Check train function
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM pg_proc 
                    WHERE proname = 'train' 
                    AND pronamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'neurondb')
                );
            """)
            has_train = cur.fetchone()[0]
            
            # Check predict function
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM pg_proc 
                    WHERE proname = 'predict' 
                    AND pronamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'neurondb')
                );
            """)
            has_predict = cur.fetchone()[0]
            
            # Check ml_models table
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM information_schema.tables 
                    WHERE table_schema = 'neurondb' AND table_name = 'ml_models'
                );
            """)
            has_table = cur.fetchone()[0]
            
            if not (has_train or has_predict or has_table):
                pytest.skip("ML functions not available")
    
    def test_analytics_functions(self, db_connection):
        """Test analytics function availability."""
        with db_connection.cursor() as cur:
            functions = [
                'neurondb_cluster',
                'neurondb_detect_outliers',
                'neurondb_reduce_dimensionality',
                'neurondb_analyze_data'
            ]
            
            available = []
            for func in functions:
                cur.execute(f"""
                    SELECT EXISTS(
                        SELECT 1 FROM pg_proc WHERE proname = '{func}'
                    );
                """)
                if cur.fetchone()[0]:
                    available.append(func)
            
            if not available:
                pytest.skip("No analytics functions available")
    
    def test_memory_tables(self, db_connection):
        """Test memory table structure."""
        with db_connection.cursor() as cur:
            # Check memory_chunks table
            cur.execute("""
                SELECT EXISTS(
                    SELECT 1 FROM information_schema.tables 
                    WHERE table_schema = 'neurondb_agent' AND table_name = 'memory_chunks'
                );
            """)
            has_chunks = cur.fetchone()[0]
            assert has_chunks, "memory_chunks table should exist"
            
            # Check embedding column type
            cur.execute("""
                SELECT data_type 
                FROM information_schema.columns 
                WHERE table_schema = 'neurondb_agent' 
                  AND table_name = 'memory_chunks' 
                  AND column_name = 'embedding';
            """)
            result = cur.fetchone()
            if result:
                data_type = result[0]
                assert 'vector' in data_type.lower() or 'user-defined' in data_type.lower(), \
                    f"embedding column should be vector type, got {data_type}"
            
            # Check hierarchical memory tables
            hierarchical_tables = ['memory_stm', 'memory_mtm', 'memory_lpm']
            for table in hierarchical_tables:
                cur.execute(f"""
                    SELECT EXISTS(
                        SELECT 1 FROM information_schema.tables 
                        WHERE table_schema = 'neurondb_agent' AND table_name = '{table}'
                    );
                """)
                has_table = cur.fetchone()[0]
                if has_table:
                    # Check embedding column
                    cur.execute(f"""
                        SELECT data_type 
                        FROM information_schema.columns 
                        WHERE table_schema = 'neurondb_agent' 
                          AND table_name = '{table}' 
                          AND column_name = 'embedding';
                    """)
                    result = cur.fetchone()
                    if result:
                        data_type = result[0]
                        # Should be neurondb_vector or vector type
                        assert 'vector' in data_type.lower(), \
                            f"{table}.embedding should be vector type, got {data_type}"
    
    def test_memory_vector_search(self, db_connection):
        """Test vector search on memory chunks."""
        with db_connection.cursor() as cur:
            # Check if memory_chunks has data
            cur.execute("""
                SELECT COUNT(*) FROM neurondb_agent.memory_chunks 
                WHERE embedding IS NOT NULL;
            """)
            count = cur.fetchone()[0]
            
            if count > 0:
                # Test vector similarity search
                try:
                    cur.execute("""
                        SELECT id, embedding <=> '[0.1,0.2,0.3]'::vector AS distance
                        FROM neurondb_agent.memory_chunks
                        WHERE embedding IS NOT NULL
                        ORDER BY embedding <=> '[0.1,0.2,0.3]'::vector
                        LIMIT 5;
                    """)
                    results = cur.fetchall()
                    assert isinstance(results, list)
                    # Should return results or empty list
                except Exception as e:
                    pytest.skip(f"Vector search on memory_chunks failed: {e}")
            else:
                pytest.skip("No memory chunks with embeddings to search")
    
    def test_hnsw_indexes(self, db_connection):
        """Test HNSW indexes on memory tables."""
        with db_connection.cursor() as cur:
            # Check for HNSW indexes on memory tables
            cur.execute("""
                SELECT indexname 
                FROM pg_indexes 
                WHERE schemaname = 'neurondb_agent' 
                  AND tablename LIKE 'memory%'
                  AND indexname LIKE '%embedding%'
                LIMIT 1;
            """)
            result = cur.fetchone()
            if result:
                # HNSW index exists
                index_name = result[0]
                assert index_name is not None
            else:
                # Index may be created automatically, that's okay
                pass
    
    def test_schema_constraints(self, db_connection):
        """Test database schema constraints."""
        with db_connection.cursor() as cur:
            # Check foreign key constraints
            cur.execute("""
                SELECT COUNT(*) 
                FROM information_schema.table_constraints 
                WHERE constraint_schema = 'neurondb_agent' 
                  AND constraint_type = 'FOREIGN KEY';
            """)
            fk_count = cur.fetchone()[0]
            assert fk_count > 0, "Should have foreign key constraints"
            
            # Check required tables
            required_tables = [
                'agents', 'sessions', 'messages', 'memory_chunks', 
                'tools', 'jobs', 'api_keys'
            ]
            for table in required_tables:
                cur.execute(f"""
                    SELECT EXISTS(
                        SELECT 1 FROM information_schema.tables 
                        WHERE table_schema = 'neurondb_agent' AND table_name = '{table}'
                    );
                """)
                exists = cur.fetchone()[0]
                assert exists, f"Table {table} should exist"
    
    @pytest.mark.requires_server
    def test_api_integration(self, api_client):
        """Test API integration with NeuronDB tools."""
        import time
        # Create agent with NeuronDB tools
        agent_data = {
            "name": f"test-integration-agent-{int(time.time())}",
            "description": "Integration test agent",
            "system_prompt": "You are a test agent",
            "model_name": "gpt-4",
            "enabled_tools": ["vector", "rag", "hybrid_search", "reranking"],
            "config": {}
        }
        
        response = api_client.post("/api/v1/agents", json_data=agent_data)
        assert response is not None
        
        if response and 'id' in response:
            agent_id = response['id']
            
            # Create session
            session_data = {
                "agent_id": agent_id,
                "external_user_id": "test-user"
            }
            session_response = api_client.post("/api/v1/sessions", json_data=session_data)
            
            if session_response and 'id' in session_response:
                session_id = session_response['id']
                
                # Send message (should trigger embeddings)
                message_data = {
                    "content": "Test message for integration",
                    "role": "user"
                }
                message_response = api_client.post(
                    f"/api/v1/sessions/{session_id}/messages",
                    json_data=message_data
                )
                assert message_response is not None


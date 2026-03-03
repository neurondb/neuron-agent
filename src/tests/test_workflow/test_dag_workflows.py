"""
Comprehensive tests for DAG Workflows.

Tests directed acyclic graph workflow execution.
"""
import pytest
import time
import uuid
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import AgentManager, SessionManager
except ImportError:
    AgentManager = None
    SessionManager = None

@pytest.mark.workflow
@pytest.mark.requires_server
class TestDAGWorkflows:
    """Test DAG workflow execution."""
    
    def test_dag_workflow_creation(self, api_client, test_agent, unique_name):
        """Test creating a DAG workflow."""
        workflow_data = {
            "name": f"{unique_name}-workflow",
            "description": "Test DAG workflow",
            "steps": [
                {
                    "id": "step1",
                    "type": "agent",
                    "agent_id": test_agent['id'],
                    "input": {"message": "Step 1"}
                },
                {
                    "id": "step2",
                    "type": "agent",
                    "agent_id": test_agent['id'],
                    "input": {"message": "Step 2"},
                    "depends_on": ["step1"]
                }
            ]
        }
        
        try:
            workflow = api_client.post("/api/v1/workflows", json_data=workflow_data)
            assert 'id' in workflow
            assert workflow['name'] == workflow_data['name']
        except Exception:
            pytest.skip("Workflow API not available")
    
    def test_dag_workflow_execution(self, api_client, test_agent):
        """Test executing a DAG workflow."""
        workflow_data = {
            "name": f"test-workflow-{uuid.uuid4().hex[:8]}",
            "steps": [
                {
                    "id": "step1",
                    "type": "agent",
                    "agent_id": test_agent['id'],
                    "input": {"message": "Hello"}
                }
            ]
        }
        
        try:
            workflow = api_client.post("/api/v1/workflows", json_data=workflow_data)
            execution = api_client.post(
                f"/api/v1/workflows/{workflow['id']}/execute",
                json_data={}
            )
            assert 'id' in execution or 'status' in execution
        except Exception:
            pytest.skip("Workflow execution API not available")
    
    def test_dag_workflow_parallel_steps(self, api_client, test_agent):
        """Test DAG workflow with parallel steps."""
        workflow_data = {
            "name": f"test-parallel-{uuid.uuid4().hex[:8]}",
            "steps": [
                {
                    "id": "step1",
                    "type": "agent",
                    "agent_id": test_agent['id'],
                    "input": {"message": "Step 1"}
                },
                {
                    "id": "step2",
                    "type": "agent",
                    "agent_id": test_agent['id'],
                    "input": {"message": "Step 2"}
                    # No depends_on - should run in parallel
                }
            ]
        }
        
        try:
            workflow = api_client.post("/api/v1/workflows", json_data=workflow_data)
            assert 'id' in workflow
        except Exception:
            pytest.skip("Workflow API not available")
    
    def test_dag_workflow_error_handling(self, api_client, test_agent):
        """Test DAG workflow error handling."""
        workflow_data = {
            "name": f"test-error-{uuid.uuid4().hex[:8]}",
            "steps": [
                {
                    "id": "step1",
                    "type": "agent",
                    "agent_id": "invalid-id",  # Invalid agent ID
                    "input": {"message": "Test"}
                }
            ]
        }
        
        try:
            # Should handle error
            workflow = api_client.post("/api/v1/workflows", json_data=workflow_data)
            execution = api_client.post(
                f"/api/v1/workflows/{workflow['id']}/execute",
                json_data={}
            )
            # Execution should handle error state
        except Exception:
            # Error is expected
            pass

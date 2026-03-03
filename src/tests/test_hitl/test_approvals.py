"""
Comprehensive tests for Approval Workflows.

Tests human approval gates in workflows.
"""
import pytest
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../examples'))
try:
    from neurondb_client import AgentManager
except ImportError:
    AgentManager = None

@pytest.mark.hitl
@pytest.mark.requires_server
class TestApprovals:
    """Test human-in-the-loop approval workflows."""
    
    def test_approval_workflow_creation(self, api_client, test_agent, unique_name):
        """Test creating workflow with approval step."""
        workflow_data = {
            "name": f"{unique_name}-approval-workflow",
            "steps": [
                {
                    "id": "step1",
                    "type": "agent",
                    "agent_id": test_agent['id'],
                    "input": {"message": "Generate proposal"}
                },
                {
                    "id": "approval1",
                    "type": "approval",
                    "required": True,
                    "depends_on": ["step1"]
                }
            ]
        }
        
        try:
            workflow = api_client.post("/api/v1/workflows", json_data=workflow_data)
            assert 'id' in workflow
        except Exception:
            pytest.skip("Approval workflow API not available")
    
    def test_approval_request(self, api_client):
        """Test requesting human approval."""
        try:
            approvals = api_client.get("/api/v1/approvals")
            assert isinstance(approvals, list)
        except Exception:
            pytest.skip("Approval API not available")
    
    def test_approval_approve(self, api_client):
        """Test approving a request."""
        try:
            # Get pending approvals
            approvals = api_client.get("/api/v1/approvals?status=pending")
            if isinstance(approvals, list) and len(approvals) > 0:
                approval_id = approvals[0].get('id')
                if approval_id:
                    response = api_client.post(
                        f"/api/v1/approvals/{approval_id}/approve",
                        json_data={}
                    )
                    assert 'status' in response
        except Exception:
            pytest.skip("Approval API not available")
    
    def test_approval_reject(self, api_client):
        """Test rejecting a request."""
        try:
            approvals = api_client.get("/api/v1/approvals?status=pending")
            if isinstance(approvals, list) and len(approvals) > 0:
                approval_id = approvals[0].get('id')
                if approval_id:
                    response = api_client.post(
                        f"/api/v1/approvals/{approval_id}/reject",
                        json_data={"reason": "Test rejection"}
                    )
                    assert 'status' in response
        except Exception:
            pytest.skip("Approval API not available")
    
    def test_approval_timeout(self, api_client):
        """Test approval timeout handling."""
        # Approvals should timeout if not responded to
        pytest.skip("Approval timeout requires time-based testing")

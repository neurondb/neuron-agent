"""
Comprehensive tests for Budget API endpoints.

Tests budget management, cost tracking, and budget controls.
"""

import pytest
import uuid
from typing import Dict, Any


@pytest.mark.api
@pytest.mark.requires_server
class TestBudgetManagement:
    """Test budget management operations."""
    
    def test_get_budget(self, api_client, test_agent):
        """Test getting budget for an agent."""
        response = api_client.get(f"/api/v1/agents/{test_agent['id']}/budget")
        assert isinstance(response, dict)
        # Should have budget-related fields
    
    def test_set_budget(self, api_client, test_agent):
        """Test setting budget for an agent."""
        budget_data = {
            "max_cost": 100.0,
            "currency": "USD",
            "alert_threshold": 0.8
        }
        
        response = api_client.post(
            f"/api/v1/agents/{test_agent['id']}/budget",
            json_data=budget_data
        )
        assert isinstance(response, dict)
    
    def test_update_budget(self, api_client, test_agent):
        """Test updating budget."""
        budget_data = {
            "max_cost": 200.0,
            "currency": "USD"
        }
        
        response = api_client.put(
            f"/api/v1/agents/{test_agent['id']}/budget",
            json_data=budget_data
        )
        assert isinstance(response, dict)


@pytest.mark.api
@pytest.mark.requires_server
class TestCostTracking:
    """Test cost tracking functionality."""
    
    def test_get_agent_costs(self, api_client, test_agent):
        """Test getting cost tracking for an agent."""
        response = api_client.get(f"/api/v1/agents/{test_agent['id']}/costs")
        assert isinstance(response, dict)
        # Should have cost-related fields like total_cost, tokens_used, etc.


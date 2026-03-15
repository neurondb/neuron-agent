/*-------------------------------------------------------------------------
 *
 * neuronagent_rbac_schema.sql
 *    RBAC expansion: principal tool permissions, workflow permissions, workspace policies.
 *
 * Idempotent. Run after neuron-agent.sql (and tenancy if used).
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 *-------------------------------------------------------------------------
 */

-- Principal-level tool permissions (allow/deny per principal, tool name)
CREATE TABLE IF NOT EXISTS neurondb_agent.principal_tool_permissions (
    principal_id UUID NOT NULL REFERENCES neurondb_agent.principals(id) ON DELETE CASCADE,
    tool_name TEXT NOT NULL,
    allowed BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (principal_id, tool_name)
);
COMMENT ON TABLE neurondb_agent.principal_tool_permissions IS 'Per-principal tool allow/deny for RBAC';

CREATE INDEX IF NOT EXISTS idx_principal_tool_permissions_principal ON neurondb_agent.principal_tool_permissions(principal_id);
CREATE INDEX IF NOT EXISTS idx_principal_tool_permissions_tool ON neurondb_agent.principal_tool_permissions(tool_name);

-- Workflow-level permissions (principal, workflow, role)
CREATE TABLE IF NOT EXISTS neurondb_agent.workflow_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    principal_id UUID NOT NULL REFERENCES neurondb_agent.principals(id) ON DELETE CASCADE,
    workflow_id UUID NOT NULL REFERENCES neurondb_agent.workflows(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('viewer', 'editor', 'executor', 'owner')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(principal_id, workflow_id)
);
COMMENT ON TABLE neurondb_agent.workflow_permissions IS 'Per-principal workflow role for RBAC';

CREATE INDEX IF NOT EXISTS idx_workflow_permissions_principal ON neurondb_agent.workflow_permissions(principal_id);
CREATE INDEX IF NOT EXISTS idx_workflow_permissions_workflow ON neurondb_agent.workflow_permissions(workflow_id);

-- Workspace policies (workspace-level policy config)
CREATE TABLE IF NOT EXISTS neurondb_agent.workspace_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL,
    policy_type TEXT NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE neurondb_agent.workspace_policies IS 'Workspace-level policy configuration for RBAC/ABAC';

CREATE INDEX IF NOT EXISTS idx_workspace_policies_workspace ON neurondb_agent.workspace_policies(workspace_id);
CREATE INDEX IF NOT EXISTS idx_workspace_policies_type ON neurondb_agent.workspace_policies(policy_type);

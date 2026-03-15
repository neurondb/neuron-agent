/*-------------------------------------------------------------------------
 *
 * neuronagent_audit_v2_schema.sql
 *    Dedicated audit tables for enterprise: events, policy_decisions,
 *    tool_invocations, workflow_executions, approval_actions.
 *    Partitioning-ready and retention-friendly.
 *
 * Idempotent. Run after neuron-agent.sql (and tenancy if used).
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 *-------------------------------------------------------------------------
 */

-- audit_events: general event log (request_id, actor_id, workspace_id, agent_id, timestamp, decision, metadata)
-- For retention: add partitioning by month in operations runbook (e.g. PARTITION BY RANGE (created_at)).
CREATE TABLE IF NOT EXISTS neurondb_agent.audit_events (
    id BIGSERIAL PRIMARY KEY,
    request_id TEXT,
    actor_id UUID,
    workspace_id UUID,
    agent_id UUID,
    session_id UUID,
    event_type TEXT NOT NULL,
    decision TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE neurondb_agent.audit_events IS 'Enterprise audit events; add monthly partitions for retention (see operations runbook)';
CREATE INDEX IF NOT EXISTS idx_audit_events_request_id ON neurondb_agent.audit_events(request_id, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_events_actor_id ON neurondb_agent.audit_events(actor_id, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_events_workspace_id ON neurondb_agent.audit_events(workspace_id, created_at) WHERE workspace_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_audit_events_agent_id ON neurondb_agent.audit_events(agent_id, created_at) WHERE agent_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_audit_events_created_at ON neurondb_agent.audit_events(created_at DESC);

-- policy_decisions: SQL/tool policy allow/deny with reason
CREATE TABLE IF NOT EXISTS neurondb_agent.policy_decisions (
    id BIGSERIAL PRIMARY KEY,
    request_id TEXT,
    actor_id UUID,
    workspace_id UUID,
    agent_id UUID,
    tool_name TEXT,
    decision TEXT NOT NULL CHECK (decision IN ('allowed', 'blocked')),
    reason_code TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE neurondb_agent.policy_decisions IS 'Policy engine decisions for audit';
CREATE INDEX IF NOT EXISTS idx_policy_decisions_request_id ON neurondb_agent.policy_decisions(request_id);
CREATE INDEX IF NOT EXISTS idx_policy_decisions_created_at ON neurondb_agent.policy_decisions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_policy_decisions_decision ON neurondb_agent.policy_decisions(decision, created_at DESC);

-- tool_invocations: each tool call with request_id, actor, workspace, agent, timestamp, metadata
CREATE TABLE IF NOT EXISTS neurondb_agent.tool_invocations (
    id BIGSERIAL PRIMARY KEY,
    request_id TEXT,
    actor_id UUID,
    workspace_id UUID,
    agent_id UUID,
    session_id UUID,
    tool_name TEXT NOT NULL,
    success BOOLEAN NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE neurondb_agent.tool_invocations IS 'Tool invocation audit trail';
CREATE INDEX IF NOT EXISTS idx_tool_invocations_request_id ON neurondb_agent.tool_invocations(request_id);
CREATE INDEX IF NOT EXISTS idx_tool_invocations_agent_id ON neurondb_agent.tool_invocations(agent_id, created_at DESC) WHERE agent_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_tool_invocations_tool_name ON neurondb_agent.tool_invocations(tool_name, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tool_invocations_created_at ON neurondb_agent.tool_invocations(created_at DESC);

-- workflow_executions: workflow run audit
CREATE TABLE IF NOT EXISTS neurondb_agent.workflow_executions_audit (
    id BIGSERIAL PRIMARY KEY,
    request_id TEXT,
    actor_id UUID,
    workspace_id UUID,
    workflow_id UUID NOT NULL,
    execution_id UUID NOT NULL,
    status TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE neurondb_agent.workflow_executions_audit IS 'Workflow execution audit trail';
CREATE INDEX IF NOT EXISTS idx_workflow_executions_audit_request_id ON neurondb_agent.workflow_executions_audit(request_id);
CREATE INDEX IF NOT EXISTS idx_workflow_executions_audit_workflow_id ON neurondb_agent.workflow_executions_audit(workflow_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_workflow_executions_audit_created_at ON neurondb_agent.workflow_executions_audit(created_at DESC);

-- approval_actions: HITL approve/reject
CREATE TABLE IF NOT EXISTS neurondb_agent.approval_actions (
    id BIGSERIAL PRIMARY KEY,
    request_id TEXT,
    actor_id UUID,
    workspace_id UUID,
    approval_id UUID NOT NULL,
    decision TEXT NOT NULL CHECK (decision IN ('approved', 'rejected')),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE neurondb_agent.approval_actions IS 'Approval gate decisions for audit';
CREATE INDEX IF NOT EXISTS idx_approval_actions_request_id ON neurondb_agent.approval_actions(request_id);
CREATE INDEX IF NOT EXISTS idx_approval_actions_approval_id ON neurondb_agent.approval_actions(approval_id);
CREATE INDEX IF NOT EXISTS idx_approval_actions_created_at ON neurondb_agent.approval_actions(created_at DESC);

-- Retention policy: document in docs/operations_runbook.txt. Example:
--   Monthly partitions for audit_events; drop partitions older than retention_days.
--   DELETE FROM neurondb_agent.audit_events WHERE created_at < NOW() - INTERVAL '90 days';
--   (Or convert table to partitioned and use partition drop.)

/*-------------------------------------------------------------------------
 *
 * neuronagent_tenancy_schema.sql
 *    Multi-tenancy: organizations, org_quotas, org_id on resources.
 *
 * Idempotent. Run after neuron-agent.sql.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 *-------------------------------------------------------------------------
 */

-- Organizations table
CREATE TABLE IF NOT EXISTS neurondb_agent.organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    billing_account_id TEXT,
    config JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE neurondb_agent.organizations IS 'Tenant organizations for multi-tenant isolation';

-- Org quotas: per-org resource limits
CREATE TABLE IF NOT EXISTS neurondb_agent.org_quotas (
    org_id UUID NOT NULL REFERENCES neurondb_agent.organizations(id) ON DELETE CASCADE,
    max_agents INT,
    max_sessions INT,
    max_workflows INT,
    max_tools INT,
    max_memory_mb INT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (org_id)
);
COMMENT ON TABLE neurondb_agent.org_quotas IS 'Per-organization resource quotas';

-- Add org_id to agents (nullable for backward compatibility)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'neurondb_agent' AND table_name = 'agents' AND column_name = 'org_id'
    ) THEN
        ALTER TABLE neurondb_agent.agents
        ADD COLUMN org_id UUID REFERENCES neurondb_agent.organizations(id) ON DELETE SET NULL;
    END IF;
END $$;
CREATE INDEX IF NOT EXISTS idx_agents_org_id ON neurondb_agent.agents(org_id) WHERE org_id IS NOT NULL;

-- Add org_id to sessions
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'neurondb_agent' AND table_name = 'sessions' AND column_name = 'org_id'
    ) THEN
        ALTER TABLE neurondb_agent.sessions
        ADD COLUMN org_id UUID REFERENCES neurondb_agent.organizations(id) ON DELETE SET NULL;
    END IF;
END $$;
CREATE INDEX IF NOT EXISTS idx_sessions_org_id ON neurondb_agent.sessions(org_id) WHERE org_id IS NOT NULL;

-- Add org_id to tools
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'neurondb_agent' AND table_name = 'tools' AND column_name = 'org_id'
    ) THEN
        ALTER TABLE neurondb_agent.tools
        ADD COLUMN org_id UUID REFERENCES neurondb_agent.organizations(id) ON DELETE SET NULL;
    END IF;
END $$;
CREATE INDEX IF NOT EXISTS idx_tools_org_id ON neurondb_agent.tools(org_id) WHERE org_id IS NOT NULL;

-- Add org_id to workflows
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'neurondb_agent' AND table_name = 'workflows' AND column_name = 'org_id'
    ) THEN
        ALTER TABLE neurondb_agent.workflows
        ADD COLUMN org_id UUID REFERENCES neurondb_agent.organizations(id) ON DELETE SET NULL;
    END IF;
END $$;
CREATE INDEX IF NOT EXISTS idx_workflows_org_id ON neurondb_agent.workflows(org_id) WHERE org_id IS NOT NULL;

-- Add org_id to audit_log
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'neurondb_agent' AND table_name = 'audit_log' AND column_name = 'org_id'
    ) THEN
        ALTER TABLE neurondb_agent.audit_log
        ADD COLUMN org_id UUID REFERENCES neurondb_agent.organizations(id) ON DELETE SET NULL;
    END IF;
END $$;
CREATE INDEX IF NOT EXISTS idx_audit_log_org_id ON neurondb_agent.audit_log(org_id) WHERE org_id IS NOT NULL;

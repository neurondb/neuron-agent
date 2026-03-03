/*-------------------------------------------------------------------------
 *
 * main.go
 *    NeuronAgent HTTP server entry point
 *
 * Wires router, middleware, and all API handlers (agents, sessions, messages,
 * tools, LLM SQL). NeuronSQL and other features are loaded as modules.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/cmd/agent-server/main.go
 *
 *-------------------------------------------------------------------------
 */

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/neurondb/NeuronAgent/internal/agent"
	"github.com/neurondb/NeuronAgent/internal/api"
	"github.com/neurondb/NeuronAgent/internal/auth"
	"github.com/neurondb/NeuronAgent/internal/config"
	"github.com/neurondb/NeuronAgent/internal/core"
	"github.com/neurondb/NeuronAgent/internal/metrics"
	"github.com/neurondb/NeuronAgent/internal/modules/neuronsql"
	"github.com/neurondb/NeuronAgent/internal/tools"
	"github.com/neurondb/NeuronAgent/internal/workflow"
	"github.com/neurondb/NeuronAgent/pkg/llm_sql"
)

var (
	version   = "latest"
	buildDate string
	vcsRef    string
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}

	app, err := core.NewApp(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "App init failed: %v\n", err)
		os.Exit(1)
	}
	defer app.Close()

	/* Register enabled modules */
	if isModuleEnabled(cfg, "neuronsql") {
		if err := app.Registry().Register(neuronsql.New()); err != nil {
			fmt.Fprintf(os.Stderr, "Module registration failed: %v\n", err)
			os.Exit(1)
		}
	}
	if err := app.Registry().Resolve(); err != nil {
		fmt.Fprintf(os.Stderr, "Module resolve failed: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := app.Registry().InitAll(ctx, app); err != nil {
		fmt.Fprintf(os.Stderr, "Module init failed: %v\n", err)
		os.Exit(1)
	}

	/* Tool registry: base handlers + module-registered handlers */
	registry := tools.NewRegistry(app.Queries(), app.DB())
	if cfg.Tools.Timeout > 0 {
		registry.SetToolTimeout(cfg.Tools.Timeout)
	}
	for name, handler := range app.RegisteredToolHandlers() {
		registry.RegisterHandler(name, handler)
	}

	runtime := agent.NewRuntime(app.DB(), app.Queries(), registry, nil)
	handlers := api.NewHandlers(app.Queries(), runtime)

	/* Workflow engine with audit and tool registry */
	workflowEngine := workflow.NewEngine(app.Queries())
	workflowEngine.SetRuntime(runtime)
	workflowEngine.SetToolRegistry(registry)
	if cfg.Workflow.MaxDuration > 0 {
		workflowEngine.SetMaxDuration(cfg.Workflow.MaxDuration)
	}
	auditLogger := auth.NewAuditLogger(app.Queries())
	workflowEngine.SetAuditLogger(auditLogger)
	workflowHandlers := api.NewWorkflowHandlers(app.Queries(), workflowEngine)

	llmBaseURL := os.Getenv("LLM_SQL_BASE_URL")
	if llmBaseURL == "" {
		llmBaseURL = "http://localhost:8000"
	}
	llmAPIKey := os.Getenv("LLM_SQL_API_KEY")
	modelClient := llm_sql.NewModelClient(llmBaseURL, llmAPIKey)
	sqlLlmHandlers := api.NewSQLLLMHandlers(modelClient)

	keyManager := auth.NewAPIKeyManager(app.Queries())
	principalManager := auth.NewPrincipalManager(app.Queries())
	rateLimiter := auth.NewRateLimiter()

	router := buildRouter(cfg, handlers, sqlLlmHandlers, app, keyManager, principalManager, rateLimiter, registry, workflowHandlers)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	app.StartWorkers(ctx)
	if err := app.Registry().StartAll(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Module start failed: %v\n", err)
		os.Exit(1)
	}

	go func() {
		metrics.InfoWithContext(context.Background(), "NeuronAgent server starting", map[string]interface{}{
			"addr":    addr,
			"version": version,
		})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			metrics.ErrorWithContext(context.Background(), "Server error", err, map[string]interface{}{"addr": addr})
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.SetKeepAlivesEnabled(false)
	_ = auditLogger.Flush(shutdownCtx)
	if err := app.Registry().StopAll(shutdownCtx); err != nil {
		metrics.WarnWithContext(context.Background(), "Module stop error", map[string]interface{}{"error": err.Error()})
	}
	if err := srv.Shutdown(shutdownCtx); err != nil {
		metrics.WarnWithContext(context.Background(), "Server shutdown error", map[string]interface{}{"error": err.Error()})
	}
	registry.Cleanup()
	metrics.InfoWithContext(context.Background(), "NeuronAgent server stopped", map[string]interface{}{})
}

func isModuleEnabled(cfg *config.Config, name string) bool {
	if cfg.Modules == nil {
		return false
	}
	ent, ok := cfg.Modules[name]
	return ok && ent.Enabled
}

func loadConfig() (*config.Config, error) {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return config.LoadConfig(path)
	}
	cfg := config.DefaultConfig()
	if err := config.LoadFromEnv(cfg); err != nil {
		return nil, err
	}
	config.ApplyProfile(cfg)
	if err := config.ValidateConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func buildRouter(cfg *config.Config, h *api.Handlers, sqlLlm *api.SQLLLMHandlers, app *core.App, keyManager *auth.APIKeyManager, principalManager *auth.PrincipalManager, rateLimiter auth.RateLimiterInterface, registry *tools.Registry, workflowHandlers *api.WorkflowHandlers) http.Handler {
	r := mux.NewRouter()

	adminHandlers := api.NewAdminHandlers(cfg, app)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)
	r.Handle("/metrics", metrics.Handler()).Methods(http.MethodGet)

	r.Use(api.RequestIDMiddleware)
	r.Use(api.RequestTimeoutMiddleware(60 * time.Second))
	r.Use(api.AuthMiddleware(keyManager, principalManager, rateLimiter))
	r.Use(api.OrgMiddleware)
	r.Use(api.RejectUnknownFieldsMiddleware(cfg))
	r.Use(api.RequestBodyLimitMiddleware(10 << 20))
	r.Use(api.SecurityHeadersMiddleware)
	r.Use(api.CORSMiddleware(cfg))
	r.Use(api.LoggingMiddleware)

	v1 := r.PathPrefix("/api/v1").Subrouter()

	/* Agents */
	v1.HandleFunc("/blueprints", h.ListBlueprints).Methods(http.MethodGet)
	v1.HandleFunc("/agents/from-blueprint", h.CreateAgentFromBlueprint).Methods(http.MethodPost)
	v1.HandleFunc("/agents", h.CreateAgent).Methods(http.MethodPost)
	v1.HandleFunc("/agents", h.ListAgents).Methods(http.MethodGet)
	v1.HandleFunc("/agents/{id}", h.GetAgent).Methods(http.MethodGet)
	v1.HandleFunc("/agents/{id}", h.UpdateAgent).Methods(http.MethodPut)
	v1.HandleFunc("/agents/{id}", h.DeleteAgent).Methods(http.MethodDelete)
	v1.HandleFunc("/agents/{id}/clone", h.CloneAgent).Methods(http.MethodPost)
	v1.HandleFunc("/agents/{id}/plan", h.GeneratePlan).Methods(http.MethodPost)
	v1.HandleFunc("/agents/{id}/reflect", h.ReflectOnResponse).Methods(http.MethodPost)
	v1.HandleFunc("/agents/{id}/delegate", h.DelegateToAgent).Methods(http.MethodPost)
	v1.HandleFunc("/agents/{id}/metrics", h.GetAgentMetrics).Methods(http.MethodGet)
	v1.HandleFunc("/agents/{id}/costs", h.GetAgentCosts).Methods(http.MethodGet)
	v1.HandleFunc("/agents/{id}/sessions", h.ListSessions).Methods(http.MethodGet)
	v1.HandleFunc("/agents/{id}/memory", h.ListMemoryChunks).Methods(http.MethodGet)
	v1.HandleFunc("/agents/{id}/memory/search", h.SearchMemory).Methods(http.MethodPost)
	v1.HandleFunc("/agents/{id}/memory/summarize", h.SummarizeMemory).Methods(http.MethodPost)

	/* Sessions */
	v1.HandleFunc("/sessions", h.CreateSession).Methods(http.MethodPost)
	v1.HandleFunc("/sessions", h.ListSessions).Methods(http.MethodGet)
	v1.HandleFunc("/sessions/{id}", h.GetSession).Methods(http.MethodGet)
	v1.HandleFunc("/sessions/{id}", h.UpdateSession).Methods(http.MethodPut)
	v1.HandleFunc("/sessions/{id}", h.DeleteSession).Methods(http.MethodDelete)
	v1.HandleFunc("/sessions/{id}/messages", h.SendMessage).Methods(http.MethodPost)
	v1.HandleFunc("/sessions/{id}/messages", h.GetMessages).Methods(http.MethodGet)
	v1.HandleFunc("/sessions/{id}/messages/{msgId}", h.GetMessage).Methods(http.MethodGet)
	v1.HandleFunc("/sessions/{id}/messages/{msgId}", h.UpdateMessage).Methods(http.MethodPut)
	v1.HandleFunc("/sessions/{id}/messages/{msgId}", h.DeleteMessage).Methods(http.MethodDelete)
	v1.HandleFunc("/sessions/{id}/reflect", h.ReflectOnResponse).Methods(http.MethodPost)
	v1.HandleFunc("/sessions/{id}/feedback", h.SubmitFeedback).Methods(http.MethodPost)

	/* Tools */
	v1.HandleFunc("/tools", h.CreateTool).Methods(http.MethodPost)
	v1.HandleFunc("/tools", h.ListTools).Methods(http.MethodGet)
	v1.HandleFunc("/tools/{id}", h.GetTool).Methods(http.MethodGet)
	v1.HandleFunc("/tools/{id}", h.UpdateTool).Methods(http.MethodPut)
	v1.HandleFunc("/tools/{id}", h.DeleteTool).Methods(http.MethodDelete)
	v1.HandleFunc("/tools/{id}/analytics", h.GetToolAnalytics).Methods(http.MethodGet)

	/* Memory */
	v1.HandleFunc("/memory/{chunkId}", h.GetMemoryChunk).Methods(http.MethodGet)
	v1.HandleFunc("/memory/{chunkId}", h.DeleteMemoryChunk).Methods(http.MethodDelete)

	/* Budget */
	v1.HandleFunc("/agents/{id}/budget", h.GetBudget).Methods(http.MethodGet)
	v1.HandleFunc("/agents/{id}/budget", h.SetBudget).Methods(http.MethodPost)
	v1.HandleFunc("/agents/{id}/budget", h.UpdateBudget).Methods(http.MethodPut)

	/* Analytics */
	v1.HandleFunc("/analytics/overview", h.GetAnalyticsOverview).Methods(http.MethodGet)
	v1.HandleFunc("/analytics/retrieval-stats", h.GetRetrievalStats).Methods(http.MethodGet)
	v1.HandleFunc("/agents/{id}/memory/consolidate", h.ConsolidateMemory).Methods(http.MethodPost)

	/* Batch */
	v1.HandleFunc("/batch/agents", h.BatchCreateAgents).Methods(http.MethodPost)
	v1.HandleFunc("/batch/agents/delete", h.BatchDeleteAgents).Methods(http.MethodPost)
	v1.HandleFunc("/batch/messages/delete", h.BatchDeleteMessages).Methods(http.MethodPost)
	v1.HandleFunc("/batch/tools/delete", h.BatchDeleteTools).Methods(http.MethodPost)

	/* Human-in-the-loop */
	v1.HandleFunc("/approvals", h.ListApprovalRequests).Methods(http.MethodGet)
	v1.HandleFunc("/approvals/{id}", h.GetApprovalRequest).Methods(http.MethodGet)
	v1.HandleFunc("/approvals/{id}/approve", h.ApproveRequest).Methods(http.MethodPost)
	v1.HandleFunc("/approvals/{id}/reject", h.RejectRequest).Methods(http.MethodPost)
	v1.HandleFunc("/feedback", h.ListFeedback).Methods(http.MethodGet)
	v1.HandleFunc("/feedback/stats", h.GetFeedbackStats).Methods(http.MethodGet)

	/* Memory management */
	v1.HandleFunc("/agents/{id}/memory/corruption", h.CheckMemoryCorruption).Methods(http.MethodPost)
	v1.HandleFunc("/agents/{id}/memory/forget", h.ForgetMemories).Methods(http.MethodPost)
	v1.HandleFunc("/agents/{id}/memory/conflicts", h.ResolveMemoryConflicts).Methods(http.MethodPost)
	v1.HandleFunc("/agents/{id}/memory/quality", h.GetMemoryQuality).Methods(http.MethodGet)
	v1.HandleFunc("/agents/{id}/memory/feedback", h.SubmitMemoryFeedback).Methods(http.MethodPost)

	/* LLM SQL (existing) */
	v1.HandleFunc("/llm/sql/generate", sqlLlm.GenerateSQL).Methods(http.MethodPost)
	v1.HandleFunc("/llm/sql/explain", sqlLlm.ExplainSQL).Methods(http.MethodPost)
	v1.HandleFunc("/llm/sql/optimize", sqlLlm.OptimizeSQL).Methods(http.MethodPost)
	v1.HandleFunc("/llm/sql/debug", sqlLlm.DebugSQL).Methods(http.MethodPost)
	v1.HandleFunc("/llm/sql/translate", sqlLlm.TranslateSQL).Methods(http.MethodPost)
	v1.HandleFunc("/llm/sql/models", sqlLlm.ListModels).Methods(http.MethodGet)

	/* Admin: config dump and diagnostics (admin RBAC) */
	v1.HandleFunc("/admin/config", adminHandlers.GetAdminConfig).Methods(http.MethodGet)
	v1.HandleFunc("/admin/diagnostics", adminHandlers.GetAdminDiagnostics).Methods(http.MethodGet)
	v1.HandleFunc("/governance/costs", h.GetGovernanceCosts).Methods(http.MethodGet)
	v1.HandleFunc("/governance/tool-risk", h.GetGovernanceToolRisk).Methods(http.MethodGet)
	v1.HandleFunc("/governance/policy-blocks", h.GetGovernancePolicyBlocks).Methods(http.MethodGet)
	v1.HandleFunc("/governance/memory-growth", h.GetGovernanceMemoryGrowth).Methods(http.MethodGet)
	v1.HandleFunc("/governance/agent-performance", h.GetGovernanceAgentPerformance).Methods(http.MethodGet)

	/* Workflows */
	v1.HandleFunc("/workflows", workflowHandlers.CreateWorkflow).Methods(http.MethodPost)
	v1.HandleFunc("/workflows", workflowHandlers.ListWorkflows).Methods(http.MethodGet)
	v1.HandleFunc("/workflows/{id}", workflowHandlers.GetWorkflow).Methods(http.MethodGet)
	v1.HandleFunc("/workflows/{id}", workflowHandlers.UpdateWorkflow).Methods(http.MethodPut)
	v1.HandleFunc("/workflows/{id}", workflowHandlers.DeleteWorkflow).Methods(http.MethodDelete)
	v1.HandleFunc("/workflows/{workflow_id}/steps", workflowHandlers.CreateWorkflowStep).Methods(http.MethodPost)
	v1.HandleFunc("/workflows/{workflow_id}/steps", workflowHandlers.ListWorkflowSteps).Methods(http.MethodGet)
	v1.HandleFunc("/workflows/{workflow_id}/execute", workflowHandlers.ExecuteWorkflow).Methods(http.MethodPost)
	v1.HandleFunc("/workflows/{workflow_id}/executions/{execution_id}", workflowHandlers.GetWorkflowExecution).Methods(http.MethodGet)
	v1.HandleFunc("/workflows/{workflow_id}/executions", workflowHandlers.ListWorkflowExecutions).Methods(http.MethodGet)
	v1.HandleFunc("/workflows/{workflow_id}/schedule", workflowHandlers.CreateWorkflowSchedule).Methods(http.MethodPost)
	v1.HandleFunc("/workflows/{workflow_id}/schedule", workflowHandlers.GetWorkflowSchedule).Methods(http.MethodGet)
	v1.HandleFunc("/workflows/{workflow_id}/schedule", workflowHandlers.UpdateWorkflowSchedule).Methods(http.MethodPut)
	v1.HandleFunc("/workflows/{workflow_id}/schedule", workflowHandlers.DeleteWorkflowSchedule).Methods(http.MethodDelete)
	v1.HandleFunc("/workflows/schedules", workflowHandlers.ListWorkflowSchedules).Methods(http.MethodGet)

	/* Claw gateway: neuronsql.* tools only */
	clawHandlers := api.NewClawHandlers(registry)
	claw := r.PathPrefix("/claw/v1").Subrouter()
	claw.HandleFunc("/tools/list", clawHandlers.ListTools).Methods(http.MethodPost)
	claw.HandleFunc("/tools/run", clawHandlers.RunTool).Methods(http.MethodPost)
	claw.HandleFunc("/health", clawHandlers.Health).Methods(http.MethodGet)

	/* Module routes (e.g. NeuronSQL at /api/v1/neuronsql/*) */
	app.ApplyModuleRoutes(v1)

	return r
}

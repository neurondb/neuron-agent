/*-------------------------------------------------------------------------
 *
 * neuronsql.go
 *    NeuronSQL CLI: ingest docs, run eval
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/cli/cmd/neuronsql.go
 *
 *-------------------------------------------------------------------------
 */

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/neurondb/NeuronAgent/internal/neuronsql/eval"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/orchestrator"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/policy"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/provider"
	neuronsqltools "github.com/neurondb/NeuronAgent/internal/neuronsql/tools"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/retrieval"
	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
	"github.com/spf13/cobra"
)

var (
	neuronsqlDocsDirs []string
	neuronsqlIndexDir string
	neuronsqlEvalDSN  string
	neuronsqlEvalSuite string
	neuronsqlEvalOut  string
)

var neuronsqlCmd = &cobra.Command{
	Use:   "neuronsql",
	Short: "NeuronSQL: SQL/PLpgSQL copilot (ingest docs, eval)",
	Long:  "NeuronSQL commands for document ingestion and evaluation.",
}

var neuronsqlIngestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Ingest docs and build BM25 index",
	Long:  "Scan docs directories, chunk to 400-800 tokens, build BM25 index, persist under index dir.",
	RunE:  runNeuronSQLIngest,
}

var neuronsqlEvalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Run NeuronSQL eval suite and write report",
	Long:  "Run a named eval suite (saas_basic, northwind_like, analytics_timeseries) and write EvalReport JSON to file for CI.",
	RunE:  runNeuronSQLEval,
}

func init() {
	neuronsqlIngestCmd.Flags().StringSliceVar(&neuronsqlDocsDirs, "docs-dir", nil, "Directories to scan for markdown/text (repeat or comma-separated)")
	neuronsqlIngestCmd.Flags().StringVar(&neuronsqlIndexDir, "index-dir", "data/neuronsql/index", "Output directory for index files")
	neuronsqlCmd.AddCommand(neuronsqlIngestCmd)
	neuronsqlEvalCmd.Flags().StringVar(&neuronsqlEvalDSN, "dsn", "", "PostgreSQL DSN for eval (e.g. postgres://user:pass@localhost/db)")
	neuronsqlEvalCmd.Flags().StringVar(&neuronsqlEvalSuite, "suite", "saas_basic", "Suite name: saas_basic, northwind_like, analytics_timeseries")
	neuronsqlEvalCmd.Flags().StringVar(&neuronsqlEvalOut, "out", "eval_report.json", "Output path for report JSON")
	neuronsqlCmd.AddCommand(neuronsqlEvalCmd)
	rootCmd.AddCommand(neuronsqlCmd)
}

func runNeuronSQLIngest(cmd *cobra.Command, args []string) error {
	dirs := neuronsqlDocsDirs
	if len(dirs) == 0 {
		dirs = []string{"docs"}
	}
	var docs []neuronsql.Document
	for _, dir := range dirs {
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext != ".md" && ext != ".txt" && ext != ".sql" {
				return nil
			}
			body, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			docs = append(docs, neuronsql.Document{
				ID:      path,
				Source:  dir,
				Path:    path,
				Content: string(body),
			})
			return nil
		}); err != nil {
			return fmt.Errorf("walk %s: %w", dir, err)
		}
	}
	if len(docs) == 0 {
		fmt.Fprintln(os.Stderr, "No documents found; use --docs-dir to specify directories.")
		return nil
	}
	ret := retrieval.NewHybridRetriever(nil, nil, nil, nil)
	if err := ret.Index(context.Background(), docs); err != nil {
		return err
	}
	if err := os.MkdirAll(neuronsqlIndexDir, 0755); err != nil {
		return err
	}
	if err := ret.SaveIndex(neuronsqlIndexDir); err != nil {
		return fmt.Errorf("save index: %w", err)
	}
	fmt.Printf("Ingested %d docs, index saved to %s\n", len(docs), neuronsqlIndexDir)
	return nil
}

func runNeuronSQLEval(cmd *cobra.Command, args []string) error {
	if neuronsqlEvalDSN == "" {
		return fmt.Errorf("--dsn required for eval (e.g. postgres://user:pass@localhost/db)")
	}
	policyEngine := policy.NewPolicyEngineImpl(nil)
	factory := neuronsqltools.NewConnectionFactory(policyEngine, neuronsqltools.DefaultSafeConnectionConfig())
	pglangCfg := provider.DefaultPGLangConfig()
	if e := os.Getenv("PGLANG_ENDPOINT"); e != "" {
		pglangCfg.Endpoint = e
	}
	llmProvider := provider.NewPGLangProvider(pglangCfg)
	orch := orchestrator.NewOrchestrator(llmProvider, policyEngine, factory, nil)
	runner := &eval.Runner{Orchestrator: orch, DSN: neuronsqlEvalDSN}
	report, err := runner.RunSuite(context.Background(), neuronsqlEvalSuite)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(neuronsqlEvalOut, b, 0644); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	fmt.Printf("Eval suite %s: pass_rate=%.2f, report written to %s\n", neuronsqlEvalSuite, report.PassRate, neuronsqlEvalOut)
	return nil
}

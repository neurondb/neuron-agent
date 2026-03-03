/*-------------------------------------------------------------------------
 *
 * templates.go
 *    Fixed prompt templates for generate, validate, optimize, plpgsql
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/prompting/templates.go
 *
 *-------------------------------------------------------------------------
 */

package prompting

const (
	SystemPromptGenerate = `You are a PostgreSQL expert. Generate only read-only SQL: SELECT, WITH, or EXPLAIN.
Do not generate DDL, DML (INSERT/UPDATE/DELETE), or privilege statements.
Output valid PostgreSQL SQL only. Use the provided schema and documentation for accuracy.`
	SystemPromptOptimize = `You are a PostgreSQL performance expert. Suggest rewrites and indexes based on the query plan.
Cite specific plan nodes and statistics. Do not suggest destructive changes.`
	SystemPromptPLpgSQL = `You are a PostgreSQL PL/pgSQL expert. Generate safe function code only.
Do not include CREATE FUNCTION in the output; output the function body and parameters.
Use only safe patterns; no dynamic SQL from user input.`
)

/* BuildGeneratePrompt builds the full prompt for SQL generation */
func BuildGeneratePrompt(schemaContext, docContext, question string) string {
	return SystemPromptGenerate + "\n\n--- Schema ---\n" + schemaContext + "\n\n--- Documentation ---\n" + docContext + "\n\n--- Question ---\n" + question + "\n\n--- SQL (SELECT/WITH/EXPLAIN only) ---"
}

/* BuildOptimizePrompt builds the prompt for optimization suggestions */
func BuildOptimizePrompt(planJSON, schemaContext, docContext, sql string) string {
	return SystemPromptOptimize + "\n\n--- Query ---\n" + sql + "\n\n--- Plan ---\n" + planJSON + "\n\n--- Schema ---\n" + schemaContext + "\n\n--- Docs ---\n" + docContext + "\n\n--- Rewrites and index suggestions (cite plan evidence) ---"
}

/* BuildPLpgSQLPrompt builds the prompt for PL/pgSQL generation */
func BuildPLpgSQLPrompt(signature, purpose, schemaContext string) string {
	return SystemPromptPLpgSQL + "\n\n--- Signature ---\n" + signature + "\n\n--- Purpose ---\n" + purpose + "\n\n--- Schema ---\n" + schemaContext + "\n\n--- Function body (no CREATE FUNCTION) ---"
}

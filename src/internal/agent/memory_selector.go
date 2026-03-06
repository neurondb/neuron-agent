/*-------------------------------------------------------------------------
 * memory_selector.go
 *    Tier-aware memory selection with token budgeting for agent context.
 *-------------------------------------------------------------------------*/

package agent

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/pkg/neurondb"
)

const (
	defaultEmbedModel = "all-MiniLM-L6-v2"
	approxTokensPerChar = 4
	recentBudgetFraction = 0.2
)

// MemoryItem is a single selected memory chunk with metadata for context building.
type MemoryItem struct {
	Content    string
	TokenCount int
	Score      float64
	Tier       string // "working", "episodic", "semantic"
	SourceID   uuid.UUID
	CreatedAt  time.Time
}

// MemorySelector selects memory across STM/MTM/LPM within a token budget.
type MemorySelector interface {
	Select(ctx context.Context, agentID, sessionID uuid.UUID, query string, budgetTokens int) ([]MemoryItem, error)
}

// MemorySelectorConfig holds weights for scoring (relevance, recency, importance).
type MemorySelectorConfig struct {
	Alpha float64 // relevance
	Beta  float64 // recency
	Gamma float64 // importance
	STMLimit int
	MTMTopK int
	LPMTopK int
	LPMMinScore float64
}

// DefaultMemorySelectorConfig returns plan-default weights.
func DefaultMemorySelectorConfig() MemorySelectorConfig {
	return MemorySelectorConfig{
		Alpha: 0.6, Beta: 0.3, Gamma: 0.1,
		STMLimit: 20, MTMTopK: 10, LPMTopK: 5, LPMMinScore: 0.7,
	}
}

// DBMemorySelector implements MemorySelector using DB and embedding client.
type DBMemorySelector struct {
	queries *db.Queries
	embed   *neurondb.EmbeddingClient
	config  MemorySelectorConfig
}

// NewDBMemorySelector creates a MemorySelector that uses the given queries and embed client.
func NewDBMemorySelector(queries *db.Queries, embed *neurondb.EmbeddingClient, config MemorySelectorConfig) *DBMemorySelector {
	if config.STMLimit <= 0 {
		config.STMLimit = 20
	}
	if config.MTMTopK <= 0 {
		config.MTMTopK = 10
	}
	if config.LPMTopK <= 0 {
		config.LPMTopK = 5
	}
	return &DBMemorySelector{queries: queries, embed: embed, config: config}
}

// Select returns memory items from working (STM), episodic (MTM), and semantic (LPM) tiers within budget.
func (s *DBMemorySelector) Select(ctx context.Context, agentID, sessionID uuid.UUID, query string, budgetTokens int) ([]MemoryItem, error) {
	var candidates []MemoryItem

	// 1. Working memory: STM by session (no vector search)
	stm, err := s.selectSTM(ctx, agentID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("memory selector STM: %w", err)
	}
	candidates = append(candidates, stm...)

	// 2. Episodic/semantic: embed query and search MTM, LPM
	queryEmbed, err := s.embed.Embed(ctx, query, defaultEmbedModel)
	if err != nil {
		return nil, fmt.Errorf("memory selector embed: %w", err)
	}

	mtm, err := s.selectMTM(ctx, agentID, queryEmbed)
	if err != nil {
		return nil, fmt.Errorf("memory selector MTM: %w", err)
	}
	candidates = append(candidates, mtm...)

	lpm, err := s.selectLPM(ctx, agentID, queryEmbed)
	if err != nil {
		return nil, fmt.Errorf("memory selector LPM: %w", err)
	}
	candidates = append(candidates, lpm...)

	// 3. Deduplicate by content hash (keep highest score per hash)
	seen := make(map[string]int)
	for i := range candidates {
		h := contentHash(candidates[i].Content)
		if j, ok := seen[h]; !ok || candidates[i].Score > candidates[j].Score {
			seen[h] = i
		}
	}
	var unique []MemoryItem
	for _, i := range seen {
		unique = append(unique, candidates[i])
	}

	// 4. Sort by score descending
	sort.Slice(unique, func(i, j int) bool { return unique[i].Score > unique[j].Score })

	// 5. Greedy pack: reserve recentBudgetFraction for working, fill rest by score
	recentBudget := int(float64(budgetTokens) * recentBudgetFraction)
	relevanceBudget := budgetTokens - recentBudget
	var selected []MemoryItem
	usedRecent, usedRelevance := 0, 0
	for _, c := range unique {
		if c.Tier == "working" && usedRecent+c.TokenCount <= recentBudget {
			selected = append(selected, c)
			usedRecent += c.TokenCount
		}
	}
	for _, c := range unique {
		if c.Tier == "working" {
			continue
		}
		if usedRelevance+c.TokenCount <= relevanceBudget {
			selected = append(selected, c)
			usedRelevance += c.TokenCount
		}
	}
	return selected, nil
}

func contentHash(s string) string {
	if len(s) > 500 {
		s = s[:500]
	}
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func estimateTokens(content string) int {
	n := len(content) / approxTokensPerChar
	if n < 1 && len(content) > 0 {
		return 1
	}
	return n
}

func (s *DBMemorySelector) selectSTM(ctx context.Context, agentID, sessionID uuid.UUID) ([]MemoryItem, error) {
	q := `SELECT id, content, importance_score, created_at
		FROM neurondb_agent.memory_stm
		WHERE agent_id = $1 AND session_id = $2 AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY created_at DESC
		LIMIT $3`
	var rows []struct {
		ID              uuid.UUID `db:"id"`
		Content         string    `db:"content"`
		ImportanceScore float64   `db:"importance_score"`
		CreatedAt       time.Time `db:"created_at"`
	}
	err := s.queries.DB.SelectContext(ctx, &rows, q, agentID, sessionID, s.config.STMLimit)
	if err != nil {
		return nil, err
	}
	out := make([]MemoryItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, MemoryItem{
			Content:    r.Content,
			TokenCount: estimateTokens(r.Content),
			Score:      1.0,
			Tier:       "working",
			SourceID:   r.ID,
			CreatedAt:  r.CreatedAt,
		})
	}
	return out, nil
}

func (s *DBMemorySelector) selectMTM(ctx context.Context, agentID uuid.UUID, queryEmbed neurondb.Vector) ([]MemoryItem, error) {
	return s.vectorSelectTier(ctx, "neurondb_agent.memory_mtm", "episodic", agentID, queryEmbed, s.config.MTMTopK, 0, true)
}

func (s *DBMemorySelector) selectLPM(ctx context.Context, agentID uuid.UUID, queryEmbed neurondb.Vector) ([]MemoryItem, error) {
	return s.vectorSelectTier(ctx, "neurondb_agent.memory_lpm", "semantic", agentID, queryEmbed, s.config.LPMTopK, s.config.LPMMinScore, false)
}

func (s *DBMemorySelector) vectorSelectTier(ctx context.Context, table, tier string, agentID uuid.UUID, queryEmbed neurondb.Vector, limit int, minScore float64, hasExpires bool) ([]MemoryItem, error) {
	vecStr := formatVectorForPostgres(queryEmbed)
	whereExpires := ""
	if hasExpires {
		whereExpires = " AND (expires_at IS NULL OR expires_at > NOW())"
	}
	q := fmt.Sprintf(`SELECT id, content, importance_score, created_at,
		1 - (embedding <=> $1::vector) AS similarity
		FROM %s
		WHERE agent_id = $2 AND embedding IS NOT NULL %s
		ORDER BY embedding <=> $1::vector
		LIMIT $3`, table, whereExpires)
	var rows []struct {
		ID              uuid.UUID `db:"id"`
		Content         string    `db:"content"`
		ImportanceScore float64   `db:"importance_score"`
		CreatedAt       time.Time `db:"created_at"`
		Similarity      float64   `db:"similarity"`
	}
	err := s.queries.DB.SelectContext(ctx, &rows, q, vecStr, agentID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]MemoryItem, 0, len(rows))
	for _, r := range rows {
		if minScore > 0 && r.Similarity < minScore {
			continue
		}
		recency := 1.0
		if !r.CreatedAt.IsZero() {
			age := time.Since(r.CreatedAt)
			if age > 7*24*time.Hour {
				recency = 0.3
			} else if age > 24*time.Hour {
				recency = 0.6
			}
		}
		score := s.config.Alpha*r.Similarity + s.config.Beta*recency + s.config.Gamma*r.ImportanceScore
		out = append(out, MemoryItem{
			Content:    r.Content,
			TokenCount: estimateTokens(r.Content),
			Score:      score,
			Tier:       tier,
			SourceID:   r.ID,
			CreatedAt:  r.CreatedAt,
		})
	}
	return out, nil
}

// formatVectorForPostgres mirrors pkg/neurondb format for vector param (avoid importing internal detail).
func formatVectorForPostgres(vec []float32) string {
	if len(vec) == 0 {
		return "[]"
	}
	b := make([]byte, 0, len(vec)*8+2)
	b = append(b, '[')
	for i, v := range vec {
		if i > 0 {
			b = append(b, ',')
		}
		b = append(b, fmt.Sprintf("%g", v)...)
	}
	b = append(b, ']')
	return string(b)
}

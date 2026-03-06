/*-------------------------------------------------------------------------
 * indexer.go
 *    Repo analysis mode: indexing, chunking, summarization, dependency graph.
 *-------------------------------------------------------------------------*/

package repomode

import (
	"context"

	"github.com/google/uuid"
)

// ChunkType is the type of code chunk.
type ChunkType string

const (
	ChunkTypeFunction ChunkType = "function"
	ChunkTypeClass    ChunkType = "class"
	ChunkTypeModule   ChunkType = "module"
	ChunkTypeFile     ChunkType = "file"
	ChunkTypeSection  ChunkType = "section"
)

// RepoChunk is a single indexed chunk (file or sub-unit).
type RepoChunk struct {
	ID          uuid.UUID
	AgentID     uuid.UUID
	RepoPath    string
	FilePath    string
	FileHash    string
	Language    string
	ChunkType   ChunkType
	ChunkName   string
	ChunkContent string
	Summary     string
	TokenCount  int
	LineStart   int
	LineEnd     int
	Metadata    map[string]interface{}
}

// Dependency is an edge in the repo dependency graph.
type Dependency struct {
	SourceFile string
	TargetFile string
	Type       string
}

// Indexer indexes a repository for analysis (chunking, embedding, dependency graph).
type Indexer interface {
	Index(ctx context.Context, agentID uuid.UUID, repoPath string) error
	GetChunks(ctx context.Context, agentID uuid.UUID, repoPath, filePath string) ([]RepoChunk, error)
	GetDependencies(ctx context.Context, agentID uuid.UUID, repoPath string) ([]Dependency, error)
}

// DefaultIndexer is a stub that performs no indexing (implement full logic per plan).
type DefaultIndexer struct{}

// NewDefaultIndexer returns a stub indexer.
func NewDefaultIndexer() *DefaultIndexer {
	return &DefaultIndexer{}
}

// Index indexes the repo (stub: no-op).
func (i *DefaultIndexer) Index(ctx context.Context, agentID uuid.UUID, repoPath string) error {
	_ = ctx
	_ = agentID
	_ = repoPath
	return nil
}

// GetChunks returns chunks for the given file (stub: empty).
func (i *DefaultIndexer) GetChunks(ctx context.Context, agentID uuid.UUID, repoPath, filePath string) ([]RepoChunk, error) {
	_ = ctx
	_ = agentID
	_ = repoPath
	_ = filePath
	return nil, nil
}

// GetDependencies returns the dependency graph (stub: empty).
func (i *DefaultIndexer) GetDependencies(ctx context.Context, agentID uuid.UUID, repoPath string) ([]Dependency, error) {
	_ = ctx
	_ = agentID
	_ = repoPath
	return nil, nil
}

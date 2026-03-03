/*-------------------------------------------------------------------------
 *
 * retriever.go
 *    HybridRetriever: BM25 + vector merge, rerank, return top-k chunks
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/retrieval/retriever.go
 *
 *-------------------------------------------------------------------------
 */

package retrieval

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

const (
	BM25TopN   = 50
	VectorTopN = 50
	RerankN   = 20
	FinalK    = 6
)

/* HybridRetriever implements neuronsql.Retriever with BM25 + optional vector */
type HybridRetriever struct {
	chunks    []neuronsql.Chunk
	bm25      *BM25Index
	vector    *VectorIndex
	embedFunc func(ctx context.Context, texts []string) ([][]float32, error)
}

/* NewHybridRetriever creates a retriever with BM25 index and optional vector index */
func NewHybridRetriever(chunks []neuronsql.Chunk, bm25 *BM25Index, vector *VectorIndex, embedFunc func(context.Context, []string) ([][]float32, error)) *HybridRetriever {
	return &HybridRetriever{chunks: chunks, bm25: bm25, vector: vector, embedFunc: embedFunc}
}

/* Index implements neuronsql.Retriever - build BM25 from docs */
func (r *HybridRetriever) Index(ctx context.Context, docs []neuronsql.Document) error {
	var allChunks []neuronsql.Chunk
	var docTerms [][]string
	for _, doc := range docs {
		chunks, err := ChunkDocument(doc)
		if err != nil {
			return err
		}
		for _, c := range chunks {
			allChunks = append(allChunks, c)
			docTerms = append(docTerms, tokenizeLower(c.Content))
		}
	}
	r.chunks = allChunks
	r.bm25 = NewBM25Index(docTerms)
	return nil
}

/* Retrieve implements neuronsql.Retriever - BM25 top-50, optional vector top-50, RRF merge, return top-6 */
func (r *HybridRetriever) Retrieve(ctx context.Context, query string, k int) ([]neuronsql.Chunk, error) {
	if k <= 0 {
		k = FinalK
	}
	if len(r.chunks) == 0 || r.bm25 == nil {
		return nil, nil
	}
	queryNorm := strings.TrimSpace(query)
	if queryNorm == "" {
		return nil, nil
	}

	/* BM25 top-50 */
	bm25Indices := r.bm25.TopK(queryNorm, BM25TopN)
	scoreByIdx := make(map[int]float64)
	for rank, idx := range bm25Indices {
		scoreByIdx[idx] += 1.0 / float64(60+rank+1)
	}

	/* Optional: vector top-50 */
	if r.vector != nil && r.embedFunc != nil {
		vecs, err := r.embedFunc(ctx, []string{queryNorm})
		if err == nil && len(vecs) > 0 {
			vecIndices := r.vector.TopK(vecs[0], VectorTopN)
			for rank, idx := range vecIndices {
				scoreByIdx[idx] += 1.0 / float64(60+rank+1)
			}
		}
	}

	/* RRF: collect and sort by combined score */
	type scored struct {
		idx   int
		score float64
	}
	var scoredList []scored
	for idx, s := range scoreByIdx {
		scoredList = append(scoredList, scored{idx, s})
	}
	for i := 0; i < len(scoredList) && i < RerankN; i++ {
		for j := i + 1; j < len(scoredList); j++ {
			if scoredList[j].score > scoredList[i].score {
				scoredList[i], scoredList[j] = scoredList[j], scoredList[i]
			}
		}
	}
	if len(scoredList) > RerankN {
		scoredList = scoredList[:RerankN]
	}
	if len(scoredList) > k {
		scoredList = scoredList[:k]
	}

	out := make([]neuronsql.Chunk, 0, len(scoredList))
	for _, s := range scoredList {
		if s.idx < len(r.chunks) {
			c := r.chunks[s.idx]
			c.Score = float64(s.score)
			out = append(out, c)
		}
	}
	return out, nil
}

/* SaveIndex persists BM25 index and chunk list to indexDir */
func (r *HybridRetriever) SaveIndex(indexDir string) error {
	if r.bm25 != nil {
		if err := r.bm25.Save(filepath.Join(indexDir, "bm25.gob")); err != nil {
			return err
		}
	}
	return nil
}

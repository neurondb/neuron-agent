/*-------------------------------------------------------------------------
 *
 * vector.go
 *    Vector index with cosine similarity (pure Go)
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/retrieval/vector.go
 *
 *-------------------------------------------------------------------------
 */

package retrieval

import (
	"encoding/gob"
	"math"
	"os"
)

/* VectorIndex stores chunk IDs and vectors for cosine similarity search */
type VectorIndex struct {
	ChunkIDs []string
	Vectors  [][]float32
}

/* NewVectorIndex creates an index from chunk IDs and vectors */
func NewVectorIndex(chunkIDs []string, vectors [][]float32) *VectorIndex {
	return &VectorIndex{ChunkIDs: chunkIDs, Vectors: vectors}
}

/* CosineSimilarity returns cosine similarity between a and b */
func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return float32(dot / (math.Sqrt(normA) * math.Sqrt(normB)))
}

/* TopK returns indices of top-k chunks by cosine similarity to query vector */
func (v *VectorIndex) TopK(query []float32, k int) []int {
	if len(v.Vectors) == 0 || len(query) == 0 {
		return nil
	}
	scores := make([]float32, len(v.Vectors))
	for i, vec := range v.Vectors {
		scores[i] = CosineSimilarity(query, vec)
	}
	indices := make([]int, len(v.Vectors))
	for i := range indices {
		indices[i] = i
	}
	for i := 0; i < k && i < len(indices); i++ {
		for j := i + 1; j < len(indices); j++ {
			if scores[indices[j]] > scores[indices[i]] {
				indices[i], indices[j] = indices[j], indices[i]
			}
		}
	}
	if k > len(indices) {
		k = len(indices)
	}
	return indices[:k]
}

/* Save persists the vector index to file */
func (v *VectorIndex) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(v)
}

/* LoadVectorIndex loads from file */
func LoadVectorIndex(path string) (*VectorIndex, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var v VectorIndex
	if err := gob.NewDecoder(f).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

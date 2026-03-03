/*-------------------------------------------------------------------------
 *
 * bm25.go
 *    Okapi BM25 scoring (pure Go)
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/retrieval/bm25.go
 *
 *-------------------------------------------------------------------------
 */

package retrieval

import (
	"encoding/gob"
	"os"
	"strings"
	"unicode"
)

const (
	BM25K1 = 1.5
	BM25B  = 0.75
)

/* BM25Index holds document terms and stats for BM25 scoring */
type BM25Index struct {
	Docs     [][]string
	DocLen   []int
	AvgDocLen float64
	IDF      map[string]float64
	N        int
}

/* NewBM25Index builds an index from tokenized documents */
func NewBM25Index(docs [][]string) *BM25Index {
	n := len(docs)
	if n == 0 {
		return &BM25Index{N: 0, IDF: make(map[string]float64)}
	}
	totalLen := 0
	df := make(map[string]int)
	for _, doc := range docs {
		seen := make(map[string]bool)
		totalLen += len(doc)
		for _, t := range doc {
			if !seen[t] {
				seen[t] = true
				df[t]++
			}
		}
	}
	avgLen := float64(totalLen) / float64(n)
	idf := make(map[string]float64)
	for term, count := range df {
		// IDF = log((N+1)/(df+0.5)+1)
		idf[term] = 1.0
		if count > 0 {
			idf[term] = float64(n+1) / (float64(count) + 0.5)
		}
	}
	docLen := make([]int, n)
	for i, doc := range docs {
		docLen[i] = len(doc)
	}
	return &BM25Index{
		Docs:      docs,
		DocLen:    docLen,
		AvgDocLen: avgLen,
		IDF:       idf,
		N:         n,
	}
}


/* Score computes BM25 score of query terms against doc at index i */
func (b *BM25Index) Score(i int, queryTerms []string) float64 {
	if i < 0 || i >= b.N || b.AvgDocLen == 0 {
		return 0
	}
	var score float64
	tf := make(map[string]int)
	for _, t := range b.Docs[i] {
		tf[t]++
	}
	for _, q := range queryTerms {
		idf, ok := b.IDF[q]
		if !ok {
			continue
		}
		f := float64(tf[q])
		score += idf * (f * (BM25K1 + 1)) / (f + BM25K1*(1-BM25B+BM25B*float64(b.DocLen[i])/b.AvgDocLen))
	}
	return score
}

/* TopK returns indices of top-k documents by BM25 score */
func (b *BM25Index) TopK(query string, k int) []int {
	terms := tokenizeLower(query)
	if len(terms) == 0 {
		return nil
	}
	scores := make([]float64, b.N)
	for i := 0; i < b.N; i++ {
		scores[i] = b.Score(i, terms)
	}
	return topKIndices(scores, k)
}

func tokenizeLower(s string) []string {
	var out []string
	for _, w := range strings.FieldsFunc(s, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) }) {
		w = strings.ToLower(strings.TrimSpace(w))
		if w != "" {
			out = append(out, w)
		}
	}
	return out
}

func topKIndices(scores []float64, k int) []int {
	if k <= 0 || len(scores) == 0 {
		return nil
	}
	if k >= len(scores) {
		k = len(scores)
	}
	idx := make([]int, len(scores))
	for i := range idx {
		idx[i] = i
	}
	for i := 0; i < k; i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[idx[j]] > scores[idx[i]] {
				idx[i], idx[j] = idx[j], idx[i]
			}
		}
	}
	return idx[:k]
}

/* Save persists the index to a file */
func (b *BM25Index) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(b)
}

/* LoadBM25Index loads an index from file */
func LoadBM25Index(path string) (*BM25Index, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var b BM25Index
	if err := gob.NewDecoder(f).Decode(&b); err != nil {
		return nil, err
	}
	return &b, nil
}

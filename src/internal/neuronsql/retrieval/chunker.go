/*-------------------------------------------------------------------------
 *
 * chunker.go
 *    Document chunking (400-800 token approximation)
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/retrieval/chunker.go
 *
 *-------------------------------------------------------------------------
 */

package retrieval

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"

	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

const (
	MinChunkTokens = 400
	MaxChunkTokens = 800
	// Approximate tokens as words (space-separated); ~1.3 tokens per word for English
	approxTokensPerWord = 1
)

/* ChunkDocument splits a document into 400-800 token (approximate) chunks with metadata */
func ChunkDocument(doc neuronsql.Document) ([]neuronsql.Chunk, error) {
	content := strings.TrimSpace(doc.Content)
	if content == "" {
		return nil, nil
	}
	words := tokenize(content)
	if len(words) == 0 {
		return nil, nil
	}
	var chunks []neuronsql.Chunk
	start := 0
	for start < len(words) {
		end := start + MaxChunkTokens
		if end > len(words) {
			end = len(words)
		}
		if end-start < MinChunkTokens && end < len(words) {
			end = start + MinChunkTokens
			if end > len(words) {
				end = len(words)
			}
		}
		chunkWords := words[start:end]
		chunkContent := strings.Join(chunkWords, " ")
		h := sha256.Sum256([]byte(chunkContent))
		chunkID := doc.ID + "_" + hex.EncodeToString(h[:8])
		if len(chunkID) > 64 {
			chunkID = chunkID[:64]
		}
		chunks = append(chunks, neuronsql.Chunk{
			ID:      chunkID,
			Content: chunkContent,
			Score:   0,
			Source:  doc.Source,
			Path:    doc.Path,
			Section: doc.Section,
		})
		start = end
	}
	return chunks, nil
}

func tokenize(s string) []string {
	var words []string
	var b strings.Builder
	for _, r := range s {
		if unicode.IsSpace(r) {
			if b.Len() > 0 {
				words = append(words, b.String())
				b.Reset()
			}
			continue
		}
		b.WriteRune(r)
	}
	if b.Len() > 0 {
		words = append(words, b.String())
	}
	return words
}

/* ChunkID generates a short ID from source, path, and content hash */
func ChunkID(source, path, content string) string {
	h := sha256.Sum256([]byte(source + path + content))
	return hex.EncodeToString(h[:])[:16]
}

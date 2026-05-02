/*-------------------------------------------------------------------------
 *
 * openapi_embed.go
 *    Embedded OpenAPI spec for /docs (sync from src/openapi via Makefile copy).
 *
 *-------------------------------------------------------------------------
 */

package api

import _ "embed"

//go:embed specdata/openapi.yaml
var embeddedOpenAPIYAML []byte

/* EmbeddedOpenAPIYAML returns the embedded openapi.yaml bytes (may be empty if missing). */
func EmbeddedOpenAPIYAML() []byte {
	return embeddedOpenAPIYAML
}

/*-------------------------------------------------------------------------
 *
 * docs_handlers.go
 *    Serves OpenAPI YAML and a minimal Redoc UI at /docs.
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"net/http"
)

const docsHTML = `<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"/><title>NeuronAgent API</title>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet"/>
<style>body{margin:0;padding:0}</style></head><body>
<redoc spec-url="/docs/openapi.yaml"></redoc>
<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
</body></html>`

/* DocsIndex serves Redoc shell pointing at OpenAPI YAML */
func DocsIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(docsHTML))
}

/* DocsOpenAPISpec serves application/yaml OpenAPI document */
func DocsOpenAPISpec(w http.ResponseWriter, _ *http.Request) {
	b := EmbeddedOpenAPIYAML()
	if len(b) == 0 {
		http.Error(w, "OpenAPI spec not embedded", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

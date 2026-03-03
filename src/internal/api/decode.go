/*-------------------------------------------------------------------------
 *
 * decode.go
 *    Request body JSON decoding with optional strict schema (reject unknown fields).
 *
 * When RejectUnknownFieldsMiddleware has set the context flag, DecodeJSON
 * uses json.DisallowUnknownFields() so that request bodies with unknown
 * fields are rejected (production hardening).
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/api/decode.go
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"encoding/json"
	"net/http"
)

/* DecodeJSON decodes the request body into v. When RejectUnknownFieldsMiddleware
 * has set the context flag (e.g. Config.RejectUnknownFields true), unknown JSON
 * fields cause a decoding error. */
func DecodeJSON(r *http.Request, v interface{}) error {
	dec := json.NewDecoder(r.Body)
	if reject, _ := r.Context().Value(rejectUnknownFieldsKey).(bool); reject {
		dec.DisallowUnknownFields()
	}
	return dec.Decode(v)
}

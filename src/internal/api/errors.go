/*-------------------------------------------------------------------------
 *
 * errors.go
 *    API error handling and error types for NeuronAgent
 *
 * Provides structured error types and error response formatting for
 * the NeuronAgent HTTP API with context and metadata support.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/api/errors.go
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Code         int
	Message      string
	Err          error
	RequestID    string
	Endpoint     string
	Method       string
	ResourceType string
	ResourceID   string
	Details      map[string]interface{}
}

func (e *APIError) Error() string {
	parts := []string{e.Message}

	if e.Endpoint != "" {
		parts = append(parts, fmt.Sprintf("endpoint='%s'", e.Endpoint))
	}
	if e.Method != "" {
		parts = append(parts, fmt.Sprintf("method='%s'", e.Method))
	}
	if e.RequestID != "" {
		parts = append(parts, fmt.Sprintf("request_id='%s'", e.RequestID))
	}
	if e.ResourceType != "" {
		part := fmt.Sprintf("resource_type='%s'", e.ResourceType)
		if e.ResourceID != "" {
			part += fmt.Sprintf(", resource_id='%s'", e.ResourceID)
		}
		parts = append(parts, part)
	}

	if e.Err != nil {
		parts = append(parts, fmt.Sprintf("error=%v", e.Err))
	}

	/* Build error message from parts */
	if len(parts) > 1 {
		return fmt.Sprintf("%s: %v", parts[0], parts[1:])
	}
	return parts[0]
}

func NewError(code int, message string, err error) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: make(map[string]interface{}),
	}
}

func NewErrorWithRequestID(code int, message string, err error, requestID string) *APIError {
	return &APIError{
		Code:      code,
		Message:   message,
		Err:       err,
		RequestID: requestID,
		Details:   make(map[string]interface{}),
	}
}

func NewErrorWithContext(code int, message string, err error, requestID, endpoint, method, resourceType, resourceID string, details map[string]interface{}) *APIError {
	if details == nil {
		details = make(map[string]interface{})
	}
	return &APIError{
		Code:         code,
		Message:      message,
		Err:          err,
		RequestID:    requestID,
		Endpoint:     endpoint,
		Method:       method,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      details,
	}
}

var (
	ErrNotFound     = NewError(http.StatusNotFound, "resource not found", nil)
	ErrBadRequest   = NewError(http.StatusBadRequest, "bad request", nil)
	ErrUnauthorized = NewError(http.StatusUnauthorized, "unauthorized", nil)
	ErrInternal     = NewError(http.StatusInternalServerError, "internal server error", nil)
)

/* DefaultDocsURL is the base URL for error documentation */
const DefaultDocsURL = "https://docs.neurondb.ai/errors"

/* ErrorCodeFromHTTP maps HTTP status to a stable error_code string */
func ErrorCodeFromHTTP(code int) string {
	switch code {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusConflict:
		return "conflict"
	case http.StatusTooManyRequests:
		return "rate_limit_exceeded"
	case http.StatusInternalServerError:
		return "internal_error"
	case http.StatusServiceUnavailable:
		return "service_unavailable"
	case http.StatusGatewayTimeout:
		return "gateway_timeout"
	default:
		if code >= 500 {
			return "server_error"
		}
		if code >= 400 {
			return "client_error"
		}
		return "unknown"
	}
}

/* RetryableFromHTTP returns true for status codes that clients may retry */
func RetryableFromHTTP(code int) bool {
	return code == http.StatusTooManyRequests ||
		code == http.StatusServiceUnavailable ||
		code == http.StatusGatewayTimeout ||
		(code >= 500 && code < 600)
}

/* WrapError wraps an error with request ID */
func WrapError(err *APIError, requestID string) *APIError {
	if err == nil {
		return nil
	}
	return NewErrorWithRequestID(err.Code, err.Message, err.Err, requestID)
}

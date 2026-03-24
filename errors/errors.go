// Package errors defines all custom error types for the Yapily SDK.
package errors

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned from the API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	TraceID    string
}

func (e *APIError) Error() string {
	if e.TraceID != "" {
		return fmt.Sprintf("api error [%d] %s: %s (trace: %s)", e.StatusCode, e.Code, e.Message, e.TraceID)
	}
	return fmt.Sprintf("api error [%d] %s: %s", e.StatusCode, e.Code, e.Message)
}

// IsNotFound returns true if the error is a 404.
func IsNotFound(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == http.StatusNotFound
	}
	return false
}

// IsUnauthorized returns true if the error is a 401.
func IsUnauthorized(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsRateLimited returns true if the error is a 429.
func IsRateLimited(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// ValidationError represents a client-side validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s' - %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// AuthError represents an authentication or token error.
type AuthError struct {
	Message string
	Cause   error
}

func (e *AuthError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("auth error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("auth error: %s", e.Message)
}

func (e *AuthError) Unwrap() error {
	return e.Cause
}

// NewAuthError creates a new AuthError.
func NewAuthError(message string, cause error) *AuthError {
	return &AuthError{Message: message, Cause: cause}
}

// RetryableError wraps an error to indicate it is safe to retry.
type RetryableError struct {
	Cause error
}

func (e *RetryableError) Error() string {
	return fmt.Sprintf("retryable error: %v", e.Cause)
}

func (e *RetryableError) Unwrap() error {
	return e.Cause
}

// IsRetryable returns true if the error is retryable.
func IsRetryable(err error) bool {
	if _, ok := err.(*RetryableError); ok {
		return true
	}
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == http.StatusTooManyRequests ||
			e.StatusCode == http.StatusServiceUnavailable ||
			e.StatusCode == http.StatusGatewayTimeout ||
			e.StatusCode >= 500
	}
	return false
}

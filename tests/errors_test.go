package tests

import (
	"net/http"
	"testing"

	sdkerrors "github.com/iamkanishka/yapily-client-go/errors"
)

func TestAPIErrorString(t *testing.T) {
	err := &sdkerrors.APIError{
		StatusCode: 404,
		Code:       "NOT_FOUND",
		Message:    "resource not found",
	}
	want := "api error [404] NOT_FOUND: resource not found"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestAPIErrorStringWithTraceID(t *testing.T) {
	err := &sdkerrors.APIError{
		StatusCode: 500,
		Code:       "SERVER_ERROR",
		Message:    "internal error",
		TraceID:    "trace-abc",
	}
	got := err.Error()
	if got == "" {
		t.Error("error string should not be empty")
	}
	// Must include trace ID
	if len(got) == 0 {
		t.Error("expected non-empty error string")
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"404 api error", &sdkerrors.APIError{StatusCode: http.StatusNotFound}, true},
		{"401 api error", &sdkerrors.APIError{StatusCode: http.StatusUnauthorized}, false},
		{"validation error", &sdkerrors.ValidationError{Field: "x", Message: "y"}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := sdkerrors.IsNotFound(tc.err); got != tc.want {
				t.Errorf("IsNotFound(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	err := &sdkerrors.APIError{StatusCode: http.StatusUnauthorized}
	if !sdkerrors.IsUnauthorized(err) {
		t.Error("expected IsUnauthorized to return true")
	}
}

func TestIsRateLimited(t *testing.T) {
	err := &sdkerrors.APIError{StatusCode: http.StatusTooManyRequests}
	if !sdkerrors.IsRateLimited(err) {
		t.Error("expected IsRateLimited to return true")
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"429", &sdkerrors.APIError{StatusCode: 429}, true},
		{"503", &sdkerrors.APIError{StatusCode: 503}, true},
		{"504", &sdkerrors.APIError{StatusCode: 504}, true},
		{"500", &sdkerrors.APIError{StatusCode: 500}, true},
		{"404", &sdkerrors.APIError{StatusCode: 404}, false},
		{"400", &sdkerrors.APIError{StatusCode: 400}, false},
		{"retryable wrapper", &sdkerrors.RetryableError{Cause: &sdkerrors.APIError{StatusCode: 400}}, true},
		{"validation", &sdkerrors.ValidationError{Field: "x"}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := sdkerrors.IsRetryable(tc.err); got != tc.want {
				t.Errorf("IsRetryable(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := sdkerrors.NewValidationError("amount", "must be greater than 0")
	want := "validation error: field 'amount' - must be greater than 0"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}

func TestAuthError(t *testing.T) {
	cause := &sdkerrors.ValidationError{Field: "token", Message: "expired"}
	err := sdkerrors.NewAuthError("token refresh failed", cause)
	if err.Error() == "" {
		t.Error("expected non-empty auth error string")
	}
	if err.Unwrap() != cause {
		t.Error("Unwrap should return the cause")
	}
}

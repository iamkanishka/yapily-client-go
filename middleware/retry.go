package middleware

import (
	"bytes"
	"context"
	"io"
	"math"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// RetryRoundTripper retries requests on transient failures with exponential backoff.
type RetryRoundTripper struct {
	inner       http.RoundTripper
	maxAttempts int
	baseDelay   time.Duration
	maxDelay    time.Duration
	logger      *zap.Logger
}

// NewRetryRoundTripper creates a RetryRoundTripper.
func NewRetryRoundTripper(inner http.RoundTripper, maxAttempts int, base, max time.Duration, logger *zap.Logger) http.RoundTripper {
	if inner == nil {
		inner = http.DefaultTransport
	}
	return &RetryRoundTripper{
		inner:       inner,
		maxAttempts: maxAttempts,
		baseDelay:   base,
		maxDelay:    max,
		logger:      logger,
	}
}

// RoundTrip executes the request, retrying on 429/5xx.
func (r *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
	}

	var (
		resp    *http.Response
		lastErr error
	)

	for attempt := 0; attempt < r.maxAttempts; attempt++ {
		if attempt > 0 {
			delay := r.backoff(attempt)
			r.logger.Debug("retry backoff",
				zap.Int("attempt", attempt),
				zap.Duration("delay", delay),
			)
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(delay):
			}
		}

		clonedReq := r.cloneRequest(req, bodyBytes)
		resp, lastErr = r.inner.RoundTrip(clonedReq)
		if lastErr != nil {
			if isContextError(req.Context()) {
				return nil, lastErr
			}
			continue
		}

		if !isRetryableStatus(resp.StatusCode) {
			return resp, nil
		}

		resp.Body.Close()
	}

	return resp, lastErr
}

func (r *RetryRoundTripper) backoff(attempt int) time.Duration {
	d := r.baseDelay * time.Duration(math.Pow(2, float64(attempt-1)))
	if d > r.maxDelay {
		d = r.maxDelay
	}
	return d
}

func (r *RetryRoundTripper) cloneRequest(req *http.Request, body []byte) *http.Request {
	clone := req.Clone(req.Context())
	if body != nil {
		clone.Body = io.NopCloser(bytes.NewReader(body))
		clone.ContentLength = int64(len(body))
	}
	return clone
}

func isRetryableStatus(code int) bool {
	return code == http.StatusTooManyRequests ||
		code == http.StatusServiceUnavailable ||
		code == http.StatusGatewayTimeout ||
		code >= 500
}

func isContextError(ctx context.Context) bool {
	return ctx.Err() != nil
}

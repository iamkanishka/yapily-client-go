// Package middleware provides composable HTTP middleware for the Yapily SDK.
package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// loggingRoundTripper wraps an http.RoundTripper with request/response logging.
type loggingRoundTripper struct {
	inner  http.RoundTripper
	logger *zap.Logger
}

// NewLoggingRoundTripper returns an http.RoundTripper that logs all requests.
func NewLoggingRoundTripper(inner http.RoundTripper, logger *zap.Logger) http.RoundTripper {
	if inner == nil {
		inner = http.DefaultTransport
	}
	return &loggingRoundTripper{inner: inner, logger: logger}
}

// RoundTrip executes the request and logs timing and status.
func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	l.logger.Debug("outgoing request",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
	)

	resp, err := l.inner.RoundTrip(req)
	elapsed := time.Since(start)

	if err != nil {
		l.logger.Error("request error",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Duration("elapsed", elapsed),
			zap.Error(err),
		)
		return nil, err
	}

	l.logger.Info("request completed",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Int("status", resp.StatusCode),
		zap.Duration("elapsed", elapsed),
	)

	return resp, nil
}

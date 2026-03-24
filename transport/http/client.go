// Package http provides the low-level HTTP transport for the Yapily SDK.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"

	sdkerrors "github.com/iamkanishka/yapily-client-go/errors"
)

// Client is the HTTP interface satisfied by *http.Client and any test double.
// Named Client rather than HTTPClient to avoid the stutter http.HTTPClient.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// RetryConfig controls exponential-backoff retry behaviour.
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultRetryConfig returns conservative retry defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   200 * time.Millisecond,
		MaxDelay:    5 * time.Second,
	}
}

// Transport is the core HTTP transport with retry, rate-limiting, and logging.
type Transport struct {
	baseURL    string
	httpClient Client
	logger     *zap.Logger
	retry      RetryConfig
	limiter    *rate.Limiter
}

// TransportOption is a functional option for Transport.
type TransportOption func(*Transport)

// WithRetryConfig sets the retry configuration.
func WithRetryConfig(cfg RetryConfig) TransportOption {
	return func(t *Transport) { t.retry = cfg }
}

// WithRateLimit configures a token-bucket rate limiter.
func WithRateLimit(rps float64, burst int) TransportOption {
	return func(t *Transport) { t.limiter = rate.NewLimiter(rate.Limit(rps), burst) }
}

// WithHTTPClient replaces the underlying HTTP client (useful in tests).
func WithHTTPClient(c Client) TransportOption {
	return func(t *Transport) { t.httpClient = c }
}

// WithLogger injects a structured logger.
func WithLogger(l *zap.Logger) TransportOption {
	return func(t *Transport) { t.logger = l }
}

// NewTransport creates a Transport. A production logger is created automatically;
// if that fails (rare), a no-op logger is used so the caller always gets a valid
// Transport.
func NewTransport(baseURL string, opts ...TransportOption) *Transport {
	logger, err := zap.NewProduction()
	if err != nil {
		logger = zap.NewNop()
	}

	t := &Transport{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
		retry:      DefaultRetryConfig(),
		limiter:    rate.NewLimiter(rate.Limit(10), 20),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Request executes an HTTP request with retry, rate-limiting, and structured
// logging. out may be *[]byte (raw body), a JSON-decodable pointer, or nil.
func (t *Transport) Request(
	ctx context.Context,
	method, path string,
	body interface{},
	headers map[string]string,
	out interface{},
) error {
	endpoint := t.baseURL + path

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
	}

	var lastErr error
	for attempt := 0; attempt < t.retry.MaxAttempts; attempt++ {
		if attempt > 0 {
			delay := t.backoff(attempt)
			t.logger.Debug("retrying request",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("attempt", attempt+1),
				zap.Duration("delay", delay),
			)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		if err := t.limiter.Wait(ctx); err != nil {
			return fmt.Errorf("rate limiter: %w", err)
		}

		req, err := t.buildRequest(ctx, method, endpoint, bodyBytes, headers)
		if err != nil {
			return err
		}

		start := time.Now()
		resp, err := t.httpClient.Do(req)
		elapsed := time.Since(start)

		if err != nil {
			t.logger.Warn("http request error",
				zap.String("method", method),
				zap.String("path", path),
				zap.Error(err),
			)
			lastErr = err
			continue
		}

		t.logger.Info("http request completed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", resp.StatusCode),
			zap.Duration("duration", elapsed),
		)

		apiErr, decodeErr := t.handleResponse(resp, out)
		if decodeErr != nil {
			return decodeErr
		}
		if apiErr != nil {
			if sdkerrors.IsRetryable(apiErr) {
				lastErr = apiErr
				continue
			}
			return apiErr
		}
		return nil
	}
	return lastErr
}

// RequestWithQuery appends query parameters to path and performs a GET.
func (t *Transport) RequestWithQuery(
	ctx context.Context,
	path string,
	params map[string]string,
	headers map[string]string,
	out interface{},
) error {
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			if v != "" {
				q.Set(k, v)
			}
		}
		if encoded := q.Encode(); encoded != "" {
			path = path + "?" + encoded
		}
	}
	return t.Request(ctx, http.MethodGet, path, nil, headers, out)
}

func (t *Transport) buildRequest(
	ctx context.Context,
	method, endpoint string,
	body []byte,
	headers map[string]string,
) (*http.Request, error) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, r)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

func (t *Transport) handleResponse(resp *http.Response, out interface{}) (*sdkerrors.APIError, error) {
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &sdkerrors.APIError{
			StatusCode: resp.StatusCode,
			TraceID:    resp.Header.Get("X-Trace-Id"),
		}
		var payload struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if jerr := json.Unmarshal(respBody, &payload); jerr == nil && payload.Error.Message != "" {
			apiErr.Code = payload.Error.Code
			apiErr.Message = payload.Error.Message
		} else {
			apiErr.Code = strconv.Itoa(resp.StatusCode)
			apiErr.Message = string(respBody)
		}
		return apiErr, nil
	}

	if out == nil || len(respBody) == 0 {
		return nil, nil
	}

	// Support raw []byte output (e.g. PDF statement downloads).
	if rawOut, ok := out.(*[]byte); ok {
		*rawOut = respBody
		return nil, nil
	}

	if jerr := json.Unmarshal(respBody, out); jerr != nil {
		return nil, fmt.Errorf("decode response: %w", jerr)
	}
	return nil, nil
}

// backoff returns the exponential delay for a given attempt (1-based).
func (t *Transport) backoff(attempt int) time.Duration {
	d := t.retry.BaseDelay * time.Duration(math.Pow(2, float64(attempt-1)))
	if d > t.retry.MaxDelay {
		d = t.retry.MaxDelay
	}
	return d
}

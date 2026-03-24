package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"

	"go.uber.org/zap"
)

func TestTransportRetryOn500(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	tr := transporthttp.NewTransport(srv.URL,
		transporthttp.WithRetryConfig(transporthttp.RetryConfig{
			MaxAttempts: 3,
			BaseDelay:   1 * time.Millisecond,
			MaxDelay:    5 * time.Millisecond,
		}),
		transporthttp.WithLogger(zap.NewNop()),
	)

	var out map[string]string
	if err := tr.Request(context.Background(), http.MethodGet, "/test", nil, nil, &out); err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 attempts, got %d", calls)
	}
}

func TestTransportNoRetryOn400(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{"code": "BAD_REQUEST", "message": "invalid"},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	tr := transporthttp.NewTransport(srv.URL,
		transporthttp.WithRetryConfig(transporthttp.RetryConfig{MaxAttempts: 3, BaseDelay: 1 * time.Millisecond}),
		transporthttp.WithLogger(zap.NewNop()),
	)

	if err := tr.Request(context.Background(), http.MethodGet, "/test", nil, nil, nil); err == nil {
		t.Fatal("expected error for 400 response")
	}
	if calls != 1 {
		t.Errorf("expected 1 attempt for 400, got %d", calls)
	}
}

func TestTransportContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		writeJSON(w, http.StatusOK, map[string]string{})
	}))
	defer srv.Close()

	tr := transporthttp.NewTransport(srv.URL, transporthttp.WithLogger(zap.NewNop()))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if err := tr.Request(ctx, http.MethodGet, "/slow", nil, nil, nil); err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestTransportQueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "10" || r.URL.Query().Get("from") != "2024-01-01" {
			http.Error(w, "wrong params", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"ok": "true"})
	}))
	defer srv.Close()

	tr := transporthttp.NewTransport(srv.URL, transporthttp.WithLogger(zap.NewNop()))
	var out map[string]string
	if err := tr.RequestWithQuery(context.Background(), "/items",
		map[string]string{"limit": "10", "from": "2024-01-01"},
		nil, &out,
	); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransportAuthHeaderForwarded(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		writeJSON(w, http.StatusOK, map[string]string{})
	}))
	defer srv.Close()

	tr := transporthttp.NewTransport(srv.URL, transporthttp.WithLogger(zap.NewNop()))
	if err := tr.Request(context.Background(), http.MethodGet, "/check",
		nil, map[string]string{"Authorization": "Basic dGVzdA=="}, nil,
	); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Basic dGVzdA==" {
		t.Errorf("auth header not forwarded: %s", gotAuth)
	}
}

func TestTransportJSONBodySent(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}
	var received payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "decode error: "+err.Error(), http.StatusBadRequest)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "wrong content-type", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{})
	}))
	defer srv.Close()

	tr := transporthttp.NewTransport(srv.URL, transporthttp.WithLogger(zap.NewNop()))
	if err := tr.Request(context.Background(), http.MethodPost, "/data",
		payload{Name: "yapily-sdk"}, nil, nil,
	); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Name != "yapily-sdk" {
		t.Errorf("expected yapily-sdk, got %s", received.Name)
	}
}

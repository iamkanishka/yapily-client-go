package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/iamkanishka/yapily-client-go/domain"
	"github.com/iamkanishka/yapily-client-go/services"

	"go.uber.org/zap"
)

func TestPaymentsCreate(t *testing.T) {
	fixture := domain.APIResponse[domain.Payment]{
		Data: domain.Payment{
			ID: "pay-001", Status: "PENDING", Amount: 50.00, Currency: "GBP",
			Recipient: domain.Recipient{Name: "Bob"}, CreatedAt: time.Now(),
		},
	}
	var gotIdem string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/payments" {
			http.NotFound(w, r)
			return
		}
		gotIdem = r.Header.Get("Idempotency-Key")
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewPaymentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	payment, err := svc.Create(context.Background(), "consent-tok", &domain.PaymentRequest{
		Amount: 50.00, Currency: "GBP",
		Recipient:      domain.Recipient{Name: "Bob"},
		IdempotencyKey: "idem-key-123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.ID != "pay-001" {
		t.Errorf("expected pay-001, got %s", payment.ID)
	}
	if gotIdem != "idem-key-123" {
		t.Errorf("idempotency key not forwarded: %q", gotIdem)
	}
}

func TestPaymentsCreateValidation(t *testing.T) {
	svc := services.NewPaymentsService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())
	cases := []struct {
		name    string
		consent string
		req     *domain.PaymentRequest
	}{
		{"empty consent", "", &domain.PaymentRequest{Amount: 10, Currency: "GBP", Recipient: domain.Recipient{Name: "a"}}},
		{"nil request", "tok", nil},
		{"zero amount", "tok", &domain.PaymentRequest{Amount: 0, Currency: "GBP", Recipient: domain.Recipient{Name: "a"}}},
		{"negative amount", "tok", &domain.PaymentRequest{Amount: -5, Currency: "GBP", Recipient: domain.Recipient{Name: "a"}}},
		{"no currency", "tok", &domain.PaymentRequest{Amount: 10, Recipient: domain.Recipient{Name: "a"}}},
		{"no recipient", "tok", &domain.PaymentRequest{Amount: 10, Currency: "GBP"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := svc.Create(context.Background(), tc.consent, tc.req); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

// TestPaymentsGet verifies GET /payments/{id} — no /details suffix (real Yapily API v12).
func TestPaymentsGet(t *testing.T) {
	fixture := domain.APIResponse[domain.Payment]{
		Data: domain.Payment{ID: "pay-999", Status: "COMPLETED", Amount: 100.00, Currency: "GBP"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/payments/pay-999" {
			t.Errorf("wrong path: %s (want /payments/pay-999)", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewPaymentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	payment, err := svc.Get(context.Background(), "consent", "pay-999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.Status != "COMPLETED" {
		t.Errorf("expected COMPLETED, got %s", payment.Status)
	}
}

func TestPaymentsCreateBodyForwarded(t *testing.T) {
	var received domain.PaymentRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "decode: "+err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, domain.APIResponse[domain.Payment]{Data: domain.Payment{ID: "p1"}})
	}))
	defer srv.Close()

	svc := services.NewPaymentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	if _, err := svc.Create(context.Background(), "consent", &domain.PaymentRequest{
		Amount: 42.50, Currency: "EUR", Recipient: domain.Recipient{Name: "Payee"},
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Amount != 42.50 {
		t.Errorf("expected 42.50, got %f", received.Amount)
	}
	if received.Currency != "EUR" {
		t.Errorf("expected EUR, got %s", received.Currency)
	}
}

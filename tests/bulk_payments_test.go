package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iamkanishka/yapily-client-go/domain"
	"github.com/iamkanishka/yapily-client-go/services"

	"go.uber.org/zap"
)

func TestBulkPaymentsCreate(t *testing.T) {
	fixture := domain.APIResponse[domain.BulkPayment]{
		Data: domain.BulkPayment{ID: "bulk-001", Status: "PENDING"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/bulk-payments" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewBulkPaymentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	bulk, err := svc.Create(context.Background(), "consent-tok", &domain.BulkPaymentRequest{
		Payments: []domain.PaymentRequest{
			{Amount: 10, Currency: "GBP", Recipient: domain.Recipient{Name: "Alice"}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bulk.ID != "bulk-001" {
		t.Errorf("expected bulk-001, got %s", bulk.ID)
	}
}

func TestBulkPaymentsCreateValidation(t *testing.T) {
	svc := services.NewBulkPaymentsService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())
	tests := []struct {
		name         string
		consentToken string
		req          *domain.BulkPaymentRequest
	}{
		{"empty consent", "", &domain.BulkPaymentRequest{Payments: []domain.PaymentRequest{{Amount: 1, Currency: "GBP", Recipient: domain.Recipient{Name: "x"}}}}},
		{"nil request", "tok", nil},
		{"empty payments", "tok", &domain.BulkPaymentRequest{Payments: []domain.PaymentRequest{}}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), tc.consentToken, tc.req)
			if err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}

// GetStatus uses GET /bulk-payments/{bulkPaymentId} (no /file suffix — real Yapily API)
func TestBulkPaymentsGetStatus(t *testing.T) {
	fixture := domain.APIResponse[domain.BulkPaymentStatus]{
		Data: domain.BulkPaymentStatus{
			ID:            "bulk-002",
			StatusDetails: &domain.BulkPaymentStatusDetail{Status: "COMPLETED"},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bulk-payments/bulk-002" {
			t.Errorf("wrong path: got %s, want /bulk-payments/bulk-002", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewBulkPaymentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	status, err := svc.GetStatus(context.Background(), "consent-tok", "bulk-002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.StatusDetails.Status != "COMPLETED" {
		t.Errorf("expected COMPLETED, got %s", status.StatusDetails.Status)
	}
}

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

func TestWebhooksGetCategories(t *testing.T) {
	fixture := domain.APIResponse[[]domain.WebhookCategory]{
		Data: []domain.WebhookCategory{
			{ID: "payments", Name: "Payments", Description: "Payment events"},
			{ID: "consents", Name: "Consents", Description: "Consent events"},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/webhooks/categories" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewWebhooksService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	cats, err := svc.GetCategories(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cats) != 2 {
		t.Errorf("expected 2 categories, got %d", len(cats))
	}
}

func TestWebhooksRegisterEvent(t *testing.T) {
	fixture := domain.APIResponse[domain.WebhookEvent]{
		Data: domain.WebhookEvent{
			ID:              "wh-evt-001",
			EventTypeID:     "payment.status.updated",
			NotificationURL: "https://app.com/wh",
		},
	}
	var capturedBody domain.RegisterWebhookRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/webhooks" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
		_ = capturedBody
	}))
	defer srv.Close()

	svc := services.NewWebhooksService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	evt, err := svc.RegisterEvent(context.Background(), &domain.RegisterWebhookRequest{
		EventTypeID:     "payment.status.updated",
		NotificationURL: "https://app.com/wh",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evt.EventTypeID != "payment.status.updated" {
		t.Errorf("expected payment.status.updated, got %s", evt.EventTypeID)
	}
}

func TestWebhooksDeleteEvent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/webhooks/payment.status.updated" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	svc := services.NewWebhooksService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	if err := svc.DeleteEvent(context.Background(), "payment.status.updated"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWebhooksResetSecret(t *testing.T) {
	fixture := domain.APIResponse[domain.ResetWebhookSecretResponse]{
		Data: domain.ResetWebhookSecretResponse{Secret: "new-secret-xyz"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/webhooks/secret/reset" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewWebhooksService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	resp, err := svc.ResetSecret(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Secret != "new-secret-xyz" {
		t.Errorf("expected new-secret-xyz, got %s", resp.Secret)
	}
}

func TestWebhooksRegisterValidation(t *testing.T) {
	svc := services.NewWebhooksService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())

	tests := []struct {
		name string
		req  *domain.RegisterWebhookRequest
	}{
		{"nil request", nil},
		{"empty eventTypeId", &domain.RegisterWebhookRequest{NotificationURL: "https://x.com"}},
		{"empty notificationUrl", &domain.RegisterWebhookRequest{EventTypeID: "pay.updated"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.RegisterEvent(context.Background(), tc.req)
			if err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}

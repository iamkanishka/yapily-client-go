package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/iamkanishka/yapily-client-go/domain"
	"github.com/iamkanishka/yapily-client-go/services"

	"go.uber.org/zap"
)

func TestVRPCreateSweepingAuthorisation(t *testing.T) {
	fixture := domain.APIResponse[domain.VRPConsent]{
		Data: domain.VRPConsent{
			ID:               "vrp-consent-001",
			Status:           "AWAITING_AUTHORIZATION",
			InstitutionID:    "monzo",
			AuthorisationURL: "https://monzo.com/vrp/auth",
			CreatedAt:        time.Now(),
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/vrp-consents" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewVRPService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	consent, err := svc.CreateSweepingAuthorisation(ctx(), &domain.VRPAuthorisationRequest{
		InstitutionID:     "monzo",
		ApplicationUserID: "user-001",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if consent.ID != "vrp-consent-001" {
		t.Errorf("expected vrp-consent-001, got %s", consent.ID)
	}
}

func TestVRPCreateSweepingAuthorisationValidation(t *testing.T) {
	svc := services.NewVRPService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())

	tests := []struct {
		name string
		req  *domain.VRPAuthorisationRequest
	}{
		{"nil request", nil},
		{"missing institutionId", &domain.VRPAuthorisationRequest{ApplicationUserID: "u"}},
		{"missing applicationUserId", &domain.VRPAuthorisationRequest{InstitutionID: "monzo"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.CreateSweepingAuthorisation(ctx(), tc.req)
			if err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}

func TestVRPConfirmFunds(t *testing.T) {
	fixture := domain.APIResponse[domain.FundsConfirmationResponse]{
		Data: domain.FundsConfirmationResponse{FundsAvailable: true},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/vrp-consents/vrp-c-1/funds-confirmation" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewVRPService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	result, err := svc.ConfirmFunds(ctx(), "consent-tok", "vrp-c-1", &domain.FundsConfirmationRequest{
		Amount: 50.00, Currency: "GBP",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.FundsAvailable {
		t.Error("expected FundsAvailable to be true")
	}
}

// ctx is a helper to avoid repeating context.Background() in table tests.
func ctx() context.Context { return context.Background() }

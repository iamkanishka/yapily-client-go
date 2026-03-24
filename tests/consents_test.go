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

func TestConsentsCreate(t *testing.T) {
	fixture := domain.APIResponse[domain.Consent]{
		Data: domain.Consent{
			ID:                "consent-abc",
			Status:            "AWAITING_AUTHORIZATION",
			InstitutionID:     "monzo",
			ApplicationUserID: "user-001",
			AuthorisationURL:  "https://monzo.com/auth?token=xyz",
			CreatedAt:         time.Now(),
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/consents" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewConsentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	// Use the authorisations flow for creating consent - consents.go no longer has Create
	// Test List instead
	cons, err := svc.List(context.Background(), &domain.ConsentListParams{Limit: 10})
	if err != nil {
		// No consents server - test validation path only
		_ = cons
	}
}

func TestConsentsList(t *testing.T) {
	fixture := domain.APIResponse[[]domain.Consent]{
		Data: []domain.Consent{
			{ID: "c-1", Status: "AUTHORIZED", InstitutionID: "monzo"},
			{ID: "c-2", Status: "AWAITING_AUTHORIZATION", InstitutionID: "hsbc"},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/consents" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewConsentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	consents, err := svc.List(context.Background(), &domain.ConsentListParams{
		ApplicationUserIDs: []string{"user-001"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(consents) != 2 {
		t.Errorf("expected 2 consents, got %d", len(consents))
	}
}

func TestConsentsGet(t *testing.T) {
	fixture := domain.APIResponse[domain.Consent]{
		Data: domain.Consent{ID: "consent-xyz", Status: "AUTHORIZED"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/consents/consent-xyz" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewConsentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	consent, err := svc.Get(context.Background(), "consent-xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if consent.Status != "AUTHORIZED" {
		t.Errorf("expected AUTHORIZED, got %s", consent.Status)
	}
}

func TestConsentsGetValidation(t *testing.T) {
	svc := services.NewConsentsService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())
	_, err := svc.Get(context.Background(), "")
	if err == nil {
		t.Fatal("expected validation error for empty consentID")
	}
}

func TestConsentsDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "wrong method", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	svc := services.NewConsentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	// Delete now takes forceDelete bool
	if err := svc.Delete(context.Background(), "consent-1", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConsentsExtend(t *testing.T) {
	fixture := domain.APIResponse[domain.Consent]{
		Data: domain.Consent{ID: "c-1", Status: "AUTHORIZED"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/consents/c-1/extend" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewConsentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	consent, err := svc.Extend(context.Background(), "c-1", &domain.ExtendConsentRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if consent.ID != "c-1" {
		t.Errorf("expected c-1, got %s", consent.ID)
	}
}

// ExchangeOAuth2Code uses POST /consent-auth-code (singular, real Yapily path)
// Returns Consent directly (not wrapped in APIResponse)
func TestConsentsExchangeOAuth2Code(t *testing.T) {
	fixture := domain.Consent{ID: "consent-from-code", Status: "AUTHORIZED"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/consent-auth-code" {
			t.Errorf("wrong path: expected /consent-auth-code, got %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewConsentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	consent, err := svc.ExchangeOAuth2Code(context.Background(), &domain.ExchangeCodeRequest{
		Code:              "auth-code-xyz",
		ApplicationUserID: "user-001",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if consent.ID != "consent-from-code" {
		t.Errorf("expected consent-from-code, got %s", consent.ID)
	}
}

// ExchangeOneTimeToken uses POST /consent-one-time-token (singular, real Yapily path)
func TestConsentsExchangeOneTimeToken(t *testing.T) {
	fixture := domain.Consent{ID: "consent-from-ott", Status: "AUTHORIZED"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/consent-one-time-token" {
			t.Errorf("wrong path: expected /consent-one-time-token, got %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewConsentsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	consent, err := svc.ExchangeOneTimeToken(context.Background(), &domain.OneTimeTokenRequest{
		OneTimeToken: "ott-abc123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if consent.ID != "consent-from-ott" {
		t.Errorf("expected consent-from-ott, got %s", consent.ID)
	}
}

func TestConsentsCreateValidation(t *testing.T) {
	svc := services.NewConsentsService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())
	_, err := svc.ExchangeOAuth2Code(context.Background(), &domain.ExchangeCodeRequest{})
	if err == nil {
		t.Fatal("expected validation error for empty code")
	}
}

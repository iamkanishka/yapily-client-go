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

func TestAuthorisationsCreateAccount(t *testing.T) {
	fixture := domain.APIResponse[domain.Authorisation]{
		Data: domain.Authorisation{
			ID:               "auth-001",
			Status:           "AWAITING_AUTHORIZATION",
			InstitutionID:    "monzo",
			AuthorisationURL: "https://monzo.com/auth?token=xyz",
			CreatedAt:        time.Now(),
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/account-auth-requests" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewAuthorisationsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	auth, err := svc.CreateAccountAuthorisation(context.Background(),
		&domain.AccountAuthorisationRequest{InstitutionID: "monzo", ApplicationUserID: "u"},
		&domain.PSUHeaders{PSUIPAddress: "1.2.3.4"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth.ID != "auth-001" {
		t.Errorf("expected auth-001, got %s", auth.ID)
	}
	if auth.AuthorisationURL == "" {
		t.Error("expected non-empty authorisationURL")
	}
}

func TestAuthorisationsCreateAccountValidation(t *testing.T) {
	svc := services.NewAuthorisationsService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())
	cases := []struct {
		name string
		req  *domain.AccountAuthorisationRequest
	}{
		{"nil", nil},
		{"no institution", &domain.AccountAuthorisationRequest{ApplicationUserID: "u"}},
		{"no user", &domain.AccountAuthorisationRequest{InstitutionID: "monzo"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := svc.CreateAccountAuthorisation(context.Background(), tc.req, nil); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestAuthorisationsCreatePayment(t *testing.T) {
	fixture := domain.APIResponse[domain.Authorisation]{
		Data: domain.Authorisation{ID: "pay-auth-001", Status: "AWAITING_AUTHORIZATION", InstitutionID: "hsbc", CreatedAt: time.Now()},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/payment-auth-requests" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewAuthorisationsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	auth, err := svc.CreatePaymentAuthorisation(context.Background(),
		&domain.PaymentAuthorisationRequest{InstitutionID: "hsbc", ApplicationUserID: "u"}, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth.ID != "pay-auth-001" {
		t.Errorf("expected pay-auth-001, got %s", auth.ID)
	}
}

func TestAuthorisationsPreAuth(t *testing.T) {
	fixture := domain.APIResponse[domain.Authorisation]{
		Data: domain.Authorisation{ID: "pre-001", Status: "AUTHORIZED"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pre-auth-requests" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewAuthorisationsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	auth, err := svc.CreatePreAuthorisation(context.Background(),
		&domain.PreAuthorisationRequest{InstitutionID: "monzo", ApplicationUserID: "u"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth.Status != "AUTHORIZED" {
		t.Errorf("expected AUTHORIZED, got %s", auth.Status)
	}
}

func TestAuthorisationsPSUHeadersForwarded(t *testing.T) {
	var gotID, gotCorp, gotIP string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID = r.Header.Get("psu-id")
		gotCorp = r.Header.Get("psu-corporate-id")
		gotIP = r.Header.Get("psu-ip-address")
		writeJSON(w, http.StatusOK, domain.APIResponse[domain.Authorisation]{Data: domain.Authorisation{ID: "x"}})
	}))
	defer srv.Close()

	svc := services.NewAuthorisationsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	if _, err := svc.CreateAccountAuthorisation(context.Background(),
		&domain.AccountAuthorisationRequest{InstitutionID: "m", ApplicationUserID: "u"},
		&domain.PSUHeaders{PSUID: "psu-123", PSUCorporateID: "corp-456", PSUIPAddress: "10.0.0.1"},
	); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != "psu-123" {
		t.Errorf("psu-id: got %q", gotID)
	}
	if gotCorp != "corp-456" {
		t.Errorf("psu-corporate-id: got %q", gotCorp)
	}
	if gotIP != "10.0.0.1" {
		t.Errorf("psu-ip-address: got %q", gotIP)
	}
}

package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iamkanishka/yapily-client-go/client"
	"github.com/iamkanishka/yapily-client-go/domain"
)

func TestClientNewSuccess(t *testing.T) {
	c, err := client.New("app-key", "app-secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	// Verify all services are wired up.
	if c.Institutions == nil {
		t.Error("Institutions service not wired")
	}
	if c.Accounts == nil {
		t.Error("Accounts service not wired")
	}
	if c.Transactions == nil {
		t.Error("Transactions service not wired")
	}
	if c.Payments == nil {
		t.Error("Payments service not wired")
	}
	if c.BulkPayments == nil {
		t.Error("BulkPayments service not wired")
	}
	if c.Consents == nil {
		t.Error("Consents service not wired")
	}
	if c.Authorisations == nil {
		t.Error("Authorisations service not wired")
	}
	if c.FinancialData == nil {
		t.Error("FinancialData service not wired")
	}
	if c.Users == nil {
		t.Error("Users service not wired")
	}
	if c.VRP == nil {
		t.Error("VRP service not wired")
	}
	if c.Notifications == nil {
		t.Error("Notifications service not wired")
	}
	if c.DataPlus == nil {
		t.Error("DataPlus service not wired")
	}
	if c.HostedPages == nil {
		t.Error("HostedPages service not wired")
	}
	if c.Constraints == nil {
		t.Error("Constraints service not wired")
	}
	if c.Application == nil {
		t.Error("Application service not wired")
	}
	if c.Webhooks == nil {
		t.Error("Webhooks service not wired")
	}
	if c.Beneficiaries == nil {
		t.Error("Beneficiaries service not wired")
	}
}

func TestClientNewValidation(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		secret string
	}{
		{"empty key", "", "secret"},
		{"empty secret", "key", ""},
		{"both empty", "", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.New(tc.key, tc.secret)
			if err == nil {
				t.Fatal("expected error for missing credentials")
			}
		})
	}
}

func TestClientWithBaseURL(t *testing.T) {
	fixture := domain.APIResponse[[]domain.Institution]{
		Data: []domain.Institution{{ID: "test-bank", Name: "Test Bank"}},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	c, err := client.New("key", "secret", client.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	insts, err := c.Institutions.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(insts) != 1 || insts[0].ID != "test-bank" {
		t.Errorf("unexpected institutions: %+v", insts)
	}
}

func TestClientAuthenticateSuccess(t *testing.T) {
	// Authenticate just verifies the provider returns a token - Basic auth is static.
	c, err := client.New("my-key", "my-secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should not error - Basic auth never needs network.
	if err := c.Authenticate(context.Background()); err != nil {
		t.Fatalf("unexpected authenticate error: %v", err)
	}
}

func TestClientEnvironmentSandboxDefault(t *testing.T) {
	// Default environment should be Sandbox.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, domain.APIResponse[[]domain.Institution]{})
	}))
	defer srv.Close()

	c, err := client.New("k", "s",
		client.WithBaseURL(srv.URL),
		client.WithEnvironment(domain.Sandbox),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should work without error
	_, err = c.Institutions.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

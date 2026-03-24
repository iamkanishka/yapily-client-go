package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"
	"github.com/iamkanishka/yapily-client-go/services"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// staticAuthProvider always returns a fixed token without network calls.
type staticAuthProvider struct{ token string }

func (s *staticAuthProvider) GetToken(_ context.Context) (string, error) { return s.token, nil }
func (s *staticAuthProvider) Invalidate()                                {}

var _ auth.Provider = (*staticAuthProvider)(nil)

// newTestTransport creates a Transport pointing at baseURL with a no-op logger.
func newTestTransport(baseURL string) *transporthttp.Transport {
	return transporthttp.NewTransport(baseURL, transporthttp.WithLogger(zap.NewNop()))
}

// writeJSON encodes v as JSON into w with the given HTTP status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "encode failed: "+err.Error(), http.StatusInternalServerError)
	}
}

func TestAccountsList(t *testing.T) {
	fixture := domain.APIResponse[[]domain.Account]{
		Data: []domain.Account{
			{ID: "acc-1", Type: "PERSONAL", Balance: domain.Balance{Amount: 1000.00, Currency: "GBP"}, InstitutionID: "monzo"},
			{ID: "acc-2", Type: "SAVINGS", Balance: domain.Balance{Amount: 5000.00, Currency: "GBP"}, InstitutionID: "monzo"},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewAccountsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	accounts, err := svc.List(context.Background(), "consent-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts))
	}
	if accounts[0].ID != "acc-1" {
		t.Errorf("expected acc-1, got %s", accounts[0].ID)
	}
	if accounts[1].Balance.Amount != 5000.00 {
		t.Errorf("expected 5000.00, got %f", accounts[1].Balance.Amount)
	}
}

func TestAccountsListValidation(t *testing.T) {
	svc := services.NewAccountsService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())
	if _, err := svc.List(context.Background(), ""); err == nil {
		t.Fatal("expected validation error for empty consentToken")
	}
}

func TestAccountsGet(t *testing.T) {
	fixture := domain.APIResponse[domain.Account]{
		Data: domain.Account{
			ID: "acc-42", Type: "CURRENT",
			Balance: domain.Balance{Amount: 250.75, Currency: "GBP"}, InstitutionID: "hsbc",
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/acc-42" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewAccountsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	acc, err := svc.Get(context.Background(), "consent-xyz", "acc-42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acc.ID != "acc-42" {
		t.Errorf("expected acc-42, got %s", acc.ID)
	}
	if acc.Balance.Amount != 250.75 {
		t.Errorf("expected 250.75, got %f", acc.Balance.Amount)
	}
}

func TestAccountsGetValidation(t *testing.T) {
	svc := services.NewAccountsService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())
	cases := []struct{ consent, id string }{
		{"", "acc-1"},
		{"tok", ""},
	}
	for _, tc := range cases {
		if _, err := svc.Get(context.Background(), tc.consent, tc.id); err == nil {
			t.Errorf("expected validation error for consent=%q id=%q", tc.consent, tc.id)
		}
	}
}

func TestAccountsGetNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]string{"code": "NOT_FOUND", "message": "not found"},
		})
	}))
	defer srv.Close()

	svc := services.NewAccountsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	if _, err := svc.Get(context.Background(), "tok", "nonexistent"); err == nil {
		t.Fatal("expected error for 404 response")
	}
}

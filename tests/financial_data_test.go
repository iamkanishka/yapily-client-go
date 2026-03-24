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

func TestFinancialDataGetBalances(t *testing.T) {
	fixture := domain.APIResponse[domain.AccountBalance]{
		Data: domain.AccountBalance{
			AccountID: "acc-1",
			Balances: []domain.BalanceDetail{
				{Type: "AVAILABLE", Amount: 1000.00, Currency: "GBP"},
				{Type: "CURRENT", Amount: 950.00, Currency: "GBP"},
			},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/acc-1/balances" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewFinancialDataService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	bal, err := svc.GetAccountBalances(context.Background(), "consent", "acc-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bal.Balances) != 2 {
		t.Errorf("expected 2 balances, got %d", len(bal.Balances))
	}
}

func TestFinancialDataGetDirectDebits(t *testing.T) {
	fixture := domain.APIResponse[[]domain.DirectDebit]{
		Data: []domain.DirectDebit{
			{ID: "dd-1", Name: "Netflix"},
			{ID: "dd-2", Name: "Gym"},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/acc-1/direct-debits" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewFinancialDataService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	dds, err := svc.GetAccountDirectDebits(context.Background(), "consent", "acc-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dds) != 2 {
		t.Errorf("expected 2 direct debits, got %d", len(dds))
	}
}

func TestFinancialDataGetIdentity(t *testing.T) {
	fixture := domain.APIResponse[domain.Identity]{
		Data: domain.Identity{
			FullName: "Jane Doe",
			Emails:   []domain.Email{{Address: "jane@example.com"}},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/identity" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewFinancialDataService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	id, err := svc.GetIdentity(context.Background(), "consent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.FullName != "Jane Doe" {
		t.Errorf("expected Jane Doe, got %s", id.FullName)
	}
}

func TestFinancialDataValidation(t *testing.T) {
	svc := services.NewFinancialDataService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())

	tests := []struct {
		name string
		fn   func() error
	}{
		{"balances empty consent", func() error { _, e := svc.GetAccountBalances(context.Background(), "", "acc-1"); return e }},
		{"balances empty account", func() error { _, e := svc.GetAccountBalances(context.Background(), "tok", ""); return e }},
		{"identity empty consent", func() error { _, e := svc.GetIdentity(context.Background(), ""); return e }},
		{"direct debits empty consent", func() error { _, e := svc.GetAccountDirectDebits(context.Background(), "", "acc"); return e }},
		{"scheduled payments empty account", func() error { _, e := svc.GetAccountScheduledPayments(context.Background(), "tok", ""); return e }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.fn(); err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}

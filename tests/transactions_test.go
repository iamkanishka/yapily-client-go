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

func TestTransactionsList(t *testing.T) {
	fixture := domain.APIResponse[[]domain.Transaction]{
		Data: []domain.Transaction{
			{ID: "tx-1", Amount: 25.50, Currency: "GBP", Description: "Coffee", Date: "2024-03-01"},
			{ID: "tx-2", Amount: -100.00, Currency: "GBP", Description: "Rent", Date: "2024-03-02"},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/acc-1/transactions" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewTransactionsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())

	txns, err := svc.List(context.Background(), "consent", "acc-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txns) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(txns))
	}
	if txns[0].ID != "tx-1" {
		t.Errorf("expected tx-1, got %s", txns[0].ID)
	}
}

func TestTransactionsListValidation(t *testing.T) {
	svc := services.NewTransactionsService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())

	tests := []struct {
		name         string
		consentToken string
		accountID    string
	}{
		{"empty consent", "", "acc-1"},
		{"empty accountID", "tok", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.List(context.Background(), tc.consentToken, tc.accountID, nil)
			if err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}

func TestTransactionsListWithQueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("limit") != "5" {
			http.Error(w, "missing limit param", http.StatusBadRequest)
			return
		}
		if q.Get("from") != "2024-01-01" {
			http.Error(w, "missing from param", http.StatusBadRequest)
			return
		}
		fixture := domain.APIResponse[[]domain.Transaction]{
			Data: []domain.Transaction{
				{ID: "tx-filtered", Amount: 10.00, Currency: "GBP", Date: "2024-01-15"},
			},
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewTransactionsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())

	params := &domain.TransactionQueryParams{
		PaginationParams: domain.PaginationParams{Limit: 5},
		From:             "2024-01-01",
	}
	txns, err := svc.List(context.Background(), "consent", "acc-1", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txns) != 1 || txns[0].ID != "tx-filtered" {
		t.Fatalf("unexpected transactions: %+v", txns)
	}
}

func TestTransactionsListPaginated(t *testing.T) {
	page1 := domain.APIResponse[[]domain.Transaction]{
		Data: []domain.Transaction{
			{ID: "tx-p1-1"}, {ID: "tx-p1-2"}, {ID: "tx-p1-3"},
		},
	}
	page2 := domain.APIResponse[[]domain.Transaction]{
		Data: []domain.Transaction{
			{ID: "tx-p2-1"}, {ID: "tx-p2-2"},
		},
	}

	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			writeJSON(w, http.StatusOK, page1)
		} else {
			writeJSON(w, http.StatusOK, page2)
		}
	}))
	defer srv.Close()

	svc := services.NewTransactionsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())

	var collected []domain.Transaction
	err := svc.ListPaginated(context.Background(), "consent", "acc-1", 3, func(page []domain.Transaction) bool {
		collected = append(collected, page...)
		return true
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// page1 has 3 (==pageSize so continues), page2 has 2 (<pageSize so stops)
	if len(collected) != 5 {
		t.Errorf("expected 5 transactions, got %d", len(collected))
	}
}

func TestTransactionsListPaginatedEarlyStop(t *testing.T) {
	page := domain.APIResponse[[]domain.Transaction]{
		Data: []domain.Transaction{{ID: "tx-1"}, {ID: "tx-2"}, {ID: "tx-3"}},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, page)
	}))
	defer srv.Close()

	svc := services.NewTransactionsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())

	calls := 0
	err := svc.ListPaginated(context.Background(), "consent", "acc-1", 3, func(_ []domain.Transaction) bool {
		calls++
		return false // stop after first page
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 callback call, got %d", calls)
	}
}

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

func TestConstraintsGetPayment(t *testing.T) {
	fixture := domain.APIResponse[[]domain.PaymentConstraint]{
		Data: []domain.PaymentConstraint{
			{InstitutionID: "monzo", PaymentType: "DOMESTIC_PAYMENT", Currencies: []string{"GBP"}},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/constraints/payment" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("institution-id") != "monzo" {
			http.Error(w, "missing institution-id", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewConstraintsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	cs, err := svc.GetPaymentConstraints(context.Background(), "monzo", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 {
		t.Errorf("expected 1 constraint, got %d", len(cs))
	}
	if cs[0].InstitutionID != "monzo" {
		t.Errorf("expected monzo, got %s", cs[0].InstitutionID)
	}
}

func TestConstraintsGetData(t *testing.T) {
	fixture := domain.APIResponse[[]domain.DataConstraint]{
		Data: []domain.DataConstraint{
			{InstitutionID: "starling", MaxDaysHistory: 365},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/constraints/data" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewConstraintsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	cs, err := svc.GetDataConstraints(context.Background(), "starling")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 1 {
		t.Errorf("expected 1 data constraint, got %d", len(cs))
	}
	if cs[0].MaxDaysHistory != 365 {
		t.Errorf("expected 365 days, got %d", cs[0].MaxDaysHistory)
	}
}

func TestConstraintsNoFilter(t *testing.T) {
	fixture := domain.APIResponse[[]domain.PaymentConstraint]{
		Data: []domain.PaymentConstraint{
			{InstitutionID: "monzo"},
			{InstitutionID: "starling"},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewConstraintsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	// No filter — should return all
	cs, err := svc.GetPaymentConstraints(context.Background(), "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cs) != 2 {
		t.Errorf("expected 2 constraints, got %d", len(cs))
	}
}

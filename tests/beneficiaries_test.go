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

func TestApplicationBeneficiariesCreate(t *testing.T) {
	fixture := domain.APIResponse[domain.ApplicationBeneficiary]{
		Data: domain.ApplicationBeneficiary{
			ID:   "bene-001",
			Name: "Corp Payroll",
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/application/beneficiaries" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewBeneficiariesService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	b, err := svc.CreateApplicationBeneficiary(context.Background(), &domain.CreateBeneficiaryRequest{
		Name: "Corp Payroll",
		AccountIdentifications: []domain.AccountIdentification{
			{Type: "SORT_CODE", Identification: "200000"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.ID != "bene-001" {
		t.Errorf("expected bene-001, got %s", b.ID)
	}
}

func TestApplicationBeneficiariesCreateValidation(t *testing.T) {
	svc := services.NewBeneficiariesService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())

	tests := []struct {
		name string
		req  *domain.CreateBeneficiaryRequest
	}{
		{"nil request", nil},
		{"empty name", &domain.CreateBeneficiaryRequest{
			AccountIdentifications: []domain.AccountIdentification{{Type: "SORT_CODE", Identification: "200000"}},
		}},
		{"empty accountIdentifications", &domain.CreateBeneficiaryRequest{Name: "Bob"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.CreateApplicationBeneficiary(context.Background(), tc.req)
			if err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}

func TestApplicationBeneficiariesDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/application/beneficiaries/bene-del" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	svc := services.NewBeneficiariesService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	if err := svc.DeleteApplicationBeneficiary(context.Background(), "bene-del"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserBeneficiariesList(t *testing.T) {
	fixture := domain.APIResponse[[]domain.UserBeneficiary]{
		Data: []domain.UserBeneficiary{
			{ID: "ub-1", Name: "Jane", Status: "ACTIVE"},
			{ID: "ub-2", Name: "Joe", Status: "PENDING"},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/uuid-001/beneficiaries" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewBeneficiariesService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	bs, err := svc.ListUserBeneficiaries(context.Background(), "uuid-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bs) != 2 {
		t.Errorf("expected 2 beneficiaries, got %d", len(bs))
	}
}

func TestUserBeneficiariesApprove(t *testing.T) {
	fixture := domain.APIResponse[domain.UserBeneficiary]{
		Data: domain.UserBeneficiary{ID: "ub-3", Status: "ACTIVE"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/uuid-001/beneficiaries/ub-3/approve" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewBeneficiariesService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	b, err := svc.ApproveBeneficiary(context.Background(), "uuid-001", "ub-3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Status != "ACTIVE" {
		t.Errorf("expected ACTIVE, got %s", b.Status)
	}
}

func TestBeneficiariesValidation(t *testing.T) {
	svc := services.NewBeneficiariesService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())

	tests := []struct {
		name string
		fn   func() error
	}{
		{"list user empty userUUID", func() error { _, e := svc.ListUserBeneficiaries(context.Background(), ""); return e }},
		{"get user bene empty beneID", func() error { _, e := svc.GetUserBeneficiary(context.Background(), "u", ""); return e }},
		{"delete app bene empty id", func() error { return svc.DeleteApplicationBeneficiary(context.Background(), "") }},
		{"approve empty userUUID", func() error { _, e := svc.ApproveBeneficiary(context.Background(), "", "b"); return e }},
		{"reject empty beneID", func() error { _, e := svc.RejectBeneficiary(context.Background(), "u", ""); return e }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.fn(); err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}

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

func TestInstitutionsList(t *testing.T) {
	fixture := domain.APIResponse[[]domain.Institution]{
		Data: []domain.Institution{
			{ID: "monzo", Name: "Monzo", FullName: "Monzo Bank Ltd"},
			{ID: "starling", Name: "Starling", FullName: "Starling Bank"},
			{ID: "hsbc", Name: "HSBC", FullName: "HSBC Bank plc"},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/institutions" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewInstitutionsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())

	insts, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(insts) != 3 {
		t.Fatalf("expected 3 institutions, got %d", len(insts))
	}

	names := map[string]bool{}
	for _, i := range insts {
		names[i.ID] = true
	}
	for _, expected := range []string{"monzo", "starling", "hsbc"} {
		if !names[expected] {
			t.Errorf("expected institution %s not found", expected)
		}
	}
}

func TestInstitutionsGet(t *testing.T) {
	fixture := domain.APIResponse[domain.Institution]{
		Data: domain.Institution{
			ID:       "monzo",
			Name:     "Monzo",
			FullName: "Monzo Bank Ltd",
			Countries: []domain.Country{
				{CountryCode: "GB", DisplayName: "United Kingdom"},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/institutions/monzo" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewInstitutionsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())

	inst, err := svc.Get(context.Background(), "monzo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inst.ID != "monzo" {
		t.Errorf("expected monzo, got %s", inst.ID)
	}
	if len(inst.Countries) != 1 || inst.Countries[0].CountryCode != "GB" {
		t.Error("expected GB country")
	}
}

func TestInstitutionsGetValidation(t *testing.T) {
	svc := services.NewInstitutionsService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())

	_, err := svc.Get(context.Background(), "")
	if err == nil {
		t.Fatal("expected validation error for empty institutionID")
	}
}

func TestInstitutionsServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]string{
				"code":    "INTERNAL_ERROR",
				"message": "something went wrong",
			},
		})
	}))
	defer srv.Close()

	// Use retry config with 1 attempt to keep test fast
	transport := newTestTransport(srv.URL)
	svc := services.NewInstitutionsService(transport, &staticAuthProvider{"tok"}, zap.NewNop())

	_, err := svc.List(context.Background())
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

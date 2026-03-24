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

func TestUsersCreate(t *testing.T) {
	fixture := domain.APIResponse[domain.User]{
		Data: domain.User{
			UUID:              "uuid-abc-123",
			ApplicationUserID: "user-demo",
			CreatedAt:         time.Now(),
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/users" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewUsersService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	user, err := svc.Create(context.Background(), &domain.CreateUserRequest{ApplicationUserID: "user-demo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.UUID != "uuid-abc-123" {
		t.Errorf("expected uuid-abc-123, got %s", user.UUID)
	}
}

func TestUsersCreateValidation(t *testing.T) {
	svc := services.NewUsersService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())
	_, err := svc.Create(context.Background(), &domain.CreateUserRequest{})
	if err == nil {
		t.Fatal("expected validation error for empty applicationUserId")
	}
}

func TestUsersGet(t *testing.T) {
	fixture := domain.APIResponse[domain.User]{
		Data: domain.User{UUID: "uuid-xyz", ApplicationUserID: "user-xyz"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/uuid-xyz" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewUsersService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	user, err := svc.Get(context.Background(), "uuid-xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.UUID != "uuid-xyz" {
		t.Errorf("expected uuid-xyz, got %s", user.UUID)
	}
}

func TestUsersDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/users/uuid-del" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	svc := services.NewUsersService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	if err := svc.Delete(context.Background(), "uuid-del"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUsersDeleteValidation(t *testing.T) {
	svc := services.NewUsersService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())
	if err := svc.Delete(context.Background(), ""); err == nil {
		t.Fatal("expected validation error for empty userUUID")
	}
}

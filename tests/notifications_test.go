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

// Notifications use /event-subscriptions (NOT /notifications/event-subscriptions)
// This is the correct Yapily API v12 path.

func TestNotificationsListEventSubscriptions(t *testing.T) {
	fixture := domain.APIResponse[[]domain.EventSubscription]{
		Data: []domain.EventSubscription{
			{ID: "evt-1", EventTypeID: "payment.status.updated"},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/event-subscriptions" {
			t.Errorf("wrong path: got %s, want /event-subscriptions", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewNotificationsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	subs, err := svc.ListEventSubscriptions(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(subs) != 1 {
		t.Errorf("expected 1 subscription, got %d", len(subs))
	}
}

func TestNotificationsCreateEventSubscription(t *testing.T) {
	fixture := domain.APIResponse[domain.EventSubscription]{
		Data: domain.EventSubscription{ID: "evt-new", EventTypeID: "payment.status.updated"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/event-subscriptions" {
			t.Errorf("wrong method/path: %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewNotificationsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	sub, err := svc.CreateEventSubscription(context.Background(), &domain.CreateEventSubscriptionRequest{
		EventTypeID:     "payment.status.updated",
		NotificationURL: "https://app.com/wh",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.EventTypeID != "payment.status.updated" {
		t.Errorf("expected payment.status.updated, got %s", sub.EventTypeID)
	}
}

func TestNotificationsGetEventSubscription(t *testing.T) {
	fixture := domain.APIResponse[domain.EventSubscription]{
		Data: domain.EventSubscription{ID: "evt-1", EventTypeID: "payment.status.updated"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/event-subscriptions/payment.status.updated" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	}))
	defer srv.Close()

	svc := services.NewNotificationsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	sub, err := svc.GetEventSubscription(context.Background(), "payment.status.updated")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.EventTypeID != "payment.status.updated" {
		t.Errorf("expected payment.status.updated, got %s", sub.EventTypeID)
	}
}

func TestNotificationsDeleteEventSubscription(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/event-subscriptions/payment.status.updated" {
			t.Errorf("wrong: %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	svc := services.NewNotificationsService(newTestTransport(srv.URL), &staticAuthProvider{"tok"}, zap.NewNop())
	if err := svc.DeleteEventSubscription(context.Background(), "payment.status.updated"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNotificationsValidation(t *testing.T) {
	svc := services.NewNotificationsService(newTestTransport("http://localhost"), &staticAuthProvider{"tok"}, zap.NewNop())
	_, err := svc.CreateEventSubscription(context.Background(), nil)
	if err == nil {
		t.Fatal("expected validation error for nil request")
	}
	_, err = svc.CreateEventSubscription(context.Background(), &domain.CreateEventSubscriptionRequest{})
	if err == nil {
		t.Fatal("expected validation error for empty eventTypeId")
	}
}

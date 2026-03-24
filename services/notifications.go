package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// NotificationsService manages event subscriptions (webhook notifications).
type NotificationsService struct{ base }

// NewNotificationsService creates a new NotificationsService.
func NewNotificationsService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *NotificationsService {
	return &NotificationsService{base{transport: t, auth: a, logger: l}}
}

// ListEventSubscriptions returns all event subscriptions for this application.
// GET /event-subscriptions.
func (s *NotificationsService) ListEventSubscriptions(ctx context.Context) ([]domain.EventSubscription, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[[]domain.EventSubscription]
	if err := s.get(ctx, "/event-subscriptions", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// CreateEventSubscription creates a new event subscription.
// POST /event-subscriptions.
func (s *NotificationsService) CreateEventSubscription(ctx context.Context, req *domain.CreateEventSubscriptionRequest) (*domain.EventSubscription, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.EventTypeID == "" {
		return nil, newValidationErr("eventTypeId", "must not be empty")
	}
	if req.NotificationURL == "" {
		return nil, newValidationErr("notificationUrl", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.EventSubscription]
	if err := s.post(ctx, "/event-subscriptions", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetEventSubscription returns a single event subscription.
// GET /event-subscriptions/{eventTypeId}.
func (s *NotificationsService) GetEventSubscription(ctx context.Context, eventTypeID string) (*domain.EventSubscription, error) {
	if eventTypeID == "" {
		return nil, newValidationErr("eventTypeId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.EventSubscription]
	if err := s.get(ctx, "/event-subscriptions/"+eventTypeID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// DeleteEventSubscription removes an event subscription.
// DELETE /event-subscriptions/{eventTypeId}.
func (s *NotificationsService) DeleteEventSubscription(ctx context.Context, eventTypeID string) error {
	if eventTypeID == "" {
		return newValidationErr("eventTypeId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return err
	}
	return s.del(ctx, "/event-subscriptions/"+eventTypeID, h)
}

package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// WebhooksService manages webhook registrations and categories.
type WebhooksService struct{ base }

// NewWebhooksService creates a new WebhooksService.
func NewWebhooksService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *WebhooksService {
	return &WebhooksService{base{transport: t, auth: a, logger: l}}
}

// GetCategories lists all available webhook event categories.
// GET /webhooks/categories.
func (s *WebhooksService) GetCategories(ctx context.Context) ([]domain.WebhookCategory, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[[]domain.WebhookCategory]
	if err := s.get(ctx, "/webhooks/categories", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ListEvents returns all registered webhook events.
// GET /webhooks.
func (s *WebhooksService) ListEvents(ctx context.Context) ([]domain.WebhookEvent, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[[]domain.WebhookEvent]
	if err := s.get(ctx, "/webhooks", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// RegisterEvent registers a new webhook event.
// POST /webhooks.
func (s *WebhooksService) RegisterEvent(ctx context.Context, req *domain.RegisterWebhookRequest) (*domain.WebhookEvent, error) {
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
	var result domain.APIResponse[domain.WebhookEvent]
	if err := s.post(ctx, "/webhooks", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// DeleteEvent removes a registered webhook event.
// DELETE /webhooks/{eventTypeId}.
func (s *WebhooksService) DeleteEvent(ctx context.Context, eventTypeID string) error {
	if eventTypeID == "" {
		return newValidationErr("eventTypeId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return err
	}
	return s.del(ctx, "/webhooks/"+eventTypeID, h)
}

// ResetSecret resets the webhook signing secret for this application.
// POST /webhooks/secret/reset.
func (s *WebhooksService) ResetSecret(ctx context.Context) (*domain.ResetWebhookSecretResponse, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.ResetWebhookSecretResponse]
	if err := s.post(ctx, "/webhooks/secret/reset", nil, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

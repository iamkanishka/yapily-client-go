package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// BulkPaymentsService handles bulk payment initiation and status checks.
type BulkPaymentsService struct{ base }

// NewBulkPaymentsService creates a new BulkPaymentsService.
func NewBulkPaymentsService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *BulkPaymentsService {
	return &BulkPaymentsService{base{transport: t, auth: a, logger: l}}
}

// Create initiates a bulk payment after obtaining authorisation.
// POST /bulk-payments.
func (s *BulkPaymentsService) Create(ctx context.Context, consentToken string, req *domain.BulkPaymentRequest) (*domain.BulkPayment, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if len(req.Payments) == 0 {
		return nil, newValidationErr("payments", "must contain at least one payment")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.BulkPayment]
	if err := s.post(ctx, "/bulk-payments", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetStatus returns the processing status of a bulk payment.
// GET /bulk-payments/{bulkPaymentId}.
func (s *BulkPaymentsService) GetStatus(ctx context.Context, consentToken, bulkPaymentID string) (*domain.BulkPaymentStatus, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	if bulkPaymentID == "" {
		return nil, newValidationErr("bulkPaymentID", "must not be empty")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.BulkPaymentStatus]
	if err := s.get(ctx, "/bulk-payments/"+bulkPaymentID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

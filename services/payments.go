package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	sdkerrors "github.com/iamkanishka/yapily-client-go/errors"
	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// PaymentsService handles single payment initiation and retrieval.
type PaymentsService struct{ base }

// NewPaymentsService creates a new PaymentsService.
func NewPaymentsService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *PaymentsService {
	return &PaymentsService{base{transport: t, auth: a, logger: l}}
}

// Create initiates a new payment.
// POST /payments.
func (s *PaymentsService) Create(ctx context.Context, consentToken string, req *domain.PaymentRequest) (*domain.Payment, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	if err := validatePaymentRequest(req); err != nil {
		return nil, err
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	if req.IdempotencyKey != "" {
		h["Idempotency-Key"] = req.IdempotencyKey
	}
	var result domain.APIResponse[domain.Payment]
	if err := s.post(ctx, "/payments", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// Get retrieves full payment details.
// GET /payments/{paymentId}.
func (s *PaymentsService) Get(ctx context.Context, consentToken, paymentID string) (*domain.Payment, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	if paymentID == "" {
		return nil, newValidationErr("paymentID", "must not be empty")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Payment]
	if err := s.get(ctx, "/payments/"+paymentID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func validatePaymentRequest(req *domain.PaymentRequest) error {
	if req == nil {
		return &sdkerrors.ValidationError{Field: "request", Message: "must not be nil"}
	}
	if req.Amount <= 0 {
		return &sdkerrors.ValidationError{Field: "amount", Message: "must be greater than 0"}
	}
	if req.Currency == "" {
		return &sdkerrors.ValidationError{Field: "currency", Message: "must not be empty"}
	}
	if req.Recipient.Name == "" {
		return &sdkerrors.ValidationError{Field: "recipient.name", Message: "must not be empty"}
	}
	return nil
}

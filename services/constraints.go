package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// ConstraintsService retrieves payment and data constraints per institution.
type ConstraintsService struct{ base }

// NewConstraintsService creates a new ConstraintsService.
func NewConstraintsService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *ConstraintsService {
	return &ConstraintsService{base{transport: t, auth: a, logger: l}}
}

// GetPaymentConstraints returns payment constraint rules, optionally filtered.
// GET /constraints/payment.
func (s *ConstraintsService) GetPaymentConstraints(ctx context.Context, institutionID, paymentType string) ([]domain.PaymentConstraint, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	q := map[string]string{}
	if institutionID != "" {
		q["institution-id"] = institutionID
	}
	if paymentType != "" {
		q["payment-type"] = paymentType
	}
	var result domain.APIResponse[[]domain.PaymentConstraint]
	if err := s.getWithQuery(ctx, "/constraints/payment", q, h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetDataConstraints returns data constraint rules, optionally filtered by institution.
// GET /constraints/data.
func (s *ConstraintsService) GetDataConstraints(ctx context.Context, institutionID string) ([]domain.DataConstraint, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	q := map[string]string{}
	if institutionID != "" {
		q["institution-id"] = institutionID
	}
	var result domain.APIResponse[[]domain.DataConstraint]
	if err := s.getWithQuery(ctx, "/constraints/data", q, h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

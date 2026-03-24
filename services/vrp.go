package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// VRPService handles Variable Recurring Payment authorisations and payments.
type VRPService struct{ base }

// NewVRPService creates a new VRPService.
func NewVRPService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *VRPService {
	return &VRPService{base{transport: t, auth: a, logger: l}}
}

// CreateSweepingAuthorisation initiates a sweeping VRP consent.
// POST /vrp-consents.
func (s *VRPService) CreateSweepingAuthorisation(ctx context.Context, req *domain.VRPAuthorisationRequest) (*domain.VRPConsent, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.InstitutionID == "" {
		return nil, newValidationErr("institutionId", "must not be empty")
	}
	if req.ApplicationUserID == "" {
		return nil, newValidationErr("applicationUserId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.VRPConsent]
	if err := s.post(ctx, "/vrp-consents", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetSweepingConsentDetails retrieves VRP consent details.
// GET /vrp-consents/{consentId}.
func (s *VRPService) GetSweepingConsentDetails(ctx context.Context, consentID string) (*domain.VRPConsent, error) {
	if consentID == "" {
		return nil, newValidationErr("consentID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.VRPConsent]
	if err := s.get(ctx, "/vrp-consents/"+consentID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreatePayment initiates a Variable Recurring Payment.
// POST /vrp-consents/{consentId}/payments.
func (s *VRPService) CreatePayment(ctx context.Context, consentToken, consentID string, req *domain.VRPPaymentRequest) (*domain.VRPPayment, error) {
	if consentID == "" {
		return nil, newValidationErr("consentID", "must not be empty")
	}
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.Amount <= 0 {
		return nil, newValidationErr("amount", "must be greater than 0")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.VRPPayment]
	if err := s.post(ctx, "/vrp-consents/"+consentID+"/payments", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetPaymentDetails retrieves VRP payment details.
// GET /vrp-consents/{consentId}/payments/{paymentId}.
func (s *VRPService) GetPaymentDetails(ctx context.Context, consentID, paymentID string) (*domain.VRPPayment, error) {
	if consentID == "" {
		return nil, newValidationErr("consentID", "must not be empty")
	}
	if paymentID == "" {
		return nil, newValidationErr("paymentID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.VRPPayment]
	if err := s.get(ctx, "/vrp-consents/"+consentID+"/payments/"+paymentID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// ConfirmFunds checks whether sufficient funds are available for a VRP payment.
// POST /vrp-consents/{consentId}/funds-confirmation.
func (s *VRPService) ConfirmFunds(ctx context.Context, consentToken, consentID string, req *domain.FundsConfirmationRequest) (*domain.FundsConfirmationResponse, error) {
	if consentID == "" {
		return nil, newValidationErr("consentID", "must not be empty")
	}
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.FundsConfirmationResponse]
	if err := s.post(ctx, "/vrp-consents/"+consentID+"/funds-confirmation", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

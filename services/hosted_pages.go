package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// HostedPagesService handles Yapily Hosted Consent and Hosted Payment page flows.
type HostedPagesService struct{ base }

// NewHostedPagesService creates a new HostedPagesService.
func NewHostedPagesService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *HostedPagesService {
	return &HostedPagesService{base{transport: t, auth: a, logger: l}}
}

// ── Hosted Consent Pages ──────────────────────────────────────────────────────

// CreateConsentRequest creates a hosted consent page session.
// POST /hosted/consent/account-auth-requests.
func (s *HostedPagesService) CreateConsentRequest(ctx context.Context, req *domain.HostedConsentRequest) (*domain.HostedConsentResponse, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.ApplicationUserID == "" {
		return nil, newValidationErr("applicationUserId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.HostedConsentResponse]
	if err := s.post(ctx, "/hosted/consent/account-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetConsentRequest retrieves a hosted consent page session.
// GET /hosted/consent/account-auth-requests/{consentRequestId}.
func (s *HostedPagesService) GetConsentRequest(ctx context.Context, consentRequestID string) (*domain.HostedConsentResponse, error) {
	if consentRequestID == "" {
		return nil, newValidationErr("consentRequestId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.HostedConsentResponse]
	if err := s.get(ctx, "/hosted/consent/account-auth-requests/"+consentRequestID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// ── Hosted Payment Pages ──────────────────────────────────────────────────────

// CreatePaymentRequest creates a hosted payment page session.
// POST /hosted/payment/payment-auth-requests.
func (s *HostedPagesService) CreatePaymentRequest(ctx context.Context, req *domain.HostedPaymentRequest) (*domain.HostedPaymentResponse, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.ApplicationUserID == "" {
		return nil, newValidationErr("applicationUserId", "must not be empty")
	}
	if req.PaymentRequest == nil {
		return nil, newValidationErr("paymentRequest", "must not be nil")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.HostedPaymentResponse]
	if err := s.post(ctx, "/hosted/payment/payment-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetPaymentRequest retrieves a hosted payment page session.
// GET /hosted/payment/payment-auth-requests/{paymentRequestId}.
func (s *HostedPagesService) GetPaymentRequest(ctx context.Context, paymentRequestID string) (*domain.HostedPaymentResponse, error) {
	if paymentRequestID == "" {
		return nil, newValidationErr("paymentRequestId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.HostedPaymentResponse]
	if err := s.get(ctx, "/hosted/payment/payment-auth-requests/"+paymentRequestID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreatePayByLink creates a Pay By Link URL.
// POST /hosted/payment/pay-by-link.
func (s *HostedPagesService) CreatePayByLink(ctx context.Context, req *domain.PayByLinkRequest) (*domain.PayByLinkResponse, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.Amount <= 0 {
		return nil, newValidationErr("amount", "must be greater than 0")
	}
	if req.Currency == "" {
		return nil, newValidationErr("currency", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.PayByLinkResponse]
	if err := s.post(ctx, "/hosted/payment/pay-by-link", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// ── Hosted VRP Pages ──────────────────────────────────────────────────────────

// GetVRPConsentRequests lists all hosted VRP consent requests.
// GET /hosted/vrp/vrp-consent-auth-requests.
func (s *HostedPagesService) GetVRPConsentRequests(ctx context.Context) ([]domain.VRPConsent, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[[]domain.VRPConsent]
	if err := s.get(ctx, "/hosted/vrp/vrp-consent-auth-requests", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// CreateVRPConsent creates a hosted VRP consent.
// POST /hosted/vrp/vrp-consent-auth-requests.
func (s *HostedPagesService) CreateVRPConsent(ctx context.Context, req *domain.VRPAuthorisationRequest) (*domain.VRPConsent, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.VRPConsent]
	if err := s.post(ctx, "/hosted/vrp/vrp-consent-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetVRPConsentRequest retrieves a single hosted VRP consent request.
// GET /hosted/vrp/vrp-consent-auth-requests/{vrpConsentRequestId}.
func (s *HostedPagesService) GetVRPConsentRequest(ctx context.Context, vrpConsentRequestID string) (*domain.VRPConsent, error) {
	if vrpConsentRequestID == "" {
		return nil, newValidationErr("vrpConsentRequestId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.VRPConsent]
	if err := s.get(ctx, "/hosted/vrp/vrp-consent-auth-requests/"+vrpConsentRequestID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// RevokeVRPConsentRequest revokes a hosted VRP consent.
// POST /hosted/vrp/vrp-consent-auth-requests/{vrpConsentRequestId}/revoke.
func (s *HostedPagesService) RevokeVRPConsentRequest(ctx context.Context, vrpConsentRequestID string) error {
	if vrpConsentRequestID == "" {
		return newValidationErr("vrpConsentRequestId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return err
	}
	return s.post(ctx, "/hosted/vrp/vrp-consent-auth-requests/"+vrpConsentRequestID+"/revoke", nil, h, nil)
}

// CreateVRPPayment creates a hosted VRP payment.
// POST /hosted/vrp/vrp-consent-auth-requests/{vrpConsentRequestId}/payments.
func (s *HostedPagesService) CreateVRPPayment(ctx context.Context, vrpConsentRequestID string, req *domain.VRPPaymentRequest) (*domain.VRPPayment, error) {
	if vrpConsentRequestID == "" {
		return nil, newValidationErr("vrpConsentRequestId", "must not be empty")
	}
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.VRPPayment]
	if err := s.post(ctx, "/hosted/vrp/vrp-consent-auth-requests/"+vrpConsentRequestID+"/payments", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetVRPPayment retrieves a hosted VRP payment.
// GET /hosted/vrp/vrp-consent-auth-requests/{vrpConsentRequestId}/payments/{paymentId}.
func (s *HostedPagesService) GetVRPPayment(ctx context.Context, vrpConsentRequestID, paymentID string) (*domain.VRPPayment, error) {
	if vrpConsentRequestID == "" {
		return nil, newValidationErr("vrpConsentRequestId", "must not be empty")
	}
	if paymentID == "" {
		return nil, newValidationErr("paymentID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.VRPPayment]
	if err := s.get(ctx, "/hosted/vrp/vrp-consent-auth-requests/"+vrpConsentRequestID+"/payments/"+paymentID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CheckFundsAvailability checks if funds are available for a hosted payment.
// POST /hosted/payment/funds-confirmation.
func (s *HostedPagesService) CheckFundsAvailability(ctx context.Context, consentToken string, req *domain.FundsConfirmationRequest) (*domain.FundsConfirmationResponse, error) {
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
	if err := s.post(ctx, "/hosted/payment/funds-confirmation", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

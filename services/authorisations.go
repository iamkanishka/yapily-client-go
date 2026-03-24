package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// AuthorisationsService handles all Yapily authorisation flows.
type AuthorisationsService struct{ base }

// NewAuthorisationsService creates a new AuthorisationsService.
func NewAuthorisationsService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *AuthorisationsService {
	return &AuthorisationsService{base{transport: t, auth: a, logger: l}}
}

// CreateAccountAuthorisation initiates an account-data consent flow.
// POST /account-auth-requests.
func (s *AuthorisationsService) CreateAccountAuthorisation(ctx context.Context, req *domain.AccountAuthorisationRequest, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.InstitutionID == "" {
		return nil, newValidationErr("institutionId", "must not be empty")
	}
	if req.ApplicationUserID == "" {
		return nil, newValidationErr("applicationUserId", "must not be empty")
	}
	h, err := s.authHeadersWithPSU(ctx, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.post(ctx, "/account-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// ReauthoriseAccountConsent re-authorises an expired account consent.
// PATCH /account-auth-requests.
func (s *AuthorisationsService) ReauthoriseAccountConsent(ctx context.Context, consentToken string, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	h, err := s.authHeadersWithConsentAndPSU(ctx, consentToken, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.patch(ctx, "/account-auth-requests", nil, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreatePaymentAuthorisation initiates a payment consent flow.
// POST /payment-auth-requests.
func (s *AuthorisationsService) CreatePaymentAuthorisation(ctx context.Context, req *domain.PaymentAuthorisationRequest, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.InstitutionID == "" {
		return nil, newValidationErr("institutionId", "must not be empty")
	}
	h, err := s.authHeadersWithPSU(ctx, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.post(ctx, "/payment-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateBulkPaymentAuthorisation initiates a bulk-payment consent flow.
// POST /bulk-payment-auth-requests.
func (s *AuthorisationsService) CreateBulkPaymentAuthorisation(ctx context.Context, req *domain.BulkPaymentAuthorisationRequest, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeadersWithPSU(ctx, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.post(ctx, "/bulk-payment-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreatePreAuthorisation creates a pre-authorisation.
// POST /pre-auth-requests.
func (s *AuthorisationsService) CreatePreAuthorisation(ctx context.Context, req *domain.PreAuthorisationRequest) (*domain.Authorisation, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.post(ctx, "/pre-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateAccountPreAuthorisation updates an existing pre-authorisation.
// PUT /account-auth-requests/{consentToken}.
func (s *AuthorisationsService) UpdateAccountPreAuthorisation(ctx context.Context, consentToken string, req *domain.AccountAuthorisationRequest, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	h, err := s.authHeadersWithConsentAndPSU(ctx, consentToken, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.put(ctx, "/account-auth-requests/"+consentToken, req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateEmbeddedAccountAuthorisation initiates an embedded account auth flow.
// POST /embedded-account-auth-requests.
func (s *AuthorisationsService) CreateEmbeddedAccountAuthorisation(ctx context.Context, req *domain.AccountAuthorisationRequest, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeadersWithPSU(ctx, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.post(ctx, "/embedded-account-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateEmbeddedAccountAuthorisation updates an embedded account auth step.
// PUT /embedded-account-auth-requests/{consentToken}.
func (s *AuthorisationsService) UpdateEmbeddedAccountAuthorisation(ctx context.Context, consentToken string, req *domain.AccountAuthorisationRequest, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	h, err := s.authHeadersWithConsentAndPSU(ctx, consentToken, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.put(ctx, "/embedded-account-auth-requests/"+consentToken, req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateEmbeddedPaymentAuthorisation initiates an embedded payment auth flow.
// POST /embedded-payment-auth-requests.
func (s *AuthorisationsService) CreateEmbeddedPaymentAuthorisation(ctx context.Context, req *domain.PaymentAuthorisationRequest, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeadersWithPSU(ctx, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.post(ctx, "/embedded-payment-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateEmbeddedPaymentAuthorisation updates an embedded payment auth step.
// PUT /embedded-payment-auth-requests/{consentToken}.
func (s *AuthorisationsService) UpdateEmbeddedPaymentAuthorisation(ctx context.Context, consentToken string, req *domain.PaymentAuthorisationRequest, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	h, err := s.authHeadersWithConsentAndPSU(ctx, consentToken, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.put(ctx, "/embedded-payment-auth-requests/"+consentToken, req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateEmbeddedBulkPaymentAuthorisation initiates an embedded bulk payment auth flow.
// POST /embedded-bulk-payment-auth-requests.
func (s *AuthorisationsService) CreateEmbeddedBulkPaymentAuthorisation(ctx context.Context, req *domain.BulkPaymentAuthorisationRequest, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeadersWithPSU(ctx, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.post(ctx, "/embedded-bulk-payment-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateEmbeddedBulkPaymentAuthorisation updates an embedded bulk payment auth step.
// PUT /embedded-bulk-payment-auth-requests/{consentToken}.
func (s *AuthorisationsService) UpdateEmbeddedBulkPaymentAuthorisation(ctx context.Context, consentToken string, req *domain.BulkPaymentAuthorisationRequest, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	h, err := s.authHeadersWithConsentAndPSU(ctx, consentToken, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.put(ctx, "/embedded-bulk-payment-auth-requests/"+consentToken, req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreatePaymentPreAuthorisation creates a payment pre-authorisation.
// POST /payment-pre-auth-requests.
func (s *AuthorisationsService) CreatePaymentPreAuthorisation(ctx context.Context, req *domain.PreAuthorisationRequest) (*domain.Authorisation, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.post(ctx, "/payment-pre-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdatePaymentPreAuthorisation updates an existing payment pre-authorisation.
// PUT /payment-pre-auth-requests.
func (s *AuthorisationsService) UpdatePaymentPreAuthorisation(ctx context.Context, consentToken string, req *domain.PaymentAuthorisationRequest, psu *domain.PSUHeaders) (*domain.Authorisation, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	h, err := s.authHeadersWithConsentAndPSU(ctx, consentToken, psu)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Authorisation]
	if err := s.put(ctx, "/payment-pre-auth-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

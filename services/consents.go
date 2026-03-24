package services

import (
	"context"
	"strconv"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// ConsentsService handles all consent lifecycle operations.
type ConsentsService struct{ base }

// NewConsentsService creates a new ConsentsService.
func NewConsentsService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *ConsentsService {
	return &ConsentsService{base{transport: t, auth: a, logger: l}}
}

// List returns all consents. At least one filter or limit must be provided.
// GET /consents.
func (s *ConsentsService) List(ctx context.Context, params *domain.ConsentListParams) ([]domain.Consent, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	q := map[string]string{}
	if params != nil {
		if len(params.ApplicationUserIDs) > 0 {
			// Yapily accepts repeated query params; encode first one — caller should use
			// raw URL for multi-value if needed.
			q["filter[applicationUserId]"] = params.ApplicationUserIDs[0]
		}
		if len(params.UserUUIDs) > 0 {
			q["filter[userUuid]"] = params.UserUUIDs[0]
		}
		if len(params.InstitutionIDs) > 0 {
			q["filter[institution]"] = params.InstitutionIDs[0]
		}
		if len(params.Statuses) > 0 {
			q["filter[status]"] = params.Statuses[0]
		}
		if params.From != "" {
			q["from"] = params.From
		}
		if params.Before != "" {
			q["before"] = params.Before
		}
		if params.Limit > 0 {
			q["limit"] = strconv.Itoa(params.Limit)
		}
		if params.Offset > 0 {
			q["offset"] = strconv.Itoa(params.Offset)
		}
	}
	var result domain.APIResponse[[]domain.Consent]
	if err := s.getWithQuery(ctx, "/consents", q, h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// Get retrieves a consent by ID.
// GET /consents/{consentId}.
func (s *ConsentsService) Get(ctx context.Context, consentID string) (*domain.Consent, error) {
	if consentID == "" {
		return nil, newValidationErr("consentID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Consent]
	if err := s.get(ctx, "/consents/"+consentID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// Delete revokes and deletes a consent.
// DELETE /consents/{consentId}.
func (s *ConsentsService) Delete(ctx context.Context, consentID string, forceDelete bool) error {
	if consentID == "" {
		return newValidationErr("consentID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return err
	}
	path := "/consents/" + consentID
	if forceDelete {
		path += "?forceDelete=true"
	}
	return s.del(ctx, path, h)
}

// Extend updates the reconfirmation timestamp of a consent.
// POST /consents/{consentId}/extend.
func (s *ConsentsService) Extend(ctx context.Context, consentID string, req *domain.ExtendConsentRequest) (*domain.Consent, error) {
	if consentID == "" {
		return nil, newValidationErr("consentID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Consent]
	if err := s.post(ctx, "/consents/"+consentID+"/extend", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// ExchangeOAuth2Code exchanges an OAuth2 authorisation code for a consent token.
// POST /consent-auth-code  (singular — matches real Yapily API).
func (s *ConsentsService) ExchangeOAuth2Code(ctx context.Context, req *domain.ExchangeCodeRequest) (*domain.Consent, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.Code == "" {
		return nil, newValidationErr("code", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	// Returns Consent directly (not wrapped in APIResponse)
	var result domain.Consent
	if err := s.post(ctx, "/consent-auth-code", req, h, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExchangeOneTimeToken exchanges a one-time token for a consent.
// POST /consent-one-time-token  (singular — matches real Yapily API).
func (s *ConsentsService) ExchangeOneTimeToken(ctx context.Context, req *domain.OneTimeTokenRequest) (*domain.Consent, error) {
	if req == nil || req.OneTimeToken == "" {
		return nil, newValidationErr("oneTimeToken", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	// Returns Consent directly (not wrapped in APIResponse)
	var result domain.Consent
	if err := s.post(ctx, "/consent-one-time-token", req, h, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// helper to satisfy compiler

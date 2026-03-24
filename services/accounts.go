package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// AccountsService handles account-related API calls.
type AccountsService struct{ base }

// NewAccountsService creates a new AccountsService.
func NewAccountsService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *AccountsService {
	return &AccountsService{base{transport: t, auth: a, logger: l}}
}

// List returns all accounts accessible under the provided consent token.
// GET /accounts.
func (s *AccountsService) List(ctx context.Context, consentToken string) ([]domain.Account, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[[]domain.Account]
	if err := s.get(ctx, "/accounts", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// Get returns a single account by ID.
// GET /accounts/{accountId}.
func (s *AccountsService) Get(ctx context.Context, consentToken, accountID string) (*domain.Account, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	if accountID == "" {
		return nil, newValidationErr("accountID", "must not be empty")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Account]
	if err := s.get(ctx, "/accounts/"+accountID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

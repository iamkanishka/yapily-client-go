package services

import (
	"context"
	"strconv"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// TransactionsService handles transaction-related API calls.
type TransactionsService struct{ base }

// NewTransactionsService creates a new TransactionsService.
func NewTransactionsService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *TransactionsService {
	return &TransactionsService{base{transport: t, auth: a, logger: l}}
}

// List returns transactions for the given account with optional filters.
// GET /accounts/{accountId}/transactions.
func (s *TransactionsService) List(ctx context.Context, consentToken, accountID string, params *domain.TransactionQueryParams) ([]domain.Transaction, error) {
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
	q := map[string]string{}
	if params != nil {
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
	var result domain.APIResponse[[]domain.Transaction]
	if err := s.getWithQuery(ctx, "/accounts/"+accountID+"/transactions", q, h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ListPaginated iterates all transaction pages via a callback.
// Return false from fn to stop early.
func (s *TransactionsService) ListPaginated(ctx context.Context, consentToken, accountID string, pageSize int, fn func([]domain.Transaction) bool) error {
	offset := 0
	for {
		params := &domain.TransactionQueryParams{
			PaginationParams: domain.PaginationParams{Limit: pageSize, Offset: offset},
		}
		txns, err := s.List(ctx, consentToken, accountID, params)
		if err != nil {
			return err
		}
		if len(txns) == 0 {
			return nil
		}
		if !fn(txns) {
			return nil
		}
		if len(txns) < pageSize {
			return nil
		}
		offset += pageSize
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

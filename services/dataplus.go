package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// DataPlusService handles transaction enrichment and categorisation.
type DataPlusService struct{ base }

// NewDataPlusService creates a new DataPlusService.
func NewDataPlusService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *DataPlusService {
	return &DataPlusService{base{transport: t, auth: a, logger: l}}
}

// Enrich submits raw transactions for categorisation and enrichment.
// POST /enrich.
func (s *DataPlusService) Enrich(ctx context.Context, req *domain.EnrichmentRequest) (*domain.EnrichmentResult, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if len(req.Transactions) == 0 {
		return nil, newValidationErr("transactions", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.EnrichmentResult]
	if err := s.post(ctx, "/enrich", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetEnrichmentResults fetches the results of a previously submitted enrichment job.
// GET /enrich/{jobId}.
func (s *DataPlusService) GetEnrichmentResults(ctx context.Context, jobID string) (*domain.EnrichmentResult, error) {
	if jobID == "" {
		return nil, newValidationErr("jobID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.EnrichmentResult]
	if err := s.get(ctx, "/enrich/"+jobID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetEnrichmentLabels returns available categorisation labels.
// GET /enrich/labels.
func (s *DataPlusService) GetEnrichmentLabels(ctx context.Context) ([]string, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[[]string]
	if err := s.get(ctx, "/enrich/labels", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// EnrichAccountTransactions submits an account's transactions for enrichment.
// and returns the enriched result synchronously.
// POST /accounts/{accountId}/transactions/categorisation.
func (s *DataPlusService) EnrichAccountTransactions(ctx context.Context, consentToken, accountID string, req *domain.EnrichmentRequest) (*domain.EnrichmentResult, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	if accountID == "" {
		return nil, newValidationErr("accountID", "must not be empty")
	}
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.EnrichmentResult]
	if err := s.post(ctx, "/accounts/"+accountID+"/transactions/categorisation", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// TransactionsAndEnrichment submits account transactions for enrichment via the new endpoint.
// POST /enrich-requests  (newer endpoint replacing POST /enrich for account-level enrichment).
func (s *DataPlusService) TransactionsAndEnrichment(ctx context.Context, consentToken, accountID string, req *domain.EnrichmentRequest) (*domain.EnrichmentResult, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	if accountID == "" {
		return nil, newValidationErr("accountID", "must not be empty")
	}
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.EnrichmentResult]
	if err := s.post(ctx, "/enrich-requests", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetEnrichmentRequestResults fetches results from the new enrich-requests endpoint.
// GET /enrich-requests/{jobId}.
func (s *DataPlusService) GetEnrichmentRequestResults(ctx context.Context, jobID string) (*domain.EnrichmentResult, error) {
	if jobID == "" {
		return nil, newValidationErr("jobID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.EnrichmentResult]
	if err := s.get(ctx, "/enrich-requests/"+jobID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

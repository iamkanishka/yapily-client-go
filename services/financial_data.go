package services

import (
	"context"
	"strconv"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// FinancialDataService handles all financial-data API calls beyond basic accounts/transactions.
type FinancialDataService struct{ base }

// NewFinancialDataService creates a new FinancialDataService.
func NewFinancialDataService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *FinancialDataService {
	return &FinancialDataService{base{transport: t, auth: a, logger: l}}
}

// GetAccountBalances retrieves detailed balance info for an account.
// GET /accounts/{accountId}/balances.
func (s *FinancialDataService) GetAccountBalances(ctx context.Context, consentToken, accountID string) (*domain.AccountBalance, error) {
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
	var result domain.APIResponse[domain.AccountBalance]
	if err := s.get(ctx, "/accounts/"+accountID+"/balances", h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetAccountBeneficiaries returns the beneficiaries linked to an account.
// GET /accounts/{accountId}/beneficiaries.
func (s *FinancialDataService) GetAccountBeneficiaries(ctx context.Context, consentToken, accountID string) ([]domain.Beneficiary, error) {
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
	var result domain.APIResponse[[]domain.Beneficiary]
	if err := s.get(ctx, "/accounts/"+accountID+"/beneficiaries", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetAccountDirectDebits returns active direct debit mandates.
// GET /accounts/{accountId}/direct-debits.
func (s *FinancialDataService) GetAccountDirectDebits(ctx context.Context, consentToken, accountID string) ([]domain.DirectDebit, error) {
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
	var result domain.APIResponse[[]domain.DirectDebit]
	if err := s.get(ctx, "/accounts/"+accountID+"/direct-debits", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetAccountScheduledPayments returns future scheduled payments.
// GET /accounts/{accountId}/scheduled-payments.
func (s *FinancialDataService) GetAccountScheduledPayments(ctx context.Context, consentToken, accountID string) ([]domain.ScheduledPayment, error) {
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
	var result domain.APIResponse[[]domain.ScheduledPayment]
	if err := s.get(ctx, "/accounts/"+accountID+"/scheduled-payments", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetAccountPeriodicPayments returns standing orders / recurring payments.
// GET /accounts/{accountId}/periodic-payments.
func (s *FinancialDataService) GetAccountPeriodicPayments(ctx context.Context, consentToken, accountID string) ([]domain.PeriodicPaymentResponse, error) {
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
	var result domain.APIResponse[[]domain.PeriodicPaymentResponse]
	if err := s.get(ctx, "/accounts/"+accountID+"/periodic-payments", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetAccountStatements returns a list of statements for an account.
// GET /accounts/{accountId}/statements.
func (s *FinancialDataService) GetAccountStatements(ctx context.Context, consentToken, accountID string, params *domain.PaginationParams) ([]domain.Statement, error) {
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
	q := buildPaginationQuery(params)
	var result domain.APIResponse[[]domain.Statement]
	if err := s.getWithQuery(ctx, "/accounts/"+accountID+"/statements", q, h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetAccountStatement returns a single statement.
// GET /accounts/{accountId}/statements/{statementId}.
func (s *FinancialDataService) GetAccountStatement(ctx context.Context, consentToken, accountID, statementID string) (*domain.Statement, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	if accountID == "" {
		return nil, newValidationErr("accountID", "must not be empty")
	}
	if statementID == "" {
		return nil, newValidationErr("statementID", "must not be empty")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Statement]
	if err := s.get(ctx, "/accounts/"+accountID+"/statements/"+statementID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetAccountStatementFile returns a PDF/CSV statement file as raw bytes.
// GET /accounts/{accountId}/statements/{statementId}/file.
func (s *FinancialDataService) GetAccountStatementFile(ctx context.Context, consentToken, accountID, statementID string) ([]byte, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	if accountID == "" {
		return nil, newValidationErr("accountID", "must not be empty")
	}
	if statementID == "" {
		return nil, newValidationErr("statementID", "must not be empty")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	// File responses are raw bytes, not JSON — use a []byte out.
	var raw []byte
	if err := s.get(ctx, "/accounts/"+accountID+"/statements/"+statementID+"/file", h, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// GetIdentity returns identity information for the authorised user.
// GET /identity.
func (s *FinancialDataService) GetIdentity(ctx context.Context, consentToken string) (*domain.Identity, error) {
	if consentToken == "" {
		return nil, newValidationErr("consentToken", "must not be empty")
	}
	h, err := s.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Identity]
	if err := s.get(ctx, "/identity", h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetRealTimeAccountTransactions fetches real-time (live) transactions.
// GET /accounts/{accountId}/transactions/real-time.
func (s *FinancialDataService) GetRealTimeAccountTransactions(ctx context.Context, consentToken, accountID string) ([]domain.Transaction, error) {
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
	var result domain.APIResponse[[]domain.Transaction]
	if err := s.get(ctx, "/accounts/"+accountID+"/transactions/real-time", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// buildPaginationQuery converts PaginationParams to query string map.
func buildPaginationQuery(p *domain.PaginationParams) map[string]string {
	q := map[string]string{}
	if p == nil {
		return q
	}
	if p.Limit > 0 {
		q["limit"] = strconv.Itoa(p.Limit)
	}
	if p.Offset > 0 {
		q["offset"] = strconv.Itoa(p.Offset)
	}
	return q
}

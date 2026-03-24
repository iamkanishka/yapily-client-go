package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// InstitutionsService handles institution-related API calls.
type InstitutionsService struct{ base }

// NewInstitutionsService creates a new InstitutionsService.
func NewInstitutionsService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *InstitutionsService {
	return &InstitutionsService{base{transport: t, auth: a, logger: l}}
}

// List returns all available institutions.
// GET /institutions.
func (s *InstitutionsService) List(ctx context.Context) ([]domain.Institution, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[[]domain.Institution]
	if err := s.get(ctx, "/institutions", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// Get returns a single institution by ID.
// GET /institutions/{institutionId}.
func (s *InstitutionsService) Get(ctx context.Context, institutionID string) (*domain.Institution, error) {
	if institutionID == "" {
		return nil, newValidationErr("institutionID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Institution]
	if err := s.get(ctx, "/institutions/"+institutionID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

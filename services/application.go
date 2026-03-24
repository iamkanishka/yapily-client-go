package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// ApplicationService manages Yapily applications and sub-applications.
type ApplicationService struct{ base }

// NewApplicationService creates a new ApplicationService.
func NewApplicationService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *ApplicationService {
	return &ApplicationService{base{transport: t, auth: a, logger: l}}
}

// GetDetails retrieves the current application's details.
// GET /application.
func (s *ApplicationService) GetDetails(ctx context.Context) (*domain.Application, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Application]
	if err := s.get(ctx, "/application", h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// Update updates the current application.
// PUT /application.
func (s *ApplicationService) Update(ctx context.Context, req *domain.ApplicationRequest) (*domain.Application, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Application]
	if err := s.put(ctx, "/application", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// Delete deletes the current application.
// DELETE /application.
func (s *ApplicationService) Delete(ctx context.Context) error {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return err
	}
	return s.del(ctx, "/application", h)
}

// ListSubApplications returns all sub-applications for the root application.
// GET /application/sub-applications.
func (s *ApplicationService) ListSubApplications(ctx context.Context) ([]domain.Application, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[[]domain.Application]
	if err := s.get(ctx, "/applications", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// CreateSubApplication creates a sub-application under the root.
// POST /application/sub-applications.
func (s *ApplicationService) CreateSubApplication(ctx context.Context, req *domain.ApplicationRequest) (*domain.Application, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.Name == "" {
		return nil, newValidationErr("name", "must not be empty")
	}
	if req.MerchantCategoryCode == "" {
		return nil, newValidationErr("merchantCategoryCode", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.Application]
	if err := s.post(ctx, "/applications", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetVRPConfiguration retrieves the VRP configuration for this application.
// GET /application/vrp-configuration.
func (s *ApplicationService) GetVRPConfiguration(ctx context.Context) (*domain.VRPConfiguration, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.VRPConfiguration]
	if err := s.get(ctx, "/application/vrp-configuration", h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// CreateVRPConfiguration creates a VRP configuration for this application.
// POST /application/vrp-configuration.
func (s *ApplicationService) CreateVRPConfiguration(ctx context.Context, req *domain.VRPConfiguration) (*domain.VRPConfiguration, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.VRPConfiguration]
	if err := s.post(ctx, "/application/vrp-configuration", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdateVRPConfiguration updates the VRP configuration for this application.
// PUT /application/vrp-configuration.
func (s *ApplicationService) UpdateVRPConfiguration(ctx context.Context, req *domain.VRPConfiguration) (*domain.VRPConfiguration, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.VRPConfiguration]
	if err := s.put(ctx, "/application/vrp-configuration", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

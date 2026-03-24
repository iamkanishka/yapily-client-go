package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// BeneficiariesService manages both application-level and user-level beneficiaries.
type BeneficiariesService struct{ base }

// NewBeneficiariesService creates a new BeneficiariesService.
func NewBeneficiariesService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *BeneficiariesService {
	return &BeneficiariesService{base{transport: t, auth: a, logger: l}}
}

// ── Application Beneficiaries ─────────────────────────────────────────────────

// ListApplicationBeneficiaries returns all beneficiaries at application level.
// GET /application/beneficiaries.
func (s *BeneficiariesService) ListApplicationBeneficiaries(ctx context.Context) ([]domain.ApplicationBeneficiary, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[[]domain.ApplicationBeneficiary]
	if err := s.get(ctx, "/application/beneficiaries", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// CreateApplicationBeneficiary creates a new application-level beneficiary.
// POST /application/beneficiaries.
func (s *BeneficiariesService) CreateApplicationBeneficiary(ctx context.Context, req *domain.CreateBeneficiaryRequest) (*domain.ApplicationBeneficiary, error) {
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.Name == "" {
		return nil, newValidationErr("name", "must not be empty")
	}
	if len(req.AccountIdentifications) == 0 {
		return nil, newValidationErr("accountIdentifications", "must contain at least one entry")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.ApplicationBeneficiary]
	if err := s.post(ctx, "/application/beneficiaries", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetApplicationBeneficiary retrieves a single application beneficiary by ID.
// GET /application/beneficiaries/{beneficiaryId}.
func (s *BeneficiariesService) GetApplicationBeneficiary(ctx context.Context, beneficiaryID string) (*domain.ApplicationBeneficiary, error) {
	if beneficiaryID == "" {
		return nil, newValidationErr("beneficiaryId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.ApplicationBeneficiary]
	if err := s.get(ctx, "/application/beneficiaries/"+beneficiaryID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// DeleteApplicationBeneficiary removes an application-level beneficiary.
// DELETE /application/beneficiaries/{beneficiaryId}.
func (s *BeneficiariesService) DeleteApplicationBeneficiary(ctx context.Context, beneficiaryID string) error {
	if beneficiaryID == "" {
		return newValidationErr("beneficiaryId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return err
	}
	return s.del(ctx, "/application/beneficiaries/"+beneficiaryID, h)
}

// ── User Beneficiaries ────────────────────────────────────────────────────────

// ListUserBeneficiaries returns all beneficiaries for a specific user.
// GET /users/{userUuid}/beneficiaries.
func (s *BeneficiariesService) ListUserBeneficiaries(ctx context.Context, userUUID string) ([]domain.UserBeneficiary, error) {
	if userUUID == "" {
		return nil, newValidationErr("userUUID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[[]domain.UserBeneficiary]
	if err := s.get(ctx, "/users/"+userUUID+"/beneficiaries", h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// CreateUserBeneficiary creates a new user-level beneficiary.
// POST /users/{userUuid}/beneficiaries.
func (s *BeneficiariesService) CreateUserBeneficiary(ctx context.Context, userUUID string, req *domain.CreateBeneficiaryRequest) (*domain.UserBeneficiary, error) {
	if userUUID == "" {
		return nil, newValidationErr("userUUID", "must not be empty")
	}
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	if req.Name == "" {
		return nil, newValidationErr("name", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.UserBeneficiary]
	if err := s.post(ctx, "/users/"+userUUID+"/beneficiaries", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetUserBeneficiary retrieves a single user beneficiary.
// GET /users/{userUuid}/beneficiaries/{beneficiaryId}.
func (s *BeneficiariesService) GetUserBeneficiary(ctx context.Context, userUUID, beneficiaryID string) (*domain.UserBeneficiary, error) {
	if userUUID == "" {
		return nil, newValidationErr("userUUID", "must not be empty")
	}
	if beneficiaryID == "" {
		return nil, newValidationErr("beneficiaryId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.UserBeneficiary]
	if err := s.get(ctx, "/users/"+userUUID+"/beneficiaries/"+beneficiaryID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// DeleteUserBeneficiary removes a user beneficiary.
// DELETE /users/{userUuid}/beneficiaries/{beneficiaryId}.
func (s *BeneficiariesService) DeleteUserBeneficiary(ctx context.Context, userUUID, beneficiaryID string) error {
	if userUUID == "" {
		return newValidationErr("userUUID", "must not be empty")
	}
	if beneficiaryID == "" {
		return newValidationErr("beneficiaryId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return err
	}
	return s.del(ctx, "/users/"+userUUID+"/beneficiaries/"+beneficiaryID, h)
}

// PatchUserBeneficiary partially updates a user beneficiary.
// PATCH /users/{userUuid}/beneficiaries/{beneficiaryId}.
func (s *BeneficiariesService) PatchUserBeneficiary(ctx context.Context, userUUID, beneficiaryID string, req *domain.PatchUserBeneficiaryRequest) (*domain.UserBeneficiary, error) {
	if userUUID == "" {
		return nil, newValidationErr("userUUID", "must not be empty")
	}
	if beneficiaryID == "" {
		return nil, newValidationErr("beneficiaryId", "must not be empty")
	}
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.UserBeneficiary]
	if err := s.patch(ctx, "/users/"+userUUID+"/beneficiaries/"+beneficiaryID, req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// ApproveBeneficiary approves a pending user beneficiary.
// POST /users/{userUuid}/beneficiaries/{beneficiaryId}/approve.
func (s *BeneficiariesService) ApproveBeneficiary(ctx context.Context, userUUID, beneficiaryID string) (*domain.UserBeneficiary, error) {
	if userUUID == "" {
		return nil, newValidationErr("userUUID", "must not be empty")
	}
	if beneficiaryID == "" {
		return nil, newValidationErr("beneficiaryId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.UserBeneficiary]
	if err := s.post(ctx, "/users/"+userUUID+"/beneficiaries/"+beneficiaryID+"/approve", nil, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// RejectBeneficiary rejects a pending user beneficiary.
// POST /users/{userUuid}/beneficiaries/{beneficiaryId}/reject.
func (s *BeneficiariesService) RejectBeneficiary(ctx context.Context, userUUID, beneficiaryID string) (*domain.UserBeneficiary, error) {
	if userUUID == "" {
		return nil, newValidationErr("userUUID", "must not be empty")
	}
	if beneficiaryID == "" {
		return nil, newValidationErr("beneficiaryId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.UserBeneficiary]
	if err := s.post(ctx, "/users/"+userUUID+"/beneficiaries/"+beneficiaryID+"/reject", nil, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

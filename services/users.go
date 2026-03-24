package services

import (
	"context"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// UsersService manages Yapily application users (PSUs).
type UsersService struct{ base }

// NewUsersService creates a new UsersService.
func NewUsersService(t *transporthttp.Transport, a auth.Provider, l *zap.Logger) *UsersService {
	return &UsersService{base{transport: t, auth: a, logger: l}}
}

// List returns all users for this application.
// GET /users.
func (s *UsersService) List(ctx context.Context, applicationUserID string) ([]domain.User, error) {
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	q := map[string]string{}
	if applicationUserID != "" {
		q["filter[applicationUserId]"] = applicationUserID
	}
	var result domain.APIResponse[[]domain.User]
	if err := s.getWithQuery(ctx, "/users", q, h, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// Create creates a new application user.
// POST /users.
func (s *UsersService) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	if req == nil || req.ApplicationUserID == "" {
		return nil, newValidationErr("applicationUserId", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.User]
	if err := s.post(ctx, "/users", req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// Get returns a single user by UUID.
// GET /users/{userUuid}.
func (s *UsersService) Get(ctx context.Context, userUUID string) (*domain.User, error) {
	if userUUID == "" {
		return nil, newValidationErr("userUUID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.User]
	if err := s.get(ctx, "/users/"+userUUID, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// Delete deletes a user by UUID.
// DELETE /users/{userUuid}.
func (s *UsersService) Delete(ctx context.Context, userUUID string) error {
	if userUUID == "" {
		return newValidationErr("userUUID", "must not be empty")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return err
	}
	return s.del(ctx, "/users/"+userUUID, h)
}

// Update updates a user's applicationUserId.
// PATCH /users/{userUuid}.
func (s *UsersService) Update(ctx context.Context, userUUID string, req *domain.UpdateUserRequest) (*domain.User, error) {
	if userUUID == "" {
		return nil, newValidationErr("userUUID", "must not be empty")
	}
	if req == nil {
		return nil, newValidationErr("request", "must not be nil")
	}
	h, err := s.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	var result domain.APIResponse[domain.User]
	if err := s.patch(ctx, "/users/"+userUUID, req, h, &result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

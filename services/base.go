// Package services contains all API service implementations.
package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"

	"go.uber.org/zap"

	sdkerrors "github.com/iamkanishka/yapily-client-go/errors"
	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

// base is embedded in every service to provide transport, auth and logging.
type base struct {
	transport *transporthttp.Transport
	auth      auth.Provider
	logger    *zap.Logger
}

// authHeaders builds the Authorization header map.
func (b *base) authHeaders(ctx context.Context) (map[string]string, error) {
	tok, err := b.auth.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get auth token: %w", err)
	}
	return map[string]string{"Authorization": tok}, nil
}

// authHeadersWithConsent adds the Consent token header.
func (b *base) authHeadersWithConsent(ctx context.Context, consentToken string) (map[string]string, error) {
	h, err := b.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	if consentToken != "" {
		h["Consent"] = consentToken
	}
	return h, nil
}

// authHeadersWithPSU adds PSU identifiers required by some institutions.
func (b *base) authHeadersWithPSU(ctx context.Context, psu *domain.PSUHeaders) (map[string]string, error) {
	h, err := b.authHeaders(ctx)
	if err != nil {
		return nil, err
	}
	if psu != nil {
		if psu.PSUID != "" {
			h["psu-id"] = psu.PSUID
		}
		if psu.PSUCorporateID != "" {
			h["psu-corporate-id"] = psu.PSUCorporateID
		}
		if psu.PSUIPAddress != "" {
			h["psu-ip-address"] = psu.PSUIPAddress
		}
	}
	return h, nil
}

// authHeadersWithConsentAndPSU combines consent + PSU headers.
func (b *base) authHeadersWithConsentAndPSU(ctx context.Context, consentToken string, psu *domain.PSUHeaders) (map[string]string, error) {
	h, err := b.authHeadersWithConsent(ctx, consentToken)
	if err != nil {
		return nil, err
	}
	if psu != nil {
		if psu.PSUID != "" {
			h["psu-id"] = psu.PSUID
		}
		if psu.PSUCorporateID != "" {
			h["psu-corporate-id"] = psu.PSUCorporateID
		}
		if psu.PSUIPAddress != "" {
			h["psu-ip-address"] = psu.PSUIPAddress
		}
	}
	return h, nil
}

// get is a convenience wrapper for GET requests.
func (b *base) get(ctx context.Context, path string, headers map[string]string, out interface{}) error {
	return b.transport.Request(ctx, http.MethodGet, path, nil, headers, out)
}

// post is a convenience wrapper for POST requests.
func (b *base) post(ctx context.Context, path string, body interface{}, headers map[string]string, out interface{}) error {
	return b.transport.Request(ctx, http.MethodPost, path, body, headers, out)
}

// put is a convenience wrapper for PUT requests.
func (b *base) put(ctx context.Context, path string, body interface{}, headers map[string]string, out interface{}) error {
	return b.transport.Request(ctx, http.MethodPut, path, body, headers, out)
}

// patch is a convenience wrapper for PATCH requests.
func (b *base) patch(ctx context.Context, path string, body interface{}, headers map[string]string, out interface{}) error {
	return b.transport.Request(ctx, http.MethodPatch, path, body, headers, out)
}

// del is a convenience wrapper for DELETE requests.
func (b *base) del(ctx context.Context, path string, headers map[string]string) error {
	return b.transport.Request(ctx, http.MethodDelete, path, nil, headers, nil)
}

// getWithQuery is a GET with query parameters.
func (b *base) getWithQuery(ctx context.Context, path string, params map[string]string, headers map[string]string, out interface{}) error {
	return b.transport.RequestWithQuery(ctx, path, params, headers, out)
}

// newValidationErr creates a ValidationError.
func newValidationErr(field, message string) *sdkerrors.ValidationError {
	return sdkerrors.NewValidationError(field, message)
}

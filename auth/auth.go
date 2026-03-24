// Package auth provides HTTP Basic Authentication for the Yapily API.
// Yapily authenticates using Application Key + Application Secret via HTTP Basic Auth.
package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/iamkanishka/yapily-client-go/domain"
	sdkerrors "github.com/iamkanishka/yapily-client-go/errors"
)

// Provider defines the interface for authentication providers.
type Provider interface {
	// GetToken returns the Authorization header value (e.g. "Basic xxx" or "Bearer yyy").
	GetToken(ctx context.Context) (string, error)
	// Invalidate clears any cached credential state.
	Invalidate()
}

// Config holds credentials for an auth provider.
type Config struct {
	ApplicationKey    string
	ApplicationSecret string
	// TokenURL is only used by the caching OAuth2 provider.
	TokenURL   string
	HTTPClient *http.Client
}

// ── Basic Auth provider (primary Yapily auth mechanism) ───────────────────────

// basicAuthProvider implements Provider using HTTP Basic Authentication.
type basicAuthProvider struct {
	token string
}

// NewBasicAuthProvider creates a Provider using Yapily's Basic Auth scheme.
// Header value = "Basic " + Base64(ApplicationKey:ApplicationSecret).
func NewBasicAuthProvider(cfg Config) Provider {
	raw := cfg.ApplicationKey + ":" + cfg.ApplicationSecret
	return &basicAuthProvider{
		token: "Basic " + base64.StdEncoding.EncodeToString([]byte(raw)),
	}
}

// NewOAuth2Provider is a backwards-compatible alias for NewBasicAuthProvider.
func NewOAuth2Provider(cfg Config) Provider {
	return NewBasicAuthProvider(cfg)
}

// GetToken returns the pre-computed Basic auth header value.
func (p *basicAuthProvider) GetToken(_ context.Context) (string, error) {
	if p.token == "" {
		return "", sdkerrors.NewAuthError("credentials not configured", nil)
	}
	return p.token, nil
}

// Invalidate is a no-op for Basic auth (static credentials don't expire).
func (p *basicAuthProvider) Invalidate() {}

// ── Caching OAuth2 provider (for Bearer token flows) ──────────────────────────

// oauth2CachingProvider fetches and caches Bearer tokens with thread-safe refresh.
type oauth2CachingProvider struct {
	config Config
	mu     sync.RWMutex
	token  *domain.Token
}

// NewCachingOAuth2Provider returns a Provider that fetches and caches Bearer tokens.
func NewCachingOAuth2Provider(cfg Config) Provider {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &oauth2CachingProvider{config: cfg}
}

// GetToken returns a valid bearer token, refreshing if expired.
func (p *oauth2CachingProvider) GetToken(ctx context.Context) (string, error) {
	p.mu.RLock()
	if p.token != nil && !p.token.IsExpired() {
		tok := p.token.AccessToken
		p.mu.RUnlock()
		return "Bearer " + tok, nil
	}
	p.mu.RUnlock()

	// Double-checked locking.
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.token != nil && !p.token.IsExpired() {
		return "Bearer " + p.token.AccessToken, nil
	}

	tok, err := p.fetchToken(ctx)
	if err != nil {
		return "", err
	}
	p.token = tok
	return "Bearer " + p.token.AccessToken, nil
}

// Invalidate clears the cached token.
func (p *oauth2CachingProvider) Invalidate() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.token = nil
}

func (p *oauth2CachingProvider) fetchToken(ctx context.Context) (*domain.Token, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.config.TokenURL, nil)
	if err != nil {
		return nil, sdkerrors.NewAuthError("failed to create token request", err)
	}
	req.SetBasicAuth(p.config.ApplicationKey, p.config.ApplicationSecret)
	req.Header.Set("Accept", "application/json")

	resp, err := p.config.HTTPClient.Do(req)
	if err != nil {
		return nil, sdkerrors.NewAuthError("token request failed", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.NewAuthError("failed to read token response", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &sdkerrors.APIError{
			StatusCode: resp.StatusCode,
			Code:       "AUTH_FAILED",
			Message:    fmt.Sprintf("token endpoint returned %d: %s", resp.StatusCode, string(body)),
		}
	}

	var tok domain.Token
	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, sdkerrors.NewAuthError("failed to decode token response", err)
	}
	tok.ExpiresAt = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	return &tok, nil
}

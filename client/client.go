// Package client is the main entry point for the Yapily Go SDK.
package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/iamkanishka/yapily-client-go/auth"
	"github.com/iamkanishka/yapily-client-go/domain"
	"github.com/iamkanishka/yapily-client-go/services"
	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

const defaultBaseURL = "https://api.yapily.com"

// Client is the fully-wired Yapily SDK client.
// Each field corresponds to a Yapily API group.
type Client struct {
	config    Config
	transport *transporthttp.Transport
	auth      auth.Provider
	logger    *zap.Logger

	// Core
	Institutions *services.InstitutionsService
	Accounts     *services.AccountsService
	Transactions *services.TransactionsService
	Payments     *services.PaymentsService
	BulkPayments *services.BulkPaymentsService
	Consents     *services.ConsentsService

	// Authorisations (all flows)
	Authorisations *services.AuthorisationsService

	// Financial data extras
	FinancialData *services.FinancialDataService

	// Users & identity
	Users *services.UsersService

	// Variable Recurring Payments
	VRP *services.VRPService

	// Notifications / event subscriptions
	Notifications *services.NotificationsService

	// Data Plus (enrichment)
	DataPlus *services.DataPlusService

	// Hosted Pages (consent + payment + VRP)
	HostedPages *services.HostedPagesService

	// Constraints
	Constraints *services.ConstraintsService

	// Application management
	Application *services.ApplicationService

	// Webhooks
	Webhooks *services.WebhooksService

	// Beneficiaries (application + user)
	Beneficiaries *services.BeneficiariesService
}

// Config holds all SDK client configuration.
type Config struct {
	ApplicationKey    string
	ApplicationSecret string
	Environment       domain.Environment
	BaseURL           string
	Timeout           time.Duration
	Logger            *zap.Logger
	RetryConfig       *transporthttp.RetryConfig
	RateLimit         *RateLimitConfig
}

// RateLimitConfig configures the token bucket rate limiter.
type RateLimitConfig struct {
	RequestsPerSecond float64
	Burst             int
}

// Option is a functional option for the client.
type Option func(*Config)

// WithEnvironment sets the target environment.
func WithEnvironment(env domain.Environment) Option {
	return func(c *Config) { c.Environment = env }
}

// WithBaseURL overrides the base URL (useful for testing or private deployments).
func WithBaseURL(u string) Option {
	return func(c *Config) { c.BaseURL = u }
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) { c.Timeout = d }
}

// WithLogger injects a custom zap logger.
func WithLogger(l *zap.Logger) Option {
	return func(c *Config) { c.Logger = l }
}

// WithRetry configures the retry policy.
func WithRetry(cfg transporthttp.RetryConfig) Option {
	return func(c *Config) { c.RetryConfig = &cfg }
}

// WithRateLimit configures rate limiting.
func WithRateLimit(rps float64, burst int) Option {
	return func(c *Config) { c.RateLimit = &RateLimitConfig{rps, burst} }
}

// New creates a fully initialised Yapily SDK client.
//
// applicationKey and applicationSecret are your Yapily console credentials.
//
// Example:
//
//	c, err := client.New("app-key", "app-secret",
//	    client.WithEnvironment(domain.Sandbox),
//	    client.WithTimeout(30*time.Second),
//	)
func New(applicationKey, applicationSecret string, opts ...Option) (*Client, error) {
	cfg := Config{
		ApplicationKey:    applicationKey,
		ApplicationSecret: applicationSecret,
		Environment:       domain.Sandbox,
		Timeout:           30 * time.Second,
	}
	for _, o := range opts {
		o(&cfg)
	}

	if cfg.ApplicationKey == "" {
		return nil, fmt.Errorf("applicationKey is required")
	}
	if cfg.ApplicationSecret == "" {
		return nil, fmt.Errorf("applicationSecret is required")
	}

	logger := cfg.Logger
	if logger == nil {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	// Yapily uses HTTP Basic Auth: Base64(ApplicationKey:ApplicationSecret).
	authProvider := auth.NewBasicAuthProvider(auth.Config{
		ApplicationKey:    cfg.ApplicationKey,
		ApplicationSecret: cfg.ApplicationSecret,
		HTTPClient:        &http.Client{Timeout: 10 * time.Second},
	})

	transportOpts := []transporthttp.TransportOption{
		transporthttp.WithLogger(logger),
	}
	if cfg.Timeout > 0 {
		transportOpts = append(transportOpts, transporthttp.WithHTTPClient(
			&http.Client{Timeout: cfg.Timeout},
		))
	}
	if cfg.RetryConfig != nil {
		transportOpts = append(transportOpts, transporthttp.WithRetryConfig(*cfg.RetryConfig))
	}
	if cfg.RateLimit != nil {
		transportOpts = append(transportOpts, transporthttp.WithRateLimit(
			cfg.RateLimit.RequestsPerSecond, cfg.RateLimit.Burst,
		))
	}

	t := transporthttp.NewTransport(baseURL, transportOpts...)

	c := &Client{
		config:    cfg,
		transport: t,
		auth:      authProvider,
		logger:    logger,
	}

	c.Institutions = services.NewInstitutionsService(t, authProvider, logger)
	c.Accounts = services.NewAccountsService(t, authProvider, logger)
	c.Transactions = services.NewTransactionsService(t, authProvider, logger)
	c.Payments = services.NewPaymentsService(t, authProvider, logger)
	c.BulkPayments = services.NewBulkPaymentsService(t, authProvider, logger)
	c.Consents = services.NewConsentsService(t, authProvider, logger)
	c.Authorisations = services.NewAuthorisationsService(t, authProvider, logger)
	c.FinancialData = services.NewFinancialDataService(t, authProvider, logger)
	c.Users = services.NewUsersService(t, authProvider, logger)
	c.VRP = services.NewVRPService(t, authProvider, logger)
	c.Notifications = services.NewNotificationsService(t, authProvider, logger)
	c.DataPlus = services.NewDataPlusService(t, authProvider, logger)
	c.HostedPages = services.NewHostedPagesService(t, authProvider, logger)
	c.Constraints = services.NewConstraintsService(t, authProvider, logger)
	c.Application = services.NewApplicationService(t, authProvider, logger)
	c.Webhooks = services.NewWebhooksService(t, authProvider, logger)
	c.Beneficiaries = services.NewBeneficiariesService(t, authProvider, logger)

	return c, nil
}

// Authenticate performs an eager check that credentials are valid by
// making a lightweight call to /institutions.
func (c *Client) Authenticate(ctx context.Context) error {
	_, err := c.auth.GetToken(ctx)
	return err
}

// Logger returns the SDK's logger instance.
func (c *Client) Logger() *zap.Logger { return c.logger }

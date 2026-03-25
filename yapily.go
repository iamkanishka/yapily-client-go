// Package yapily provides a Go SDK for interacting with the Yapily API.
//
// Example:
//
//	client, err := yapily.New(
//	    "your-app-key",
//	    "your-app-secret",
//	    yapily.WithEnvironment(yapily.Sandbox),
//	)
//	if err != nil {
//	    panic(err)
//	}
//
//	accounts, err := client.Accounts.GetAccounts(...)
package yapily

import (
	"github.com/iamkanishka/yapily-client-go/client"
	"github.com/iamkanishka/yapily-client-go/domain"
)

// Re-export core types.
type Client = client.Client
type Config = client.Config
type Option = client.Option
type RateLimitConfig = client.RateLimitConfig

// Re-export environment type.
type Environment = domain.Environment

// Re-export environments.
const (
	Sandbox = domain.Sandbox
	Live    = domain.Production
)

// Re-export constructor and options.
var (
	New = client.New

	WithEnvironment = client.WithEnvironment
	WithBaseURL     = client.WithBaseURL
	WithTimeout     = client.WithTimeout
	WithLogger      = client.WithLogger
	WithRetry       = client.WithRetry
	WithRateLimit   = client.WithRateLimit
)

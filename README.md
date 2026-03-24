# yapily-client-go

[![CI](https://github.com/iamkanishka/yapily-client-go/actions/workflows/ci.yml/badge.svg)](https://github.com/iamkanishka/yapily-client-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/iamkanishka/yapily-client-go.svg)](https://pkg.go.dev/github.com/iamkanishka/yapily-client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/iamkanishka/yapily-client-go)](https://goreportcard.com/report/github.com/iamkanishka/yapily-client-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](./VERSION)

A production-grade Go SDK for the [Yapily Open Banking API](https://docs.yapily.com) (v11.5.x).
Full coverage of every API group: institutions, accounts, transactions, payments, consents,
authorisations, VRP, bulk payments, data enrichment, hosted pages, webhooks, and more.

---

## Features

- **Complete API coverage** — all 17 Yapily API groups implemented
- **HTTP Basic Auth** — correct Yapily auth scheme (`Base64(key:secret)`)
- **Retry + exponential backoff** — configurable attempts, delays
- **Token-bucket rate limiting** — prevent 429s automatically
- **Generic `APIResponse[T]`** — no boilerplate wrappers
- **Full context propagation** — every call accepts `context.Context`
- **PSU headers** — `psu-id`, `psu-corporate-id`, `psu-ip-address` on all auth flows
- **Idempotency keys** — safe payment retries via `Idempotency-Key` header
- **Webhook HMAC-SHA256 verification**
- **Structured logging** — `go.uber.org/zap`
- **Pagination helpers** — `PageIterator[T]` and `ListPaginated` callback API
- **Composable middleware** — logging + retry `http.RoundTripper` wrappers

---

## Installation

```bash
go get github.com/iamkanishka/yapily-client-go@latest
```

Requires **Go 1.21+**.

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/iamkanishka/yapily-client-go/client"
    "github.com/iamkanishka/yapily-client-go/domain"
)

func main() {
    c, err := client.New(
        "your-application-key",
        "your-application-secret",
        client.WithEnvironment(domain.Sandbox),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    institutions, err := c.Institutions.List(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d institutions\n", len(institutions))
}
```

---

## Authentication

Yapily uses **HTTP Basic Authentication**. Your Application Key and Secret are
sent as `Base64(key:secret)` on every request — no OAuth2 token exchange needed.

```go
c, err := client.New(
    os.Getenv("YAPILY_APP_KEY"),
    os.Getenv("YAPILY_APP_SECRET"),
    client.WithEnvironment(domain.Production),
)
```

> Never hard-code credentials. Use environment variables or a secrets manager.

---

## API Coverage

### Institutions
```go
institutions, err := c.Institutions.List(ctx)
inst, err        := c.Institutions.Get(ctx, "monzo")
```

### Users
```go
user, err  := c.Users.Create(ctx, &domain.CreateUserRequest{ApplicationUserID: "user-123"})
users, err := c.Users.List(ctx, "")
user, err  = c.Users.Get(ctx, user.UUID)
err        = c.Users.Update(ctx, user.UUID, &domain.UpdateUserRequest{ApplicationUserID: "new-id"})
err        = c.Users.Delete(ctx, user.UUID)
```

### Authorisations
```go
// Account data access
auth, err := c.Authorisations.CreateAccountAuthorisation(ctx,
    &domain.AccountAuthorisationRequest{
        InstitutionID:    "monzo",
        ApplicationUserID: "user-123",
        FeatureScopeList: []string{"ACCOUNTS", "TRANSACTIONS"},
        Callback:         "https://yourapp.com/callback",
    },
    &domain.PSUHeaders{PSUIPAddress: "1.2.3.4"},
)
// → redirect user to auth.AuthorisationURL

// Payment authorisation
auth, err = c.Authorisations.CreatePaymentAuthorisation(ctx, &domain.PaymentAuthorisationRequest{
    InstitutionID:    "monzo",
    ApplicationUserID: "user-123",
    PaymentRequest:   &domain.PaymentRequest{...},
}, nil)

// Embedded (SCA) flows
auth, err = c.Authorisations.CreateEmbeddedAccountAuthorisation(ctx, req, psu)
auth, err = c.Authorisations.UpdateEmbeddedAccountAuthorisation(ctx, consentToken, req, psu)

// Bulk payment authorisation
auth, err = c.Authorisations.CreateBulkPaymentAuthorisation(ctx, req, psu)

// Pre-authorisations
auth, err = c.Authorisations.CreatePreAuthorisation(ctx, req)
auth, err = c.Authorisations.UpdateAccountPreAuthorisation(ctx, consentToken, req, psu)
```

### Consents
```go
// Exchange callback token
consent, err := c.Consents.ExchangeOneTimeToken(ctx,
    &domain.OneTimeTokenRequest{OneTimeToken: "token-from-callback"})

// Or exchange OAuth2 code
consent, err = c.Consents.ExchangeOAuth2Code(ctx,
    &domain.ExchangeCodeRequest{Code: "auth-code", ApplicationUserID: "user-123"})

// List / Get / Delete / Extend
consents, err := c.Consents.List(ctx, "user-123", "", nil)
consent, err   = c.Consents.Get(ctx, consent.ID)
consent, err   = c.Consents.Extend(ctx, consent.ID, &domain.ExtendConsentRequest{...})
err            = c.Consents.Delete(ctx, consent.ID)
```

### Accounts & Financial Data
```go
consentToken := consent.ID  // or the token from ExchangeOneTimeToken

accounts, err := c.Accounts.List(ctx, consentToken)
account, err  := c.Accounts.Get(ctx, consentToken, accountID)

// Balances, direct debits, scheduled/periodic payments, statements, identity
balances, err := c.FinancialData.GetAccountBalances(ctx, consentToken, accountID)
dds, err      := c.FinancialData.GetAccountDirectDebits(ctx, consentToken, accountID)
sched, err    := c.FinancialData.GetAccountScheduledPayments(ctx, consentToken, accountID)
periodic, err := c.FinancialData.GetAccountPeriodicPayments(ctx, consentToken, accountID)
stmts, err    := c.FinancialData.GetAccountStatements(ctx, consentToken, accountID, nil)
stmt, err     := c.FinancialData.GetAccountStatement(ctx, consentToken, accountID, stmtID)
file, err     := c.FinancialData.GetAccountStatementFile(ctx, consentToken, accountID, stmtID)
identity, err := c.FinancialData.GetIdentity(ctx, consentToken)
rtxns, err    := c.FinancialData.GetRealTimeAccountTransactions(ctx, consentToken, accountID)
```

### Transactions
```go
txns, err := c.Transactions.List(ctx, consentToken, accountID, &domain.TransactionQueryParams{
    PaginationParams: domain.PaginationParams{Limit: 50},
    From:   "2024-01-01",
    Before: "2024-12-31",
})

// Auto-paginate through all pages
err = c.Transactions.ListPaginated(ctx, consentToken, accountID, 50,
    func(page []domain.Transaction) bool {
        // process page; return false to stop
        return true
    },
)
```

### Payments
```go
import "github.com/iamkanishka/yapily-client-go/utils"

payment, err := c.Payments.Create(ctx, consentToken, &domain.PaymentRequest{
    Amount:   100.00,
    Currency: "GBP",
    Recipient: domain.Recipient{
        Name: "Jane Smith",
        AccountIdentifications: []domain.AccountIdentification{
            {Type: "SORT_CODE",      Identification: "200000"},
            {Type: "ACCOUNT_NUMBER", Identification: "55779911"},
        },
    },
    Reference:      "Invoice-001",
    IdempotencyKey: utils.IdempotencyKey(""),
})

payment, err = c.Payments.Get(ctx, consentToken, payment.ID)
```

### Bulk Payments
```go
bulk, err := c.BulkPayments.Create(ctx, consentToken, &domain.BulkPaymentRequest{
    Payments: []domain.PaymentRequest{{...}, {...}},
})
status, err := c.BulkPayments.GetFileStatus(ctx, consentToken, bulk.ID)
```

### Variable Recurring Payments (VRP)
```go
// Create sweeping VRP consent
vrpConsent, err := c.VRP.CreateSweepingAuthorisation(ctx, &domain.VRPAuthorisationRequest{
    InstitutionID:    "monzo",
    ApplicationUserID: "user-123",
    ControlParameters: &domain.VRPControlParameters{
        MaximumIndividualAmount: 500.00,
        Currency:                "GBP",
        PeriodicLimits: []domain.VRPPeriodicLimit{{
            MaximumAmount:   2000.00,
            Currency:        "GBP",
            PeriodType:      "Month",
            PeriodAlignment: "Calendar",
        }},
    },
})

// Check funds, then make a VRP payment
funds, err := c.VRP.ConfirmFunds(ctx, consentToken, vrpConsent.ID,
    &domain.FundsConfirmationRequest{Amount: 100.00, Currency: "GBP"})

vrpPayment, err := c.VRP.CreatePayment(ctx, consentToken, vrpConsent.ID,
    &domain.VRPPaymentRequest{Amount: 100.00, Currency: "GBP", Recipient: domain.Recipient{...}})
```

### Data Plus (Transaction Enrichment)
```go
result, err := c.DataPlus.Enrich(ctx, &domain.EnrichmentRequest{
    Transactions: []domain.RawTransaction{
        {ID: "t1", Amount: 4.50, Currency: "GBP", Description: "Starbucks", Date: "2024-03-01"},
    },
})
// result.Transactions[0].Category == "Food & Drink"

labels, err := c.DataPlus.GetEnrichmentLabels(ctx)
```

### Hosted Pages
```go
// Hosted consent page
hosted, err := c.HostedPages.CreateConsentRequest(ctx, &domain.HostedConsentRequest{
    ApplicationUserID: "user-123",
    FeatureScopeList:  []string{"ACCOUNTS"},
})
// → redirect to hosted.HostedURL

// Hosted payment page
hostedPay, err := c.HostedPages.CreatePaymentRequest(ctx, &domain.HostedPaymentRequest{
    ApplicationUserID: "user-123",
    PaymentRequest:    &domain.PaymentRequest{Amount: 10.00, Currency: "GBP", ...},
})

// Pay By Link
link, err := c.HostedPages.CreatePayByLink(ctx, &domain.PayByLinkRequest{
    Amount: 50.00, Currency: "GBP", Reference: "Order-99",
})
```

### Notifications & Event Subscriptions
```go
sub, err := c.Notifications.CreateEventSubscription(ctx, &domain.CreateEventSubscriptionRequest{
    EventTypeID:     "payment.status.updated",
    NotificationURL: "https://yourapp.com/webhooks",
})
subs, err  := c.Notifications.ListEventSubscriptions(ctx)
err         = c.Notifications.DeleteEventSubscription(ctx, sub.EventTypeID)
```

### Webhooks
```go
cats, err   := c.Webhooks.GetCategories(ctx)
events, err := c.Webhooks.ListEvents(ctx)
event, err  := c.Webhooks.RegisterEvent(ctx, &domain.RegisterWebhookRequest{
    EventTypeID:     "payment.status.updated",
    NotificationURL: "https://yourapp.com/webhooks/payments",
})
secret, err := c.Webhooks.ResetSecret(ctx)
err          = c.Webhooks.DeleteEvent(ctx, event.EventTypeID)

// Verify incoming webhook signatures
import "github.com/iamkanishka/yapily-client-go/middleware"
err = middleware.VerifyWebhookSignature(body, secret.Secret, r.Header.Get("X-Yapily-Signature"))
```

### Constraints
```go
payConstraints, err  := c.Constraints.GetPaymentConstraints(ctx, "monzo", "DOMESTIC_PAYMENT")
dataConstraints, err := c.Constraints.GetDataConstraints(ctx, "monzo")
```

### Application Management
```go
app, err    := c.Application.GetDetails(ctx)
app, err     = c.Application.Update(ctx, &domain.ApplicationRequest{...})
subs, err   := c.Application.ListSubApplications(ctx)
sub, err    := c.Application.CreateSubApplication(ctx, &domain.ApplicationRequest{
    Name: "My Sub-App", MerchantCategoryCode: "5411",
})
vrpCfg, err := c.Application.GetVRPConfiguration(ctx)
```

### Beneficiaries
```go
// Application-level
b, err  := c.Beneficiaries.CreateApplicationBeneficiary(ctx, &domain.CreateBeneficiaryRequest{...})
bs, err := c.Beneficiaries.ListApplicationBeneficiaries(ctx)
b, err   = c.Beneficiaries.GetApplicationBeneficiary(ctx, b.ID)
err      = c.Beneficiaries.DeleteApplicationBeneficiary(ctx, b.ID)

// User-level
ub, err  := c.Beneficiaries.CreateUserBeneficiary(ctx, userUUID, &domain.CreateBeneficiaryRequest{...})
ub, err   = c.Beneficiaries.ApproveBeneficiary(ctx, userUUID, ub.ID)
ub, err   = c.Beneficiaries.PatchUserBeneficiary(ctx, userUUID, ub.ID, &domain.PatchUserBeneficiaryRequest{...})
```

---

## Configuration

```go
c, err := client.New(appKey, appSecret,
    client.WithEnvironment(domain.Sandbox),      // or domain.Production
    client.WithBaseURL("https://custom.api.com"), // optional override
    client.WithTimeout(30 * time.Second),
    client.WithRetry(transporthttp.RetryConfig{
        MaxAttempts: 3,
        BaseDelay:   200 * time.Millisecond,
        MaxDelay:    5 * time.Second,
    }),
    client.WithRateLimit(10, 20), // 10 req/s, burst 20
    client.WithLogger(zapLogger),
)
```

---

## Error Handling

```go
import sdkerrors "github.com/iamkanishka/yapily-client-go/errors"

_, err := c.Accounts.List(ctx, consentToken)
if err != nil {
    switch {
    case sdkerrors.IsNotFound(err):
        // 404
    case sdkerrors.IsUnauthorized(err):
        // 401 — check credentials or consent token
    case sdkerrors.IsRateLimited(err):
        // 429 — already retried; back off further
    default:
        if apiErr, ok := err.(*sdkerrors.APIError); ok {
            fmt.Printf("HTTP %d [%s]: %s (trace: %s)\n",
                apiErr.StatusCode, apiErr.Code, apiErr.Message, apiErr.TraceID)
        }
    }
}
```

---

## Running the Example

```bash
export YAPILY_APP_KEY=your-key
export YAPILY_APP_SECRET=your-secret
export YAPILY_CONSENT_TOKEN=your-consent-token  # optional, for account/payment demo

go run ./examples/basic/
```

---

## Running Tests

```bash
go test ./... -v -race
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Versioning

This SDK follows [Semantic Versioning](https://semver.org/).

```bash
# Tag and publish a new release
git tag v1.0.0
git push origin v1.0.0
# → CI builds cross-platform binaries and creates a GitHub Release automatically
```

---

## Project Structure

```
yapily-client-go/
├── client/                 # Main entry point — client.New() wires all services
├── auth/                   # HTTP Basic Auth provider (+ caching OAuth2 variant)
├── transport/http/         # Retry, rate-limit, logging transport
├── domain/                 # All data models (850+ lines, full Yapily v11.5 schema)
├── errors/                 # Typed errors: APIError, ValidationError, AuthError
├── middleware/             # Logging, retry RoundTripper, webhook HMAC verification
├── utils/                  # PageIterator[T], IdempotencyKey helper
├── services/
│   ├── institutions.go     # GET /institutions
│   ├── accounts.go         # GET /accounts
│   ├── transactions.go     # GET /accounts/{id}/transactions + ListPaginated
│   ├── payments.go         # POST/GET /payments
│   ├── bulk_payments.go    # POST /bulk-payments
│   ├── consents.go         # Full consent lifecycle (7 endpoints)
│   ├── authorisations.go   # All auth flows (12 endpoints)
│   ├── financial_data.go   # Balances, DDs, scheduled, statements, identity, RT txns
│   ├── users.go            # CRUD /users
│   ├── vrp.go              # Variable Recurring Payments
│   ├── notifications.go    # Event subscriptions
│   ├── dataplus.go         # Transaction enrichment
│   ├── hosted_pages.go     # Hosted consent + payment + VRP pages
│   ├── constraints.go      # Payment + data constraints
│   ├── application.go      # Application + sub-application management
│   ├── webhooks.go         # Webhook registration + secret rotation
│   └── beneficiaries.go    # Application + user beneficiaries
├── examples/basic/         # End-to-end runnable example
├── tests/                  # Table-driven tests with httptest servers
└── .github/workflows/      # CI (matrix) + release (cross-platform binaries)
```

---

## License

[MIT](./LICENSE) © yapily-client-go contributors

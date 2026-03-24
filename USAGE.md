# Usage Guide — Yapily Go SDK

Complete reference for every service, method, request type, and response type.
Aligned with [Yapily API v11.5.x](https://docs.yapily.com/api-reference/introduction).
95 methods across 17 services.

---

## Table of Contents

1. [Installation](#installation)
2. [Client Setup](#client-setup)
3. [Authentication](#authentication)
4. [Error Handling](#error-handling)
5. [Institutions](#institutions)
6. [Users](#users)
7. [Authorisations](#authorisations)
8. [Consents](#consents)
9. [Accounts](#accounts)
10. [Financial Data](#financial-data)
11. [Transactions](#transactions)
12. [Payments](#payments)
13. [Bulk Payments](#bulk-payments)
14. [Variable Recurring Payments (VRP)](#variable-recurring-payments-vrp)
15. [Data Plus — Enrichment](#data-plus--enrichment)
16. [Hosted Pages](#hosted-pages)
17. [Notifications](#notifications)
18. [Webhooks](#webhooks)
19. [Constraints](#constraints)
20. [Application Management](#application-management)
21. [Beneficiaries](#beneficiaries)
22. [Pagination](#pagination)
23. [Idempotency](#idempotency)
24. [Webhook Signature Verification](#webhook-signature-verification)
25. [Middleware](#middleware)
26. [Configuration Reference](#configuration-reference)

---

## Installation

```bash
go get github.com/iamkanishka/yapily-client-go@latest
```

Requires Go 1.25+.

---

## Client Setup

```go
import (
    "github.com/iamkanishka/yapily-client-go/client"
    "github.com/iamkanishka/yapily-client-go/domain"
    transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
    "time"
)

c, err := client.New(
    os.Getenv("YAPILY_APP_KEY"),
    os.Getenv("YAPILY_APP_SECRET"),
    client.WithEnvironment(domain.Sandbox),   // domain.Sandbox | domain.Production
    client.WithTimeout(30 * time.Second),
    client.WithRetry(transporthttp.RetryConfig{
        MaxAttempts: 3,
        BaseDelay:   200 * time.Millisecond,
        MaxDelay:    5 * time.Second,
    }),
    client.WithRateLimit(10, 20),   // 10 req/s, burst 20
)
if err != nil {
    log.Fatal(err)
}
```

The `*client.Client` exposes 17 service fields — one per API group:

| Field | Service |
|-------|---------|
| `c.Institutions` | Bank/institution discovery |
| `c.Users` | Application user (PSU) management |
| `c.Authorisations` | All consent initiation flows |
| `c.Consents` | Consent lifecycle (exchange, list, extend, delete) |
| `c.Accounts` | Account listing and retrieval |
| `c.FinancialData` | Balances, DDs, statements, identity, real-time txns |
| `c.Transactions` | Transaction listing + pagination |
| `c.Payments` | Single payment initiation |
| `c.BulkPayments` | Bulk payment initiation |
| `c.VRP` | Variable Recurring Payments |
| `c.DataPlus` | Transaction enrichment / categorisation |
| `c.HostedPages` | Hosted consent, payment, and VRP pages |
| `c.Notifications` | Event subscription management |
| `c.Webhooks` | Webhook registration and secret rotation |
| `c.Constraints` | Payment and data constraint rules |
| `c.Application` | Application and sub-application management |
| `c.Beneficiaries` | Application-level and user-level beneficiaries |

---

## Authentication

Yapily authenticates using **HTTP Basic Auth**: `Base64(ApplicationKey:ApplicationSecret)`.
The SDK applies this automatically on every request.

```go
// Optional eager verification — useful at startup
if err := c.Authenticate(ctx); err != nil {
    log.Fatalf("invalid credentials: %v", err)
}
```

> Never hard-code credentials. Use environment variables or a secrets manager.

---

## Error Handling

```go
import sdkerrors "github.com/iamkanishka/yapily-client-go/errors"

_, err := c.Accounts.List(ctx, consentToken)
if err != nil {
    switch {
    case sdkerrors.IsNotFound(err):
        // 404 — resource does not exist
    case sdkerrors.IsUnauthorized(err):
        // 401 — bad credentials or expired consent
    case sdkerrors.IsRateLimited(err):
        // 429 — already retried; apply additional back-off
    default:
        if e, ok := err.(*sdkerrors.APIError); ok {
            // Full structured error from Yapily
            fmt.Printf("HTTP %d [%s] %s  traceId=%s\n",
                e.StatusCode, e.Code, e.Message, e.TraceID)
        }
        if e, ok := err.(*sdkerrors.ValidationError); ok {
            // Client-side — bad input caught before HTTP call
            fmt.Printf("field '%s': %s\n", e.Field, e.Message)
        }
    }
}
```

**Error types:**

| Type | When |
|------|------|
| `*APIError` | Non-2xx from Yapily. Fields: `StatusCode`, `Code`, `Message`, `TraceID` |
| `*ValidationError` | Missing/invalid input, caught client-side. Fields: `Field`, `Message` |
| `*AuthError` | Credential or token problem. Fields: `Message`, `Cause` |
| `*RetryableError` | Wraps a transient error after all retries exhausted |

**Predicates:** `IsNotFound`, `IsUnauthorized`, `IsRateLimited`, `IsRetryable`

---

## Institutions

```go
// List all institutions configured in your application.
institutions, err := c.Institutions.List(ctx)

// Get one institution by its Yapily ID.
inst, err := c.Institutions.Get(ctx, "monzo")
// inst.ID, inst.Name, inst.FullName, inst.Countries, inst.Features, inst.EnvironmentType
```

---

## Users

```go
// Create an application user (PSU).
user, err := c.Users.Create(ctx, &domain.CreateUserRequest{
    ApplicationUserID: "your-internal-user-id",
})

// List users (pass empty string for all, or an applicationUserId to filter).
users, err := c.Users.List(ctx, "")

// Get a user by Yapily UUID.
user, err = c.Users.Get(ctx, user.UUID)

// Update a user's applicationUserId.
user, err = c.Users.Update(ctx, user.UUID, &domain.UpdateUserRequest{
    ApplicationUserID: "new-internal-id",
})

// Delete a user.
err = c.Users.Delete(ctx, user.UUID)
```

---

## Authorisations

All consent initiation flows. Every method accepts an optional `*domain.PSUHeaders`.

```go
psu := &domain.PSUHeaders{
    PSUID:          "end-user-id",    // forwarded as psu-id header
    PSUCorporateID: "corp-id",        // optional — psu-corporate-id
    PSUIPAddress:   "192.168.1.1",    // strongly recommended — psu-ip-address
}
```

### Account authorisation

```go
auth, err := c.Authorisations.CreateAccountAuthorisation(ctx,
    &domain.AccountAuthorisationRequest{
        InstitutionID:     "monzo",
        ApplicationUserID: "user-123",
        FeatureScopeList:  []string{"ACCOUNTS", "TRANSACTIONS", "IDENTITY"},
        Callback:          "https://yourapp.com/callback",
        OneTimeToken:      true,
    }, psu)
// Redirect user to: auth.AuthorisationURL

// Re-authorise after consent expires.
auth, err = c.Authorisations.ReauthoriseAccountConsent(ctx, consentToken, psu)

// Update an existing pre-authorisation.
auth, err = c.Authorisations.UpdateAccountPreAuthorisation(ctx, consentToken, req, psu)
```

### Payment authorisation

```go
auth, err = c.Authorisations.CreatePaymentAuthorisation(ctx,
    &domain.PaymentAuthorisationRequest{
        InstitutionID:     "monzo",
        ApplicationUserID: "user-123",
        PaymentRequest: &domain.PaymentRequest{
            Amount: 10.00, Currency: "GBP",
            Recipient: domain.Recipient{Name: "Vendor"},
        },
    }, psu)
```

### Bulk payment authorisation

```go
auth, err = c.Authorisations.CreateBulkPaymentAuthorisation(ctx,
    &domain.BulkPaymentAuthorisationRequest{
        InstitutionID:     "monzo",
        ApplicationUserID: "user-123",
        BulkPaymentRequest: &domain.BulkPaymentRequest{
            Payments: []domain.PaymentRequest{
                {Amount: 100, Currency: "GBP", Recipient: domain.Recipient{Name: "Alice"}},
            },
        },
    }, psu)
```

### Pre-authorisation

```go
auth, err = c.Authorisations.CreatePreAuthorisation(ctx,
    &domain.PreAuthorisationRequest{
        InstitutionID: "monzo", ApplicationUserID: "user-123",
    })
```

### Embedded (SCA) flows

```go
// Step 1 — initiate
auth, err = c.Authorisations.CreateEmbeddedAccountAuthorisation(ctx, req, psu)

// Step 2 — submit SCA (OTP / password)
auth, err = c.Authorisations.UpdateEmbeddedAccountAuthorisation(ctx, auth.ID,
    &domain.AccountAuthorisationRequest{
        UserCredentials: &domain.UserCredentials{Password: "123456"},
    }, psu)

// Same pattern:
// CreateEmbeddedPaymentAuthorisation / UpdateEmbeddedPaymentAuthorisation
// CreateEmbeddedBulkPaymentAuthorisation / UpdateEmbeddedBulkPaymentAuthorisation
```

---

## Consents

```go
// After the user is redirected back from the bank:

// Option A — exchange a one-time token from the callback query string.
consent, err := c.Consents.ExchangeOneTimeToken(ctx,
    &domain.OneTimeTokenRequest{OneTimeToken: r.URL.Query().Get("one-time-token")})

// Option B — exchange an OAuth2 authorisation code.
consent, err = c.Consents.ExchangeOAuth2Code(ctx,
    &domain.ExchangeCodeRequest{
        Code:              r.URL.Query().Get("code"),
        ApplicationUserID: "user-123",
        RedirectURL:       "https://yourapp.com/callback",
    })

consentToken := consent.ID  // use for all subsequent data/payment calls

// List consents (filter by user, institution, or paginate).
consents, err := c.Consents.List(ctx, "user-123", "monzo",
    &domain.PaginationParams{Limit: 20})

// Get consent status.
consent, err = c.Consents.Get(ctx, consent.ID)

// Extend consent validity period.
consent, err = c.Consents.Extend(ctx, consent.ID,
    &domain.ExtendConsentRequest{ExpiresAt: &newExpiry})

// Revoke and permanently delete a consent.
err = c.Consents.Delete(ctx, consent.ID)
```

---

## Accounts

```go
// List all accounts accessible under the consent.
accounts, err := c.Accounts.List(ctx, consentToken)
for _, acc := range accounts {
    fmt.Printf("[%s] %s  %.2f %s\n",
        acc.Type, acc.ID, acc.Balance.Amount, acc.Balance.Currency)
}

// Get a single account.
account, err := c.Accounts.Get(ctx, consentToken, accountID)
```

**`domain.Account` key fields:** `ID`, `Type`, `Balance.Amount`, `Balance.Currency`, `InstitutionID`, `AccountIdentifications`

---

## Financial Data

```go
// Detailed balance breakdown.
balances, err := c.FinancialData.GetAccountBalances(ctx, consentToken, accountID)
for _, b := range balances.Balances {
    fmt.Printf("%-20s  %.2f %s\n", b.Type, b.Amount, b.Currency)
}

// Beneficiaries linked to the account.
benes, err := c.FinancialData.GetAccountBeneficiaries(ctx, consentToken, accountID)

// Active direct debit mandates.
dds, err := c.FinancialData.GetAccountDirectDebits(ctx, consentToken, accountID)

// Scheduled (future-dated) payments.
sched, err := c.FinancialData.GetAccountScheduledPayments(ctx, consentToken, accountID)

// Standing orders / recurring payments.
periodic, err := c.FinancialData.GetAccountPeriodicPayments(ctx, consentToken, accountID)

// Statement list.
stmts, err := c.FinancialData.GetAccountStatements(ctx, consentToken, accountID,
    &domain.PaginationParams{Limit: 12})

// Single statement metadata.
stmt, err := c.FinancialData.GetAccountStatement(ctx, consentToken, accountID, statementID)

// Statement file as raw bytes (PDF or CSV depending on institution).
fileBytes, err := c.FinancialData.GetAccountStatementFile(ctx, consentToken, accountID, statementID)

// Identity / KYC data for the account holder.
identity, err := c.FinancialData.GetIdentity(ctx, consentToken)
// identity.FullName, identity.Addresses, identity.Emails, identity.Phones, identity.BirthDate

// Real-time (live) transaction stream.
rtxns, err := c.FinancialData.GetRealTimeAccountTransactions(ctx, consentToken, accountID)
```

---

## Transactions

```go
// Fetch a page with date filters and pagination.
txns, err := c.Transactions.List(ctx, consentToken, accountID,
    &domain.TransactionQueryParams{
        PaginationParams: domain.PaginationParams{Limit: 50, Offset: 0},
        From:   "2024-01-01",  // ISO 8601 date
        Before: "2024-12-31",
    })

// Auto-paginate through ALL transactions with a callback.
err = c.Transactions.ListPaginated(ctx, consentToken, accountID, 50,
    func(page []domain.Transaction) bool {
        for _, tx := range page {
            fmt.Printf("[%s] %+.2f %s  %s\n",
                tx.Date, tx.Amount, tx.Currency, tx.Description)
        }
        return true  // return false to stop after this page
    })
```

**`domain.Transaction` key fields:** `ID`, `Date`, `Amount`, `Currency`, `Description`, `Status`, `Reference`, `Balance`, `Merchant`

---

## Payments

```go
import "github.com/iamkanishka/yapily-client-go/utils"

// Single domestic payment.
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
    Reference:      "Invoice-2024-001",
    IdempotencyKey: utils.IdempotencyKey("order-abc-123"),
    Type:           "DOMESTIC_PAYMENT",  // optional
})

// Scheduled payment.
payment, err = c.Payments.Create(ctx, consentToken, &domain.PaymentRequest{
    Amount: 50.00, Currency: "GBP",
    Recipient:   domain.Recipient{Name: "Landlord"},
    ScheduledAt: &scheduledTime,
    Type:        "DOMESTIC_SCHEDULED_PAYMENT",
})

// Standing order / periodic payment.
payment, err = c.Payments.Create(ctx, consentToken, &domain.PaymentRequest{
    Amount: 200.00, Currency: "GBP",
    Recipient: domain.Recipient{Name: "Savings"},
    Type:      "DOMESTIC_PERIODIC_PAYMENT",
    PeriodicPayment: &domain.PeriodicPaymentRequest{
        Frequency:        "Monthly",
        NumberOfPayments: 12,
        NextPaymentDate:  "2024-02-01",
        FinalPaymentDate: "2025-01-01",
    },
})

// Get payment status.
payment, err = c.Payments.Get(ctx, consentToken, payment.ID)
fmt.Printf("Status: %s  LifecycleID: %s\n", payment.Status, payment.PaymentLifecycleID)
```

---

## Bulk Payments

```go
bulk, err := c.BulkPayments.Create(ctx, consentToken, &domain.BulkPaymentRequest{
    Reference: "Payroll-March-2024",
    Payments: []domain.PaymentRequest{
        {Amount: 3000, Currency: "GBP", Recipient: domain.Recipient{Name: "Alice"}},
        {Amount: 2500, Currency: "GBP", Recipient: domain.Recipient{Name: "Bob"}},
    },
})

// Poll file processing status.
status, err := c.BulkPayments.GetFileStatus(ctx, consentToken, bulk.ID)
```

---

## Variable Recurring Payments (VRP)

```go
// Step 1: Create a sweeping VRP consent (redirect user to auth URL).
vrpConsent, err := c.VRP.CreateSweepingAuthorisation(ctx,
    &domain.VRPAuthorisationRequest{
        InstitutionID:     "monzo",
        ApplicationUserID: "user-123",
        Callback:          "https://yourapp.com/vrp-callback",
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
// Redirect user to: vrpConsent.AuthorisationURL

// Step 2: Retrieve consent details.
vrpConsent, err = c.VRP.GetSweepingConsentDetails(ctx, vrpConsent.ID)

// Step 3 (optional): Check funds.
funds, err := c.VRP.ConfirmFunds(ctx, consentToken, vrpConsent.ID,
    &domain.FundsConfirmationRequest{Amount: 100.00, Currency: "GBP"})
if !funds.FundsAvailable {
    log.Fatal("insufficient funds")
}

// Step 4: Create a VRP payment.
vrpPay, err := c.VRP.CreatePayment(ctx, consentToken, vrpConsent.ID,
    &domain.VRPPaymentRequest{
        Amount:    100.00,
        Currency:  "GBP",
        Reference: "Monthly sweep",
        Recipient: domain.Recipient{
            Name: "Savings Account",
            AccountIdentifications: []domain.AccountIdentification{
                {Type: "SORT_CODE",      Identification: "200000"},
                {Type: "ACCOUNT_NUMBER", Identification: "12345678"},
            },
        },
    })

// Step 5: Get payment details.
details, err := c.VRP.GetPaymentDetails(ctx, vrpConsent.ID, vrpPay.ID)
```

---

## Data Plus — Enrichment

```go
// Submit raw transactions for enrichment (category + merchant detection).
result, err := c.DataPlus.Enrich(ctx, &domain.EnrichmentRequest{
    InstitutionID: "monzo",  // optional
    Transactions: []domain.RawTransaction{
        {ID: "t1", Amount: 4.50,  Currency: "GBP", Description: "Starbucks Canary Wharf", Date: "2024-03-01"},
        {ID: "t2", Amount: 89.00, Currency: "GBP", Description: "Amazon Prime",            Date: "2024-03-02"},
        {ID: "t3", Amount: 12.50, Currency: "GBP", Description: "TfL travel charge",       Date: "2024-03-03"},
    },
})
for _, t := range result.Transactions {
    fmt.Printf("%-40s → %s (%s)\n", t.Description, t.Category, t.Merchant.MerchantName)
}

// Poll async enrichment job.
result, err = c.DataPlus.GetEnrichmentResults(ctx, result.JobID)

// Retrieve all available category labels.
labels, err := c.DataPlus.GetEnrichmentLabels(ctx)
// → ["Shopping", "Food & Drink", "Transport", "Bills & Utilities", …]

// Enrich an account's existing transactions synchronously.
result, err = c.DataPlus.EnrichAccountTransactions(ctx, consentToken, accountID,
    &domain.EnrichmentRequest{Transactions: rawTxns})
```

---

## Hosted Pages

```go
// ── Hosted Consent Page ──────────────────────────────────────────────────────

hosted, err := c.HostedPages.CreateConsentRequest(ctx, &domain.HostedConsentRequest{
    ApplicationUserID:   "user-123",
    FeatureScopeList:    []string{"ACCOUNTS", "TRANSACTIONS"},
    AllowedCountries:    []string{"GB"},
    RedirectURL:         "https://yourapp.com/callback",
})
// Redirect user to: hosted.HostedURL

session, err := c.HostedPages.GetConsentRequest(ctx, hosted.ID)

// ── Hosted Payment Page ──────────────────────────────────────────────────────

hostedPay, err := c.HostedPages.CreatePaymentRequest(ctx, &domain.HostedPaymentRequest{
    ApplicationUserID: "user-123",
    RedirectURL:       "https://yourapp.com/pay-callback",
    PaymentRequest: &domain.PaymentRequest{
        Amount: 25.00, Currency: "GBP",
        Recipient:  domain.Recipient{Name: "Merchant Ltd"},
        Reference:  "Order-9876",
    },
})
// Redirect user to: hostedPay.HostedURL

fetchedPay, err := c.HostedPages.GetPaymentRequest(ctx, hostedPay.ID)

// Check funds (before showing hosted page).
funds, err := c.HostedPages.CheckFundsAvailability(ctx, consentToken,
    &domain.FundsConfirmationRequest{Amount: 25.00, Currency: "GBP"})

// ── Pay By Link ──────────────────────────────────────────────────────────────

link, err := c.HostedPages.CreatePayByLink(ctx, &domain.PayByLinkRequest{
    Amount:      50.00,
    Currency:    "GBP",
    Reference:   "Invoice-123",
    RedirectURL: "https://yourapp.com/thank-you",
    ExpiresAt:   &expiryTime,
})
fmt.Printf("Share: %s\n", link.URL)

// ── Hosted VRP ───────────────────────────────────────────────────────────────

vrp,      err := c.HostedPages.CreateVRPConsent(ctx, &domain.VRPAuthorisationRequest{...})
list,     err  = c.HostedPages.GetVRPConsentRequests(ctx)
vrpReq,   err  = c.HostedPages.GetVRPConsentRequest(ctx, vrp.ID)
vrpPay,   err  = c.HostedPages.CreateVRPPayment(ctx, vrp.ID, &domain.VRPPaymentRequest{...})
fetchPay, err  = c.HostedPages.GetVRPPayment(ctx, vrp.ID, vrpPay.ID)
err             = c.HostedPages.RevokeVRPConsentRequest(ctx, vrp.ID)
```

---

## Notifications

```go
// Create a webhook event subscription.
sub, err := c.Notifications.CreateEventSubscription(ctx,
    &domain.CreateEventSubscriptionRequest{
        EventTypeID:     "payment.status.updated",
        NotificationURL: "https://yourapp.com/webhooks/payments",
    })

subs, err := c.Notifications.ListEventSubscriptions(ctx)
sub,  err  = c.Notifications.GetEventSubscription(ctx, "payment.status.updated")
err        = c.Notifications.DeleteEventSubscription(ctx, "payment.status.updated")
```

---

## Webhooks

```go
// Browse available event types.
categories, err := c.Webhooks.GetCategories(ctx)

// List, register, and delete webhooks.
events, err := c.Webhooks.ListEvents(ctx)
event,  err  = c.Webhooks.RegisterEvent(ctx, &domain.RegisterWebhookRequest{
    EventTypeID:     "payment.status.updated",
    NotificationURL: "https://yourapp.com/webhooks",
})
err = c.Webhooks.DeleteEvent(ctx, "payment.status.updated")

// Rotate the webhook signing secret.
newSecret, err := c.Webhooks.ResetSecret(ctx)
// Store newSecret.Secret securely for signature verification.
```

---

## Constraints

```go
// Payment constraints per institution (limits, currencies, supported types).
pc, err := c.Constraints.GetPaymentConstraints(ctx,
    "monzo",            // institutionID — empty for all
    "DOMESTIC_PAYMENT", // paymentType   — empty for all
)
for _, c := range pc {
    fmt.Printf("%s: max=%.2f currencies=%v\n", c.InstitutionID, *c.MaxAmount, c.Currencies)
}

// Data constraints per institution (max history, available features).
dc, err := c.Constraints.GetDataConstraints(ctx, "monzo")
```

---

## Application Management

```go
// Your application.
app, err := c.Application.GetDetails(ctx)
app, err  = c.Application.Update(ctx, &domain.ApplicationRequest{
    Name: "My Updated App", MerchantCategoryCode: "5411",
})

// Sub-applications (multi-tenant / marketplace).
subs, err := c.Application.ListSubApplications(ctx)
sub,  err  = c.Application.CreateSubApplication(ctx, &domain.ApplicationRequest{
    Name:                 "Client Sub-App",
    MerchantCategoryCode: "5411",
    CallbackURLs:         []string{"https://client.example.com/cb"},
})

// VRP configuration.
cfg, err := c.Application.GetVRPConfiguration(ctx)
cfg, err  = c.Application.CreateVRPConfiguration(ctx, &domain.VRPConfiguration{
    MaxAmount: 1000.00, Currency: "GBP",
    SupportedCurrencies: []string{"GBP", "EUR"},
})
cfg, err = c.Application.UpdateVRPConfiguration(ctx, cfg)
```

---

## Beneficiaries

```go
// ── Application-level (shared across all users) ───────────────────────────────

bene, err := c.Beneficiaries.CreateApplicationBeneficiary(ctx,
    &domain.CreateBeneficiaryRequest{
        Name: "Corporate Supplier",
        AccountIdentifications: []domain.AccountIdentification{
            {Type: "SORT_CODE",      Identification: "200000"},
            {Type: "ACCOUNT_NUMBER", Identification: "12345678"},
        },
    })

all,  err := c.Beneficiaries.ListApplicationBeneficiaries(ctx)
bene, err  = c.Beneficiaries.GetApplicationBeneficiary(ctx, bene.ID)
err        = c.Beneficiaries.DeleteApplicationBeneficiary(ctx, bene.ID)

// ── User-level (per PSU) ─────────────────────────────────────────────────────

ub, err := c.Beneficiaries.CreateUserBeneficiary(ctx, userUUID,
    &domain.CreateBeneficiaryRequest{
        Name: "Jane's Savings",
        AccountIdentifications: []domain.AccountIdentification{
            {Type: "IBAN", Identification: "GB29NWBK60161331926819"},
        },
    })

list, err := c.Beneficiaries.ListUserBeneficiaries(ctx, userUUID)
ub,   err  = c.Beneficiaries.GetUserBeneficiary(ctx, userUUID, ub.ID)
ub,   err  = c.Beneficiaries.PatchUserBeneficiary(ctx, userUUID, ub.ID,
    &domain.PatchUserBeneficiaryRequest{Name: "Jane's Main Savings"})
ub,   err  = c.Beneficiaries.ApproveBeneficiary(ctx, userUUID, ub.ID)
ub,   err  = c.Beneficiaries.RejectBeneficiary(ctx, userUUID, ub.ID)
err        = c.Beneficiaries.DeleteUserBeneficiary(ctx, userUUID, ub.ID)
```

---

## Pagination

### Generic `PageIterator[T]`

```go
import "github.com/iamkanishka/yapily-client-go/utils"

iter := utils.NewPageIterator[domain.Transaction](50,
    func(p domain.PaginationParams) ([]domain.Transaction, error) {
        return c.Transactions.List(ctx, consentToken, accountID,
            &domain.TransactionQueryParams{PaginationParams: p})
    })

// Collect all pages into one slice.
all, err := iter.Collect()

// Or iterate page by page.
for {
    page, ok := iter.Next()
    if !ok { break }
    fmt.Printf("Page of %d transactions\n", len(page))
}
```

### Callback pagination (transactions)

```go
err = c.Transactions.ListPaginated(ctx, consentToken, accountID, 50,
    func(page []domain.Transaction) bool {
        // process page
        return true  // false = stop after this page
    })
```

---

## Idempotency

Always use idempotency keys for payment creation to make retries safe.

```go
import "github.com/iamkanishka/yapily-client-go/utils"

// Use your own stable key (recommended — tied to your internal order ID).
key := utils.IdempotencyKey("order-abc-123")

// Or auto-generate (timestamp-based, unique within a process).
key = utils.IdempotencyKey("")  // → "idem-1711024800000000000"
```

If you retry a payment with the same idempotency key, Yapily returns the original
response rather than processing a duplicate payment.

---

## Webhook Signature Verification

```go
import "github.com/iamkanishka/yapily-client-go/middleware"

http.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    sig := r.Header.Get("X-Yapily-Signature")  // format: "sha256=<hex>"

    if err := middleware.VerifyWebhookSignature(body, webhookSecret, sig); err != nil {
        // middleware.ErrInvalidWebhookSignature
        http.Error(w, "invalid signature", http.StatusUnauthorized)
        return
    }

    var event map[string]interface{}
    json.Unmarshal(body, &event)
    // Process the verified event...
    w.WriteHeader(http.StatusOK)
})
```

Rotate your secret periodically with `c.Webhooks.ResetSecret(ctx)`.

---

## Middleware

```go
import (
    "github.com/iamkanishka/yapily-client-go/middleware"
    "go.uber.org/zap"
)

logger, _ := zap.NewProduction()

// Compose: logging (outer) → retry (inner) → net/http default
transport := middleware.NewLoggingRoundTripper(
    middleware.NewRetryRoundTripper(
        http.DefaultTransport,
        3,                      // maxAttempts
        200*time.Millisecond,   // baseDelay
        5*time.Second,          // maxDelay
        logger,
    ),
    logger,
)

// Pass to the SDK client's transport option.
c, _ := client.New("key", "secret",
    client.WithBaseURL("https://api.yapily.com"),
    // Inject via transport option if needed for custom http.Client use cases.
)
```

---

## Configuration Reference

### `client.New` options

| Option | Default | Description |
|--------|---------|-------------|
| `WithEnvironment(env)` | `domain.Sandbox` | Switch between `Sandbox` and `Production` |
| `WithBaseURL(url)` | `https://api.yapily.com` | Override the API base URL |
| `WithTimeout(d)` | `30s` | Per-request HTTP timeout |
| `WithRetry(cfg)` | 3 attempts, 200ms–5s backoff | Retry policy for transient failures |
| `WithRateLimit(rps, burst)` | 10 rps / 20 burst | Token-bucket rate limiter |
| `WithLogger(l)` | zap production logger | Custom `*zap.Logger` |

### `RetryConfig` fields

| Field | Default | Description |
|-------|---------|-------------|
| `MaxAttempts` | `3` | Total attempts (first + retries) |
| `BaseDelay` | `200ms` | Initial delay before first retry |
| `MaxDelay` | `5s` | Cap on exponential backoff delay |

Retries trigger on HTTP `429`, `500`, `503`, `504` only. `4xx` errors (except 429) are not retried.

### `PSUHeaders` fields

| Field | HTTP Header | When to use |
|-------|-------------|-------------|
| `PSUID` | `psu-id` | End-user identifier at the bank |
| `PSUCorporateID` | `psu-corporate-id` | Business/corporate account flows |
| `PSUIPAddress` | `psu-ip-address` | Required by many institutions |

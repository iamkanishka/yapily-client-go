# Changelog

All notable changes to the OpenBanking Go SDK are documented here.

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
Format based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

---

## [Unreleased]

> Changes staged for the next release.

---

## [1.0.2] — 2026-03-23

### Security

- **GO-2026-4601** — `net/url` incorrect parsing of IPv6 host literals.
  Fixed by updating `go` directive in `go.mod` from `1.21` to `1.25.8`,
  where this vulnerability is patched. Affected call sites:
  `transport/http/client.go` `buildRequest` and `Request`.
- **GO-2026-4337** — `crypto/tls` unexpected session resumption.
  Fixed by the same `go 1.25.8` update. Affected call sites:
  `middleware/retry.go` `RoundTrip`, `examples/basic/main.go`.

### Fixed (golangci-lint)

- **`transport/http/client.go`** — Renamed exported interface `HTTPClient` →
  `Client` to avoid the stutter `http.HTTPClient` (`revive` rule).
  All call sites updated (`client.go`, `middleware/retry.go`,
  `middleware/logging.go`).
- **`transport/http/client.go`** — `zap.NewProduction()` error return now
  handled; falls back to `zap.NewNop()` on failure (`errcheck`).
- **`tests/accounts_test.go`** — `zap.NewDevelopment()` replaced with
  `zap.NewNop()` (no error to discard); `json.Encoder.Encode` error now
  handled (`errcheck`).
- **`tests/auth_test.go`** — `provider.GetToken` return values no longer
  discarded with `_` (`errcheck`).
- **`tests/authorisations_test.go`** — `svc.CreateAccountAuthorisation`
  return value now checked (`errcheck`).
- **`tests/payments_test.go`** — `json.Decoder.Decode` and `svc.Create`
  errors now handled (`errcheck`).
- **`tests/transport_test.go`** — `json.Encoder.Encode`, `json.Decoder.Decode`,
  and `transport.Request` errors now handled (`errcheck`).
- **`services/consents.go`** — Removed dead function `validateConsentRequest`
  (`unused`).
- **`errors/errors.go`** — Trailing whitespace removed to satisfy `gofmt`.
- **`services/*.go`** — All exported doc comments now end with a period
  (`godot`), ~95 comments fixed.
- **`.golangci.yml`** — Replaced deprecated `exportloopref` linter (removed
  since Go 1.22) with `copyloopvar`; updated `go` version to `1.25`.

### Changed

- **Module path** renamed from `github.com/yourusername/openbanking-go-sdk`
  to `github.com/iamkanishka/yapily` across all 42 Go files and `go.mod`.

---

## [1.0.1] — 2026-03-21

### Fixed (discovered via live Yapily OpenAPI spec v12.0.0)

- **`consents.go`** `ExchangeOAuth2Code` — endpoint corrected from
  `POST /consent-auth-codes` → `POST /consent-auth-code` (singular).
  The old path returns 404 on real Yapily.
- **`consents.go`** `ExchangeOneTimeToken` — endpoint corrected from
  `POST /consent-one-time-tokens` → `POST /consent-one-time-token` (singular).
- **`payments.go`** `Get` — endpoint corrected from
  `GET /payments/{id}/details` → `GET /payments/{id}`.
- **`bulk_payments.go`** `GetStatus` — endpoint corrected from
  `GET /bulk-payments/{id}/file` → `GET /bulk-payments/{id}`.
  Method renamed `GetFileStatus` → `GetStatus` to match.
- **`notifications.go`** — all endpoints corrected from
  `/notifications/event-subscriptions` prefix → `/event-subscriptions`
  (the real Yapily path has no `/notifications/` prefix).
- **`application.go`** sub-application endpoints corrected from
  `GET/POST /application/sub-applications` → `GET/POST /applications`.
- **`consents.go`** removed non-existent `POST /consents` endpoint (no such
  path exists in the Yapily spec).
- **`consents.go`** `Delete` now accepts `forceDelete bool` parameter,
  forwarded as `?forceDelete=true` query param per spec.
- **`consents.go`** `List` signature updated to `ConsentListParams` struct
  supporting all Yapily filter parameters: `filter[applicationUserId]`,
  `filter[userUuid]`, `filter[institution]`, `filter[status]`, `from`,
  `before`, `limit`, `offset`.
- **`consents.go`** `ExchangeOAuth2Code` and `ExchangeOneTimeToken` now
  decode response as `*domain.Consent` directly (not wrapped in
  `APIResponse`) — matching actual Yapily response format.
- **`bulk_payments.go`** response type corrected to `BulkPaymentStatus`
  containing `StatusDetails.Status` — matching real API response shape.

### Added

- **`authorisations.go`** `CreatePaymentPreAuthorisation` — `POST /payment-pre-auth-requests`
- **`authorisations.go`** `UpdatePaymentPreAuthorisation` — `PUT /payment-pre-auth-requests`
- **`dataplus.go`** `TransactionsAndEnrichment` — `POST /enrich-requests`
  (new Yapily v12 endpoint)
- **`dataplus.go`** `GetEnrichmentRequestResults` — `GET /enrich-requests/{jobId}`
- **`domain`** `ConsentListParams` — typed filter struct for `Consents.List`
- **`domain`** `BulkPaymentStatus` and `BulkPaymentStatusDetail` — types
  matching real bulk payment status response

---

## [1.0.0] — 2026-03-21

Initial stable release. Full coverage of the Yapily API v11.5.x.

### Added

#### Core Infrastructure
- `client.New()` — main SDK entry point with functional options pattern
- HTTP Basic Authentication (`Base64(ApplicationKey:ApplicationSecret)`) matching Yapily's actual auth scheme
- Caching OAuth2 provider (`NewCachingOAuth2Provider`) for Bearer token flows
- Exponential backoff retry — configurable `MaxAttempts`, `BaseDelay`, `MaxDelay`; retries on `429`, `500`, `503`, `504`
- Token-bucket rate limiter (`golang.org/x/time/rate`) — default 10 req/s, burst 20
- Context propagation on every HTTP call
- Structured logging via `go.uber.org/zap`
- Generic `APIResponse[T]` wrapper — eliminates per-endpoint response structs
- Full `go.mod` module at `github.com/yourusername/openbanking-go-sdk`

#### Error Types (`errors` package)
- `APIError` — non-2xx response; fields: `StatusCode`, `Code`, `Message`, `TraceID`
- `ValidationError` — client-side input validation; fields: `Field`, `Message`
- `AuthError` — credential/token failure; fields: `Message`, `Cause`
- `RetryableError` — wraps a transient failure after all retries exhausted
- Predicates: `IsNotFound`, `IsUnauthorized`, `IsRateLimited`, `IsRetryable`

#### Domain Models (`domain` package)
- 50+ types covering the full Yapily v11.5 schema
- `Institution`, `Country`, `Media`, `Monitoring`
- `Account`, `AccountBalance`, `BalanceDetail`, `CreditLine`, `AccountIdentification`
- `Transaction`, `IsoBankTxCode`, `ProprietaryBankTxCode`, `Merchant`, `TransactionEnrichment`
- `PaymentRequest`, `Payment`, `PeriodicPaymentRequest`, `InternationalPaymentRequest`, `ExchangeRateInfo`
- `BulkPaymentRequest`, `BulkPayment`
- `Recipient`, `Address`, `ChargeInfo`, `PaymentStatusDetail`
- `ConsentRequest`, `Consent`, `ExchangeCodeRequest`, `OneTimeTokenRequest`, `ExtendConsentRequest`
- `AccountAuthorisationRequest`, `PaymentAuthorisationRequest`, `BulkPaymentAuthorisationRequest`
- `Authorisation`, `PreAuthorisationRequest`, `UserCredentials`, `ForwardParameter`
- `PSUHeaders` — `psu-id`, `psu-corporate-id`, `psu-ip-address` header forwarding
- `User`, `CreateUserRequest`, `UpdateUserRequest`
- `Beneficiary`, `DirectDebit`, `ScheduledPayment`, `PeriodicPaymentResponse`, `Statement`, `Identity`
- `VRPAuthorisationRequest`, `VRPControlParameters`, `VRPPeriodicLimit`, `VRPConsent`
- `VRPPaymentRequest`, `VRPPayment`, `FundsConfirmationRequest`, `FundsConfirmationResponse`
- `EventSubscription`, `CreateEventSubscriptionRequest`
- `EnrichmentRequest`, `RawTransaction`, `EnrichedTransaction`, `EnrichmentResult`
- `HostedConsentRequest`, `HostedConsentResponse`, `HostedPaymentRequest`, `HostedPaymentResponse`
- `PayByLinkRequest`, `PayByLinkResponse`
- `PaymentConstraint`, `DataConstraint`
- `ApplicationRequest`, `Application`, `VRPConfiguration`
- `WebhookCategory`, `WebhookEvent`, `RegisterWebhookRequest`, `ResetWebhookSecretResponse`
- `ApplicationBeneficiary`, `UserBeneficiary`, `CreateBeneficiaryRequest`, `PatchUserBeneficiaryRequest`
- `PaginationParams`, `TransactionQueryParams`

#### Services (95 methods across 17 services)

**`InstitutionsService`** — 2 methods
- `List(ctx)` → `GET /institutions`
- `Get(ctx, institutionID)` → `GET /institutions/{id}`

**`UsersService`** — 5 methods
- `List(ctx, applicationUserID)` → `GET /users`
- `Create(ctx, req)` → `POST /users`
- `Get(ctx, userUUID)` → `GET /users/{uuid}`
- `Update(ctx, userUUID, req)` → `PATCH /users/{uuid}`
- `Delete(ctx, userUUID)` → `DELETE /users/{uuid}`

**`AuthorisationsService`** — 12 methods
- `CreateAccountAuthorisation(ctx, req, psu)` → `POST /account-auth-requests`
- `ReauthoriseAccountConsent(ctx, consentToken, psu)` → `PATCH /account-auth-requests`
- `UpdateAccountPreAuthorisation(ctx, consentToken, req, psu)` → `PUT /account-auth-requests/{token}`
- `CreatePaymentAuthorisation(ctx, req, psu)` → `POST /payment-auth-requests`
- `CreateBulkPaymentAuthorisation(ctx, req, psu)` → `POST /bulk-payment-auth-requests`
- `CreatePreAuthorisation(ctx, req)` → `POST /pre-auth-requests`
- `CreateEmbeddedAccountAuthorisation(ctx, req, psu)` → `POST /embedded-account-auth-requests`
- `UpdateEmbeddedAccountAuthorisation(ctx, token, req, psu)` → `PUT /embedded-account-auth-requests/{token}`
- `CreateEmbeddedPaymentAuthorisation(ctx, req, psu)` → `POST /embedded-payment-auth-requests`
- `UpdateEmbeddedPaymentAuthorisation(ctx, token, req, psu)` → `PUT /embedded-payment-auth-requests/{token}`
- `CreateEmbeddedBulkPaymentAuthorisation(ctx, req, psu)` → `POST /embedded-bulk-payment-auth-requests`
- `UpdateEmbeddedBulkPaymentAuthorisation(ctx, token, req, psu)` → `PUT /embedded-bulk-payment-auth-requests/{token}`

**`ConsentsService`** — 7 methods
- `ExchangeOneTimeToken(ctx, req)` → `POST /consent-one-time-tokens`
- `ExchangeOAuth2Code(ctx, req)` → `POST /consent-auth-codes`
- `List(ctx, userID, institutionID, params)` → `GET /consents`
- `Get(ctx, consentID)` → `GET /consents/{id}`
- `Extend(ctx, consentID, req)` → `POST /consents/{id}/extend`
- `Delete(ctx, consentID)` → `DELETE /consents/{id}`
- `Create(ctx, req)` → `POST /consents` *(legacy helper)*

**`AccountsService`** — 2 methods
- `List(ctx, consentToken)` → `GET /accounts`
- `Get(ctx, consentToken, accountID)` → `GET /accounts/{id}`

**`FinancialDataService`** — 10 methods
- `GetAccountBalances` → `GET /accounts/{id}/balances`
- `GetAccountBeneficiaries` → `GET /accounts/{id}/beneficiaries`
- `GetAccountDirectDebits` → `GET /accounts/{id}/direct-debits`
- `GetAccountScheduledPayments` → `GET /accounts/{id}/scheduled-payments`
- `GetAccountPeriodicPayments` → `GET /accounts/{id}/periodic-payments`
- `GetAccountStatements` → `GET /accounts/{id}/statements`
- `GetAccountStatement` → `GET /accounts/{id}/statements/{stmtId}`
- `GetAccountStatementFile` → `GET /accounts/{id}/statements/{stmtId}/file`
- `GetIdentity` → `GET /identity`
- `GetRealTimeAccountTransactions` → `GET /accounts/{id}/transactions/real-time`

**`TransactionsService`** — 2 methods
- `List(ctx, consentToken, accountID, params)` → `GET /accounts/{id}/transactions`
- `ListPaginated(ctx, consentToken, accountID, pageSize, fn)` — auto-paginates with callback

**`PaymentsService`** — 2 methods
- `Create(ctx, consentToken, req)` → `POST /payments` (includes `Idempotency-Key` header)
- `Get(ctx, consentToken, paymentID)` → `GET /payments/{id}/details`

**`BulkPaymentsService`** — 2 methods
- `Create(ctx, consentToken, req)` → `POST /bulk-payments`
- `GetFileStatus(ctx, consentToken, paymentID)` → `GET /bulk-payments/{id}/file`

**`VRPService`** — 5 methods
- `CreateSweepingAuthorisation(ctx, req)` → `POST /vrp-consents`
- `GetSweepingConsentDetails(ctx, consentID)` → `GET /vrp-consents/{id}`
- `CreatePayment(ctx, consentToken, consentID, req)` → `POST /vrp-consents/{id}/payments`
- `GetPaymentDetails(ctx, consentID, paymentID)` → `GET /vrp-consents/{id}/payments/{payId}`
- `ConfirmFunds(ctx, consentToken, consentID, req)` → `POST /vrp-consents/{id}/funds-confirmation`

**`DataPlusService`** — 4 methods
- `Enrich(ctx, req)` → `POST /enrich`
- `GetEnrichmentResults(ctx, jobID)` → `GET /enrich/{jobId}`
- `GetEnrichmentLabels(ctx)` → `GET /enrich/labels`
- `EnrichAccountTransactions(ctx, consentToken, accountID, req)` → `POST /accounts/{id}/transactions/categorisation`

**`HostedPagesService`** — 11 methods
- `CreateConsentRequest` → `POST /hosted/consent/account-auth-requests`
- `GetConsentRequest` → `GET /hosted/consent/account-auth-requests/{id}`
- `CreatePaymentRequest` → `POST /hosted/payment/payment-auth-requests`
- `GetPaymentRequest` → `GET /hosted/payment/payment-auth-requests/{id}`
- `CreatePayByLink` → `POST /hosted/payment/pay-by-link`
- `CheckFundsAvailability` → `POST /hosted/payment/funds-confirmation`
- `GetVRPConsentRequests` → `GET /hosted/vrp/vrp-consent-auth-requests`
- `CreateVRPConsent` → `POST /hosted/vrp/vrp-consent-auth-requests`
- `GetVRPConsentRequest` → `GET /hosted/vrp/vrp-consent-auth-requests/{id}`
- `RevokeVRPConsentRequest` → `POST /hosted/vrp/vrp-consent-auth-requests/{id}/revoke`
- `CreateVRPPayment` → `POST /hosted/vrp/vrp-consent-auth-requests/{id}/payments`
- `GetVRPPayment` → `GET /hosted/vrp/vrp-consent-auth-requests/{id}/payments/{payId}`

**`NotificationsService`** — 4 methods
- `ListEventSubscriptions` → `GET /notifications/event-subscriptions`
- `CreateEventSubscription` → `POST /notifications/event-subscriptions`
- `GetEventSubscription` → `GET /notifications/event-subscriptions/{id}`
- `DeleteEventSubscription` → `DELETE /notifications/event-subscriptions/{id}`

**`WebhooksService`** — 5 methods
- `GetCategories` → `GET /webhooks/categories`
- `ListEvents` → `GET /webhooks`
- `RegisterEvent` → `POST /webhooks`
- `DeleteEvent` → `DELETE /webhooks/{id}`
- `ResetSecret` → `POST /webhooks/secret/reset`

**`ConstraintsService`** — 2 methods
- `GetPaymentConstraints(ctx, institutionID, paymentType)` → `GET /constraints/payment`
- `GetDataConstraints(ctx, institutionID)` → `GET /constraints/data`

**`ApplicationService`** — 8 methods
- `GetDetails` → `GET /application`
- `Update` → `PUT /application`
- `Delete` → `DELETE /application`
- `ListSubApplications` → `GET /application/sub-applications`
- `CreateSubApplication` → `POST /application/sub-applications`
- `GetVRPConfiguration` → `GET /application/vrp-configuration`
- `CreateVRPConfiguration` → `POST /application/vrp-configuration`
- `UpdateVRPConfiguration` → `PUT /application/vrp-configuration`

**`BeneficiariesService`** — 11 methods
- `ListApplicationBeneficiaries` → `GET /application/beneficiaries`
- `CreateApplicationBeneficiary` → `POST /application/beneficiaries`
- `GetApplicationBeneficiary` → `GET /application/beneficiaries/{id}`
- `DeleteApplicationBeneficiary` → `DELETE /application/beneficiaries/{id}`
- `ListUserBeneficiaries` → `GET /users/{uuid}/beneficiaries`
- `CreateUserBeneficiary` → `POST /users/{uuid}/beneficiaries`
- `GetUserBeneficiary` → `GET /users/{uuid}/beneficiaries/{id}`
- `DeleteUserBeneficiary` → `DELETE /users/{uuid}/beneficiaries/{id}`
- `PatchUserBeneficiary` → `PATCH /users/{uuid}/beneficiaries/{id}`
- `ApproveBeneficiary` → `POST /users/{uuid}/beneficiaries/{id}/approve`
- `RejectBeneficiary` → `POST /users/{uuid}/beneficiaries/{id}/reject`

#### Utilities (`utils` package)
- `PageIterator[T]` — generic offset-based paginator with `Next()` and `Collect()`
- `IdempotencyKey(input string)` — returns `input` if non-empty, else generates `idem-<nanoseconds>`

#### Middleware (`middleware` package)
- `NewLoggingRoundTripper` — logs method, URL, status code, and elapsed time
- `NewRetryRoundTripper` — retries `429`/`5xx` with exponential backoff; safe body replay
- `VerifyWebhookSignature(payload, secret, header)` — HMAC-SHA256 `sha256=<hex>` verification
- `WebhookHandler(secret, inner)` — `http.Handler` wrapper that validates signatures

#### Testing
- 19 test files, 113 test scenarios
- `httptest`-server-based integration tests for all 17 services
- Concurrent auth provider safety test (double-checked locking)
- Pagination: page-by-page, early-stop, empty-result, collect-all
- Transport: retry-on-500, no-retry-on-400, context cancellation, query params
- Table-driven validation tests for every service method

#### CI/CD
- `.github/workflows/ci.yml` — Go matrix (1.21, 1.22), `go vet`, `golangci-lint`, race detector, Codecov
- `.github/workflows/release.yml` — triggered on `v*` tags; builds Linux/macOS/Windows binaries, creates GitHub Release with checksums
- `.golangci.yml` — enabled linters: `errcheck`, `staticcheck`, `bodyclose`, `noctx`, `gocritic`, `revive`, `goimports`, `godot`

---

## Versioning Policy

This SDK follows **Semantic Versioning**:

| Version bump | Reason |
|---|---|
| **Patch** `x.y.Z` | Bug fixes, documentation corrections, test additions |
| **Minor** `x.Y.0` | New endpoints, new options, backwards-compatible additions |
| **Major** `X.0.0` | Breaking API changes (method signature, removed field, renamed type) |

### Releasing a new version

```bash
# Update VERSION file
echo "1.1.0" > VERSION

# Update go.mod if needed
# Commit
git add -A && git commit -m "chore: release v1.1.0"

# Tag
git tag v1.1.0
git push origin main --tags
# → CI builds cross-platform binaries and publishes a GitHub Release automatically
```

### Updating in a project

```bash
# Pin to a specific version
go get github.com/yourusername/openbanking-go-sdk@v1.1.0

# Upgrade to latest
go get github.com/yourusername/openbanking-go-sdk@latest

# View available versions
go list -m -versions github.com/yourusername/openbanking-go-sdk
```

---

## Migration Guide

### From pre-release / prototype to v1.0.0

| Before | After | Notes |
|--------|-------|-------|
| `client.New(clientID, clientSecret, ...)` | `client.New(appKey, appSecret, ...)` | Param names reflect Yapily terminology |
| `auth.NewOAuth2Provider(cfg)` | `auth.NewBasicAuthProvider(cfg)` | Yapily uses Basic auth, not OAuth2. `NewOAuth2Provider` still exists as an alias. |
| `cfg.ClientID` / `cfg.ClientSecret` | `cfg.ApplicationKey` / `cfg.ApplicationSecret` | Matches Yapily Console field names |
| `c.Consents.Create(ctx, req)` | `c.Authorisations.CreateAccountAuthorisation(ctx, req, psu)` | Prefer the Authorisations service for new integrations |
| `payment.Get(ctx, token, id)` | `payment.Get(ctx, token, id)` — path now `/payments/{id}/details` | Aligned with Yapily v11 spec |

---

[Unreleased]: https://github.com/iamkanishka/yapily/compare/v1.0.2...HEAD
[1.0.2]: https://github.com/iamkanishka/yapily/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/iamkanishka/yapily/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/iamkanishka/yapily/releases/tag/v1.0.0

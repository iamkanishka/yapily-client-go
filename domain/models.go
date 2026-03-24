// Package domain contains all core data models for the Yapily SDK,
// aligned with the Yapily API v11.5.x.
package domain

import "time"

// ─── Environment ─────────────────────────────────────────────────────────────

// Environment represents the API environment.
type Environment string

const (
	Sandbox    Environment = "sandbox"
	Production Environment = "production"
)

// ─── Generic wrapper ──────────────────────────────────────────────────────────

// APIResponse is a generic wrapper for all paginated API responses.
type APIResponse[T any] struct {
	Meta  *ResponseMeta `json:"meta,omitempty"`
	Data  T             `json:"data"`
	Links *Links        `json:"links,omitempty"`
}

// ResponseMeta is the meta block returned on every Yapily response.
type ResponseMeta struct {
	TracingID  string      `json:"tracingId,omitempty"`
	Count      int         `json:"count,omitempty"`
	TotalCount int         `json:"totalCount,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination holds cursor-based pagination info.
type Pagination struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

// Links holds hypermedia links.
type Links struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

// ─── Auth token ───────────────────────────────────────────────────────────────

// Token represents an OAuth2 access token.
type Token struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	ExpiresAt   time.Time `json:"-"`
}

// IsExpired returns true if the token has expired (with 30s buffer).
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt.Add(-30 * time.Second))
}

// ─── Institutions ─────────────────────────────────────────────────────────────

// Institution represents a bank or financial institution.
type Institution struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	FullName          string      `json:"fullName"`
	Countries         []Country   `json:"countries"`
	Features          []string    `json:"features"`
	Media             []Media     `json:"media"`
	CredentialsType   string      `json:"credentialsType"`
	EnvironmentType   string      `json:"environmentType"`
	AuthorisationFlow string      `json:"authorisationFlow,omitempty"`
	Monitoring        *Monitoring `json:"monitoring,omitempty"`
}

// Country represents a country where an institution operates.
type Country struct {
	CountryCode string `json:"countryCode"`
	DisplayName string `json:"displayName"`
}

// Media holds logo/icon URLs for an institution.
type Media struct {
	Source string `json:"source"`
	Type   string `json:"type"`
}

// Monitoring holds institution uptime/monitoring data.
type Monitoring struct {
	Uptime   float64 `json:"uptime,omitempty"`
	Response float64 `json:"response,omitempty"`
}

// ─── Accounts ─────────────────────────────────────────────────────────────────

// Balance holds monetary balance information.
type Balance struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// Account represents a bank account.
type Account struct {
	ID                             string                          `json:"id"`
	Type                           string                          `json:"type"`
	Balance                        Balance                         `json:"balance"`
	InstitutionID                  string                          `json:"institutionId"`
	Currency                       string                          `json:"currency,omitempty"`
	UsageType                      string                          `json:"usageType,omitempty"`
	AccountNames                   []AccountName                   `json:"accountNames,omitempty"`
	AccountIdentifications         []AccountIdentification         `json:"accountIdentifications,omitempty"`
	ConsolidatedAccountInformation *ConsolidatedAccountInformation `json:"consolidatedAccountInformation,omitempty"`
}

// AccountName holds display name info for an account.
type AccountName struct {
	Name string `json:"name"`
}

// AccountIdentification holds account number identifiers.
type AccountIdentification struct {
	Type           string `json:"type"`
	Identification string `json:"identification"`
}

// ConsolidatedAccountInformation holds aggregated account details.
type ConsolidatedAccountInformation struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// AccountBalance holds detailed balance breakdown.
type AccountBalance struct {
	AccountID string          `json:"accountId"`
	DateTime  time.Time       `json:"dateTime"`
	Balances  []BalanceDetail `json:"balances"`
}

// BalanceDetail is a single balance entry.
type BalanceDetail struct {
	Type        string       `json:"type"`
	DateTime    string       `json:"dateTime,omitempty"`
	Amount      float64      `json:"amount"`
	Currency    string       `json:"currency"`
	CreditLines []CreditLine `json:"creditLines,omitempty"`
}

// CreditLine holds credit line detail on a balance.
type CreditLine struct {
	Agreed   bool    `json:"agreed"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// ─── Transactions ─────────────────────────────────────────────────────────────

// Transaction represents a bank transaction.
type Transaction struct {
	ID                             string                 `json:"id"`
	Date                           string                 `json:"date"`
	BookingDateTime                *time.Time             `json:"bookingDateTime,omitempty"`
	ValueDateTime                  *time.Time             `json:"valueDateTime,omitempty"`
	Amount                         float64                `json:"amount"`
	Currency                       string                 `json:"currency"`
	Description                    string                 `json:"description"`
	TransactionType                string                 `json:"transactionType,omitempty"`
	Status                         string                 `json:"status,omitempty"`
	Reference                      string                 `json:"reference,omitempty"`
	Balance                        *Balance               `json:"balance,omitempty"`
	IsoBankTransactionCode         *IsoBankTxCode         `json:"isoBankTransactionCode,omitempty"`
	ProprietaryBankTransactionCode *ProprietaryBankTxCode `json:"proprietaryBankTransactionCode,omitempty"`
	Merchant                       *Merchant              `json:"merchant,omitempty"`
	Enrichment                     *TransactionEnrichment `json:"enrichment,omitempty"`
}

// IsoBankTxCode holds ISO 20022 transaction code.
type IsoBankTxCode struct {
	DomainCode *DomainCode `json:"domainCode,omitempty"`
}

// DomainCode holds domain/family/sub-family transaction classification.
type DomainCode struct {
	Code   string      `json:"code"`
	Family *FamilyCode `json:"family,omitempty"`
}

// FamilyCode holds family and sub-family codes.
type FamilyCode struct {
	Code      string         `json:"code"`
	SubFamily *SubFamilyCode `json:"subFamily,omitempty"`
}

// SubFamilyCode holds sub-family code.
type SubFamilyCode struct {
	Code string `json:"code"`
}

// ProprietaryBankTxCode holds bank-specific transaction codes.
type ProprietaryBankTxCode struct {
	Code   string `json:"code"`
	Issuer string `json:"issuer,omitempty"`
}

// Merchant holds merchant information on a transaction.
type Merchant struct {
	MerchantName         string `json:"merchantName,omitempty"`
	MerchantCategoryCode string `json:"merchantCategoryCode,omitempty"`
}

// TransactionEnrichment holds Data Plus enrichment on a transaction.
type TransactionEnrichment struct {
	TransactionID string `json:"transactionId,omitempty"`
	Category      string `json:"category,omitempty"`
	MerchantName  string `json:"merchantName,omitempty"`
}

// ─── Payments ─────────────────────────────────────────────────────────────────

// Recipient represents the recipient of a payment.
type Recipient struct {
	Name                   string                  `json:"name"`
	AccountNumber          string                  `json:"accountNumber,omitempty"`
	SortCode               string                  `json:"sortCode,omitempty"`
	AccountIdentifications []AccountIdentification `json:"accountIdentifications,omitempty"`
	Address                *Address                `json:"address,omitempty"`
}

// Address holds a physical address.
type Address struct {
	AddressLines []string `json:"addressLines,omitempty"`
	Street       string   `json:"street,omitempty"`
	City         string   `json:"city,omitempty"`
	PostalCode   string   `json:"postalCode,omitempty"`
	Country      string   `json:"country,omitempty"`
}

// PaymentRequest is the payload to initiate a payment.
type PaymentRequest struct {
	Amount               float64                      `json:"amount"`
	Currency             string                       `json:"currency"`
	Recipient            Recipient                    `json:"recipient"`
	Reference            string                       `json:"reference,omitempty"`
	IdempotencyKey       string                       `json:"idempotencyKey,omitempty"`
	ContextType          string                       `json:"contextType,omitempty"`
	Type                 string                       `json:"type,omitempty"`
	ScheduledAt          *time.Time                   `json:"scheduledAt,omitempty"`
	PeriodicPayment      *PeriodicPaymentRequest      `json:"periodicPayment,omitempty"`
	InternationalPayment *InternationalPaymentRequest `json:"internationalPayment,omitempty"`
}

// PeriodicPaymentRequest holds scheduling details for a periodic payment.
type PeriodicPaymentRequest struct {
	Frequency         string  `json:"frequency"`
	NumberOfPayments  int     `json:"numberOfPayments,omitempty"`
	NextPaymentAmount float64 `json:"nextPaymentAmount,omitempty"`
	NextPaymentDate   string  `json:"nextPaymentDate,omitempty"`
	FinalPaymentDate  string  `json:"finalPaymentDate,omitempty"`
}

// InternationalPaymentRequest holds additional fields for international payments.
type InternationalPaymentRequest struct {
	CurrencyOfTransfer string            `json:"currencyOfTransfer"`
	InstructedAmount   float64           `json:"instructedAmount,omitempty"`
	InstructedCurrency string            `json:"instructedCurrency,omitempty"`
	ExchangeRateInfo   *ExchangeRateInfo `json:"exchangeRateInfo,omitempty"`
}

// ExchangeRateInfo holds currency exchange rate details.
type ExchangeRateInfo struct {
	UnitCurrency           string  `json:"unitCurrency"`
	ExchangeRate           float64 `json:"exchangeRate"`
	RateType               string  `json:"rateType"`
	ContractIdentification string  `json:"contractIdentification,omitempty"`
}

// Payment represents a payment response.
type Payment struct {
	ID                   string                `json:"id"`
	Status               string                `json:"status"`
	Amount               float64               `json:"amount"`
	Currency             string                `json:"currency"`
	Recipient            Recipient             `json:"recipient"`
	Reference            string                `json:"reference,omitempty"`
	CreatedAt            time.Time             `json:"createdAt"`
	IdempotencyKey       string                `json:"idempotencyKey,omitempty"`
	PaymentLifecycleID   string                `json:"paymentLifecycleId,omitempty"`
	StatusDetails        []PaymentStatusDetail `json:"statusDetails,omitempty"`
	InstitutionPaymentID string                `json:"institutionPaymentId,omitempty"`
	ChargesInfo          []ChargeInfo          `json:"chargesInfo,omitempty"`
}

// PaymentStatusDetail holds a status transition record.
type PaymentStatusDetail struct {
	Status               string    `json:"status"`
	StatusUpdateDateTime time.Time `json:"statusUpdateDateTime"`
}

// ChargeInfo holds charge detail for a payment.
type ChargeInfo struct {
	ChargeBearer string  `json:"chargeBearer"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
}

// BulkPaymentRequest is the payload for bulk payment initiation.
type BulkPaymentRequest struct {
	Payments     []PaymentRequest `json:"payments"`
	OriginatorID string           `json:"originatorId,omitempty"`
	Reference    string           `json:"reference,omitempty"`
}

// BulkPaymentStatusDetail holds status info for a bulk payment.
type BulkPaymentStatusDetail struct {
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BulkPaymentStatus is the response for GET /bulk-payments/{id}.
type BulkPaymentStatus struct {
	ID            string                   `json:"id"`
	ConsentID     string                   `json:"consentId,omitempty"`
	StatusDetails *BulkPaymentStatusDetail `json:"statusDetails,omitempty"`
	CreatedAt     time.Time                `json:"createdAt"`
}

// BulkPayment is the response for bulk payment initiation.
type BulkPayment struct {
	ID       string    `json:"id"`
	Status   string    `json:"status"`
	Payments []Payment `json:"payments,omitempty"`
}

// ─── Consents ─────────────────────────────────────────────────────────────────

// ConsentRequest is the payload to create a new consent.
type ConsentRequest struct {
	InstitutionID     string   `json:"institutionId"`
	ApplicationUserID string   `json:"applicationUserId"`
	Permissions       []string `json:"permissions,omitempty"`
	RedirectURL       string   `json:"redirectUrl,omitempty"`
}

// Consent represents a user's consent for data access.
type Consent struct {
	ID                string     `json:"id"`
	Status            string     `json:"status"`
	InstitutionID     string     `json:"institutionId"`
	ApplicationUserID string     `json:"applicationUserId"`
	ExpiresAt         *time.Time `json:"expiresAt,omitempty"`
	AuthorisationURL  string     `json:"authorisationUrl,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UserID            string     `json:"userId,omitempty"`
	FeatureScopeList  []string   `json:"featureScopeList,omitempty"`
}

// ExchangeCodeRequest exchanges an OAuth2 code for a consent token.
type ExchangeCodeRequest struct {
	ApplicationUserID string `json:"applicationUserId"`
	Code              string `json:"code"`
	RedirectURL       string `json:"redirectUrl,omitempty"`
}

// OneTimeTokenRequest exchanges a one-time token for a consent.
type OneTimeTokenRequest struct {
	OneTimeToken string `json:"oneTimeToken"`
}

// ExtendConsentRequest extends the life of an existing consent.
type ExtendConsentRequest struct {
	TransactionFrom *time.Time `json:"transactionFrom,omitempty"`
	TransactionTo   *time.Time `json:"transactionTo,omitempty"`
	ExpiresAt       *time.Time `json:"expiresAt,omitempty"`
}

// ─── Authorisations ───────────────────────────────────────────────────────────

// AccountAuthorisationRequest creates an account-access authorisation.
type AccountAuthorisationRequest struct {
	InstitutionID     string             `json:"institutionId"`
	ApplicationUserID string             `json:"applicationUserId"`
	Callback          string             `json:"callback,omitempty"`
	RedirectURL       string             `json:"redirectUrl,omitempty"`
	OneTimeToken      bool               `json:"oneTimeToken,omitempty"`
	FeatureScopeList  []string           `json:"featureScopeList,omitempty"`
	UserCredentials   *UserCredentials   `json:"userCredentials,omitempty"`
	ForwardParameters []ForwardParameter `json:"forwardParameters,omitempty"`
}

// PaymentAuthorisationRequest creates a payment authorisation.
type PaymentAuthorisationRequest struct {
	InstitutionID     string           `json:"institutionId"`
	ApplicationUserID string           `json:"applicationUserId"`
	Callback          string           `json:"callback,omitempty"`
	RedirectURL       string           `json:"redirectUrl,omitempty"`
	OneTimeToken      bool             `json:"oneTimeToken,omitempty"`
	PaymentRequest    *PaymentRequest  `json:"paymentRequest,omitempty"`
	UserCredentials   *UserCredentials `json:"userCredentials,omitempty"`
}

// BulkPaymentAuthorisationRequest creates a bulk payment authorisation.
type BulkPaymentAuthorisationRequest struct {
	InstitutionID      string              `json:"institutionId"`
	ApplicationUserID  string              `json:"applicationUserId"`
	Callback           string              `json:"callback,omitempty"`
	BulkPaymentRequest *BulkPaymentRequest `json:"bulkPaymentRequest,omitempty"`
}

// Authorisation is the response for any authorisation request.
type Authorisation struct {
	ID                string     `json:"id"`
	UserID            string     `json:"userId,omitempty"`
	ApplicationUserID string     `json:"applicationUserId,omitempty"`
	Status            string     `json:"status"`
	InstitutionID     string     `json:"institutionId"`
	AuthorisationURL  string     `json:"authorisationUrl,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	ExpiresAt         *time.Time `json:"expiresAt,omitempty"`
	ConsentToken      string     `json:"consentToken,omitempty"`
	FeatureScopeList  []string   `json:"featureScopeList,omitempty"`
	QrCodeURL         string     `json:"qrCodeUrl,omitempty"`
}

// UserCredentials holds SCA / embedded flow credentials.
type UserCredentials struct {
	ID                    string            `json:"id,omitempty"`
	CorpID                string            `json:"corpId,omitempty"`
	Password              string            `json:"password,omitempty"`
	AdditionalFieldValues map[string]string `json:"additionalFieldValues,omitempty"`
}

// ForwardParameter holds forward parameters for the authorisation redirect.
type ForwardParameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PreAuthorisationRequest is the payload for a pre-authorisation.
type PreAuthorisationRequest struct {
	InstitutionID     string `json:"institutionId"`
	ApplicationUserID string `json:"applicationUserId"`
	Callback          string `json:"callback,omitempty"`
	Scope             string `json:"scope,omitempty"`
	OneTimeToken      bool   `json:"oneTimeToken,omitempty"`
}

// ─── Users ────────────────────────────────────────────────────────────────────

// User represents an application user (PSU).
type User struct {
	UUID              string    `json:"uuid"`
	ApplicationUserID string    `json:"applicationUserId"`
	ApplicationUUID   string    `json:"applicationUuid,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	LastModifiedAt    time.Time `json:"lastModifiedAt,omitempty"`
}

// CreateUserRequest is the payload to create a new user.
type CreateUserRequest struct {
	ApplicationUserID string `json:"applicationUserId"`
}

// UpdateUserRequest is the payload to update a user.
type UpdateUserRequest struct {
	ApplicationUserID string `json:"applicationUserId"`
}

// ─── Financial Data extras ────────────────────────────────────────────────────

// Beneficiary represents an account beneficiary (payee).
type Beneficiary struct {
	ID                     string                  `json:"id"`
	Reference              string                  `json:"reference,omitempty"`
	AccountIdentifications []AccountIdentification `json:"accountIdentifications,omitempty"`
	CreditorAgent          *CreditorAgent          `json:"creditorAgent,omitempty"`
}

// CreditorAgent holds the institution of the creditor.
type CreditorAgent struct {
	BICFI string `json:"bicFi,omitempty"`
	Name  string `json:"name,omitempty"`
}

// DirectDebit represents a direct debit mandate.
type DirectDebit struct {
	ID                    string  `json:"id"`
	Reference             string  `json:"reference,omitempty"`
	Status                string  `json:"status,omitempty"`
	Name                  string  `json:"name,omitempty"`
	PreviousPaymentAmount float64 `json:"previousPaymentAmount,omitempty"`
	PreviousPaymentDate   string  `json:"previousPaymentDate,omitempty"`
}

// ScheduledPayment represents a future scheduled payment.
type ScheduledPayment struct {
	ID              string     `json:"id"`
	Reference       string     `json:"reference,omitempty"`
	Amount          float64    `json:"amount"`
	Currency        string     `json:"currency"`
	ScheduledDate   string     `json:"scheduledDate,omitempty"`
	ScheduledType   string     `json:"scheduledType,omitempty"`
	CreditorAccount *Recipient `json:"creditorAccount,omitempty"`
}

// PeriodicPaymentResponse represents a standing order / recurring payment.
type PeriodicPaymentResponse struct {
	ID              string            `json:"id"`
	Reference       string            `json:"reference,omitempty"`
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	Frequency       *FrequencyDetails `json:"frequency,omitempty"`
	CreditorAccount *Recipient        `json:"creditorAccount,omitempty"`
	Status          string            `json:"status,omitempty"`
}

// FrequencyDetails holds recurrence details.
type FrequencyDetails struct {
	Type          string `json:"type"`
	IntervalMonth int    `json:"intervalMonth,omitempty"`
	PointInTime   int    `json:"pointInTime,omitempty"`
}

// Statement represents an account statement.
type Statement struct {
	ID        string `json:"id"`
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	Status    string `json:"status,omitempty"`
	Type      string `json:"type,omitempty"`
}

// Identity holds identity information for KYC/AML.
type Identity struct {
	FullName  string            `json:"fullName,omitempty"`
	Addresses []IdentityAddress `json:"addresses,omitempty"`
	Phones    []Phone           `json:"phones,omitempty"`
	Emails    []Email           `json:"emails,omitempty"`
	BirthDate string            `json:"birthDate,omitempty"`
}

// IdentityAddress holds an address in an identity record.
type IdentityAddress struct {
	AddressType  string   `json:"addressType,omitempty"`
	AddressLines []string `json:"addressLines,omitempty"`
	City         string   `json:"city,omitempty"`
	PostalCode   string   `json:"postalCode,omitempty"`
	Country      string   `json:"country,omitempty"`
}

// Phone holds a phone number entry.
type Phone struct {
	Number string `json:"number"`
	Type   string `json:"type,omitempty"`
}

// Email holds an email entry.
type Email struct {
	Address string `json:"address"`
	Type    string `json:"type,omitempty"`
}

// ─── Variable Recurring Payments ──────────────────────────────────────────────

// VRPAuthorisationRequest creates a VRP consent authorisation.
type VRPAuthorisationRequest struct {
	InstitutionID     string                `json:"institutionId"`
	ApplicationUserID string                `json:"applicationUserId"`
	Callback          string                `json:"callback,omitempty"`
	OneTimeToken      bool                  `json:"oneTimeToken,omitempty"`
	ControlParameters *VRPControlParameters `json:"controlParameters,omitempty"`
	InitiatingParty   *InitiatingParty      `json:"initiatingParty,omitempty"`
}

// VRPControlParameters sets the allowed limits for a VRP.
type VRPControlParameters struct {
	ValidFromDateTime       *time.Time         `json:"validFromDateTime,omitempty"`
	ValidToDateTime         *time.Time         `json:"validToDateTime,omitempty"`
	MaximumIndividualAmount float64            `json:"maximumIndividualAmount,omitempty"`
	Currency                string             `json:"currency,omitempty"`
	PeriodicLimits          []VRPPeriodicLimit `json:"periodicLimits,omitempty"`
}

// VRPPeriodicLimit sets a periodic spend limit.
type VRPPeriodicLimit struct {
	MaximumAmount   float64 `json:"maximumAmount"`
	Currency        string  `json:"currency"`
	PeriodType      string  `json:"periodType"`
	PeriodAlignment string  `json:"periodAlignment"`
}

// InitiatingParty holds details about the party initiating a VRP.
type InitiatingParty struct {
	FinancialID string   `json:"financialId,omitempty"`
	Address     *Address `json:"address,omitempty"`
}

// VRPConsent is the response for a VRP authorisation.
type VRPConsent struct {
	ID                string                `json:"id"`
	Status            string                `json:"status"`
	InstitutionID     string                `json:"institutionId"`
	AuthorisationURL  string                `json:"authorisationUrl,omitempty"`
	CreatedAt         time.Time             `json:"createdAt"`
	ControlParameters *VRPControlParameters `json:"controlParameters,omitempty"`
}

// VRPPaymentRequest creates a Variable Recurring Payment.
type VRPPaymentRequest struct {
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Recipient Recipient `json:"recipient"`
	Reference string    `json:"reference,omitempty"`
}

// VRPPayment is the response for a VRP payment.
type VRPPayment struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"createdAt"`
}

// FundsConfirmationRequest checks if funds are available.
type FundsConfirmationRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// FundsConfirmationResponse is the result of a funds check.
type FundsConfirmationResponse struct {
	FundsAvailable bool `json:"fundsAvailable"`
}

// ─── Notifications / Event Subscriptions ──────────────────────────────────────

// EventSubscription represents a webhook event subscription.
type EventSubscription struct {
	ID              string `json:"id"`
	ApplicationID   string `json:"applicationId,omitempty"`
	EventTypeID     string `json:"eventTypeId"`
	NotificationURL string `json:"notificationUrl"`
	Subscribed      bool   `json:"subscribed,omitempty"`
}

// CreateEventSubscriptionRequest creates a new event subscription.
type CreateEventSubscriptionRequest struct {
	EventTypeID     string `json:"eventTypeId"`
	NotificationURL string `json:"notificationUrl"`
}

// ─── Data Plus (Enrichment) ───────────────────────────────────────────────────

// EnrichmentRequest submits transactions for enrichment.
type EnrichmentRequest struct {
	InstitutionID      string           `json:"institutionId,omitempty"`
	Transactions       []RawTransaction `json:"transactions"`
	CategorisationType string           `json:"categorisationType,omitempty"`
}

// RawTransaction is an un-enriched transaction submitted for analysis.
type RawTransaction struct {
	ID          string  `json:"id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

// EnrichedTransaction is a transaction returned after enrichment.
type EnrichedTransaction struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Category    string    `json:"category,omitempty"`
	Merchant    *Merchant `json:"merchant,omitempty"`
}

// EnrichmentResult holds enrichment job results.
type EnrichmentResult struct {
	JobID        string                `json:"jobId,omitempty"`
	Status       string                `json:"status,omitempty"`
	Transactions []EnrichedTransaction `json:"transactions,omitempty"`
}

// ─── Hosted Pages ─────────────────────────────────────────────────────────────

// HostedConsentRequest creates a hosted consent page session.
type HostedConsentRequest struct {
	InstitutionID       string   `json:"institutionId,omitempty"`
	ApplicationUserID   string   `json:"applicationUserId"`
	RedirectURL         string   `json:"redirectUrl,omitempty"`
	AllowedCountries    []string `json:"allowedCountries,omitempty"`
	AllowedInstitutions []string `json:"allowedInstitutions,omitempty"`
	FeatureScopeList    []string `json:"featureScopeList,omitempty"`
}

// HostedConsentResponse is the response for a hosted consent page creation.
type HostedConsentResponse struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	HostedURL string     `json:"hostedUrl,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

// HostedPaymentRequest creates a hosted payment page session.
type HostedPaymentRequest struct {
	InstitutionID     string          `json:"institutionId,omitempty"`
	ApplicationUserID string          `json:"applicationUserId"`
	RedirectURL       string          `json:"redirectUrl,omitempty"`
	PaymentRequest    *PaymentRequest `json:"paymentRequest"`
	AllowedCountries  []string        `json:"allowedCountries,omitempty"`
}

// HostedPaymentResponse is the response for a hosted payment page creation.
type HostedPaymentResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	HostedURL string    `json:"hostedUrl,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// PayByLinkRequest creates a pay-by-link.
type PayByLinkRequest struct {
	Amount      float64    `json:"amount"`
	Currency    string     `json:"currency"`
	Reference   string     `json:"reference,omitempty"`
	RedirectURL string     `json:"redirectUrl,omitempty"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}

// PayByLinkResponse is the response for pay-by-link creation.
type PayByLinkResponse struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

// ─── Constraints ──────────────────────────────────────────────────────────────

// PaymentConstraint represents payment constraints for an institution.
type PaymentConstraint struct {
	InstitutionID string   `json:"institutionId"`
	PaymentType   string   `json:"paymentType"`
	FeatureScope  []string `json:"featureScope,omitempty"`
	Currencies    []string `json:"currencies,omitempty"`
	MaxAmount     *float64 `json:"maxAmount,omitempty"`
	MinAmount     *float64 `json:"minAmount,omitempty"`
}

// DataConstraint represents data constraints for an institution.
type DataConstraint struct {
	InstitutionID  string   `json:"institutionId"`
	FeatureScope   []string `json:"featureScope,omitempty"`
	MaxDaysHistory int      `json:"maxDaysHistory,omitempty"`
}

// ─── Application Management ───────────────────────────────────────────────────

// ApplicationRequest creates or updates a sub-application.
type ApplicationRequest struct {
	Name                 string   `json:"name"`
	MerchantCategoryCode string   `json:"merchantCategoryCode"`
	CallbackURLs         []string `json:"callbackUrls,omitempty"`
	IsContractPresent    bool     `json:"isContractPresent,omitempty"`
	RootApplicationID    string   `json:"rootApplicationId,omitempty"`
}

// Application is the response representation of an application.
type Application struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	MerchantCategoryCode string   `json:"merchantCategoryCode,omitempty"`
	CallbackURLs         []string `json:"callbackUrls,omitempty"`
	Active               bool     `json:"active,omitempty"`
}

// VRPConfiguration holds VRP configuration for an application.
type VRPConfiguration struct {
	ApplicationID       string   `json:"applicationId"`
	SupportedCurrencies []string `json:"supportedCurrencies,omitempty"`
	MaxAmount           float64  `json:"maxAmount,omitempty"`
	Currency            string   `json:"currency,omitempty"`
}

// ─── Webhooks ─────────────────────────────────────────────────────────────────

// WebhookCategory holds a webhook event category.
type WebhookCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// WebhookEvent represents a registered webhook event.
type WebhookEvent struct {
	ID              string `json:"id"`
	EventTypeID     string `json:"eventTypeId"`
	NotificationURL string `json:"notificationUrl"`
}

// RegisterWebhookRequest registers a new webhook.
type RegisterWebhookRequest struct {
	EventTypeID     string `json:"eventTypeId"`
	NotificationURL string `json:"notificationUrl"`
}

// ResetWebhookSecretResponse holds the new webhook secret.
type ResetWebhookSecretResponse struct {
	Secret string `json:"secret"`
}

// ─── Beneficiaries ────────────────────────────────────────────────────────────

// ApplicationBeneficiary is a beneficiary stored at application level.
type ApplicationBeneficiary struct {
	ID                     string                  `json:"id"`
	Name                   string                  `json:"name"`
	AccountIdentifications []AccountIdentification `json:"accountIdentifications,omitempty"`
}

// CreateBeneficiaryRequest creates a beneficiary.
type CreateBeneficiaryRequest struct {
	Name                   string                  `json:"name"`
	AccountIdentifications []AccountIdentification `json:"accountIdentifications"`
}

// UserBeneficiary is a beneficiary stored at user level.
type UserBeneficiary struct {
	ID                     string                  `json:"id"`
	Name                   string                  `json:"name"`
	Status                 string                  `json:"status,omitempty"`
	AccountIdentifications []AccountIdentification `json:"accountIdentifications,omitempty"`
}

// PatchUserBeneficiaryRequest partially updates a user beneficiary.
type PatchUserBeneficiaryRequest struct {
	Name                   string                  `json:"name,omitempty"`
	AccountIdentifications []AccountIdentification `json:"accountIdentifications,omitempty"`
}

// ─── Pagination helpers ───────────────────────────────────────────────────────

// ConsentListParams holds query filters for listing consents.
// At least one of ApplicationUserIDs, UserUUIDs, or Limit must be set.
type ConsentListParams struct {
	ApplicationUserIDs []string
	UserUUIDs          []string
	InstitutionIDs     []string
	Statuses           []string
	From               string
	Before             string
	Limit              int
	Offset             int
}

// PaginationParams holds query params for paginated list requests.
type PaginationParams struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// TransactionQueryParams extends pagination with transaction-specific date filters.
type TransactionQueryParams struct {
	PaginationParams
	From   string `json:"from,omitempty"`
	Before string `json:"before,omitempty"`
}

// PSUHeaders carries optional PSU (Payment Service User) identifiers.
type PSUHeaders struct {
	PSUID          string
	PSUCorporateID string
	PSUIPAddress   string
}

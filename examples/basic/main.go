// Package main demonstrates a comprehensive Yapily SDK walkthrough.
//
// Run:
//
//	YAPILY_APP_KEY=xxx YAPILY_APP_SECRET=yyy go run ./examples/basic/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/iamkanishka/yapily-client-go/client"
	"github.com/iamkanishka/yapily-client-go/domain"
	"github.com/iamkanishka/yapily-client-go/utils"

	transporthttp "github.com/iamkanishka/yapily-client-go/transport/http"
)

func main() {
	appKey := os.Getenv("YAPILY_APP_KEY")
	appSecret := os.Getenv("YAPILY_APP_SECRET")
	if appKey == "" || appSecret == "" {
		log.Fatal("YAPILY_APP_KEY and YAPILY_APP_SECRET must be set")
	}

	// ── 1. Create the SDK client ──────────────────────────────────────────────
	c, err := client.New(appKey, appSecret,
		client.WithEnvironment(domain.Sandbox),
		client.WithTimeout(30*time.Second),
		client.WithRetry(transporthttp.RetryConfig{
			MaxAttempts: 3,
			BaseDelay:   200 * time.Millisecond,
			MaxDelay:    5 * time.Second,
		}),
		client.WithRateLimit(10, 20),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// ── 2. Application details ────────────────────────────────────────────────
	fmt.Println("→ Fetching application details...")
	app, err := c.Application.GetDetails(ctx)
	if err != nil {
		log.Printf("  application details unavailable: %v", err)
	} else {
		fmt.Printf("✓ Application: %s\n", app.Name)
	}

	// ── 3. List institutions ──────────────────────────────────────────────────
	fmt.Println("\n→ Fetching institutions...")
	institutions, err := c.Institutions.List(ctx)
	if err != nil {
		log.Fatalf("failed to list institutions: %v", err)
	}
	fmt.Printf("✓ Found %d institutions\n", len(institutions))
	for i, inst := range institutions {
		if i >= 3 {
			fmt.Printf("  ... and %d more\n", len(institutions)-3)
			break
		}
		fmt.Printf("  - %-30s (%s)\n", inst.Name, inst.ID)
	}

	// ── 4. Create a user ──────────────────────────────────────────────────────
	fmt.Println("\n→ Creating application user...")
	user, err := c.Users.Create(ctx, &domain.CreateUserRequest{
		ApplicationUserID: fmt.Sprintf("sdk-demo-%d", time.Now().Unix()),
	})
	if err != nil {
		log.Printf("  create user unavailable: %v", err)
	} else {
		fmt.Printf("✓ User: %s\n", user.UUID)
	}

	// ── 5. Initiate account authorisation ─────────────────────────────────────
	fmt.Println("\n→ Creating account authorisation...")
	auth, err := c.Authorisations.CreateAccountAuthorisation(ctx,
		&domain.AccountAuthorisationRequest{
			InstitutionID:     "monzo",
			ApplicationUserID: "demo-user-001",
			FeatureScopeList:  []string{"ACCOUNTS", "TRANSACTIONS"},
			Callback:          "https://yourapp.com/callback",
		},
		&domain.PSUHeaders{PSUIPAddress: "127.0.0.1"},
	)
	if err != nil {
		log.Printf("  authorisation unavailable: %v", err)
	} else {
		fmt.Printf("✓ Authorisation: %s  status: %s\n", auth.ID, auth.Status)
		if auth.AuthorisationURL != "" {
			fmt.Printf("  → Direct user to: %s\n", auth.AuthorisationURL)
		}
	}

	// ── 6. Payment constraints ────────────────────────────────────────────────
	fmt.Println("\n→ Fetching payment constraints...")
	constraints, err := c.Constraints.GetPaymentConstraints(ctx, "monzo", "")
	if err != nil {
		log.Printf("  constraints unavailable: %v", err)
	} else {
		fmt.Printf("✓ Found %d payment constraint rules\n", len(constraints))
	}

	// ── 7. Webhook categories ─────────────────────────────────────────────────
	fmt.Println("\n→ Fetching webhook categories...")
	cats, err := c.Webhooks.GetCategories(ctx)
	if err != nil {
		log.Printf("  webhooks unavailable: %v", err)
	} else {
		fmt.Printf("✓ Found %d webhook categories\n", len(cats))
	}

	// ── 8. Consent-gated flow ─────────────────────────────────────────────────
	// In production: after the user visits auth.AuthorisationURL and is redirected back,
	// exchange the one-time token to get a consent token:
	//
	//   consent, err := c.Consents.ExchangeOneTimeToken(ctx,
	//       &domain.OneTimeTokenRequest{OneTimeToken: "token-from-callback"})
	//   consentToken := consent.ID
	//
	consentToken := os.Getenv("YAPILY_CONSENT_TOKEN")
	if consentToken == "" {
		fmt.Println("\nℹ️  Set YAPILY_CONSENT_TOKEN to demo account/payment flows.")
		fmt.Println("✓ Basic SDK demo complete.")
		return
	}

	// ── 9. List accounts ──────────────────────────────────────────────────────
	fmt.Println("\n→ Fetching accounts...")
	accounts, err := c.Accounts.List(ctx, consentToken)
	if err != nil {
		log.Fatalf("failed to list accounts: %v", err)
	}
	fmt.Printf("✓ Found %d accounts\n", len(accounts))
	for _, acc := range accounts {
		fmt.Printf("  - [%-12s] %-36s  %10.2f %s\n",
			acc.Type, acc.ID, acc.Balance.Amount, acc.Balance.Currency)
	}
	if len(accounts) == 0 {
		fmt.Println("No accounts — done.")
		return
	}
	accountID := accounts[0].ID

	// ── 10. Account balances ──────────────────────────────────────────────────
	fmt.Printf("\n→ Fetching detailed balances for %s...\n", accountID)
	balances, err := c.FinancialData.GetAccountBalances(ctx, consentToken, accountID)
	if err != nil {
		log.Printf("  balances unavailable: %v", err)
	} else {
		for _, b := range balances.Balances {
			fmt.Printf("  %-20s  %10.2f %s\n", b.Type, b.Amount, b.Currency)
		}
	}

	// ── 11. Direct debits ─────────────────────────────────────────────────────
	fmt.Printf("\n→ Fetching direct debits...\n")
	dds, err := c.FinancialData.GetAccountDirectDebits(ctx, consentToken, accountID)
	if err != nil {
		log.Printf("  direct debits unavailable: %v", err)
	} else {
		fmt.Printf("✓ %d direct debits\n", len(dds))
	}

	// ── 12. Transactions ──────────────────────────────────────────────────────
	fmt.Printf("\n→ Fetching transactions for %s...\n", accountID)
	txns, err := c.Transactions.List(ctx, consentToken, accountID, &domain.TransactionQueryParams{
		PaginationParams: domain.PaginationParams{Limit: 10},
	})
	if err != nil {
		log.Fatalf("failed to list transactions: %v", err)
	}
	fmt.Printf("✓ Found %d transactions (page 1)\n", len(txns))
	for i, tx := range txns {
		if i >= 5 {
			break
		}
		fmt.Printf("  [%s]  %10.2f %-4s  %s\n", tx.Date, tx.Amount, tx.Currency, tx.Description)
	}

	// ── 13. Identity ──────────────────────────────────────────────────────────
	fmt.Println("\n→ Fetching identity...")
	identity, err := c.FinancialData.GetIdentity(ctx, consentToken)
	if err != nil {
		log.Printf("  identity unavailable: %v", err)
	} else {
		fmt.Printf("✓ Identity: %s\n", identity.FullName)
	}

	// ── 14. Create payment ────────────────────────────────────────────────────
	fmt.Println("\n→ Creating payment...")
	payment, err := c.Payments.Create(ctx, consentToken, &domain.PaymentRequest{
		Amount:   10.00,
		Currency: "GBP",
		Recipient: domain.Recipient{
			Name: "Jane Smith",
			AccountIdentifications: []domain.AccountIdentification{
				{Type: "SORT_CODE", Identification: "200000"},
				{Type: "ACCOUNT_NUMBER", Identification: "55779911"},
			},
		},
		Reference:      "SDK-Demo",
		IdempotencyKey: utils.IdempotencyKey(""),
	})
	if err != nil {
		log.Printf("  payment unavailable: %v", err)
	} else {
		fmt.Printf("✓ Payment: %s  status: %s\n", payment.ID, payment.Status)
	}

	// ── 15. Data enrichment ───────────────────────────────────────────────────
	if len(txns) > 0 {
		fmt.Println("\n→ Enriching transactions...")
		rawTxns := make([]domain.RawTransaction, 0, len(txns))
		for _, tx := range txns {
			rawTxns = append(rawTxns, domain.RawTransaction{
				ID:          tx.ID,
				Amount:      tx.Amount,
				Currency:    tx.Currency,
				Description: tx.Description,
				Date:        tx.Date,
			})
		}
		enriched, err := c.DataPlus.Enrich(ctx, &domain.EnrichmentRequest{
			Transactions: rawTxns,
		})
		if err != nil {
			log.Printf("  enrichment unavailable: %v", err)
		} else {
			fmt.Printf("✓ Enriched %d transactions\n", len(enriched.Transactions))
		}
	}

	fmt.Println("\n✓ Full SDK demo complete!")
}

package tests

import (
	"context"
	"testing"

	"github.com/iamkanishka/yapily-client-go/auth"
)

func TestBasicAuthProviderGetToken(t *testing.T) {
	p := auth.NewBasicAuthProvider(auth.Config{ApplicationKey: "k", ApplicationSecret: "s"})
	tok, err := p.GetToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tok) < 6 || tok[:6] != "Basic " {
		t.Errorf("expected Basic prefix, got: %s", tok)
	}
}

func TestBasicAuthProviderConsistency(t *testing.T) {
	p := auth.NewBasicAuthProvider(auth.Config{ApplicationKey: "k", ApplicationSecret: "s"})
	tok1, err1 := p.GetToken(context.Background())
	tok2, err2 := p.GetToken(context.Background())
	if err1 != nil || err2 != nil {
		t.Fatalf("errors: %v %v", err1, err2)
	}
	if tok1 != tok2 {
		t.Error("expected identical tokens on repeated calls")
	}
}

func TestBasicAuthProviderInvalidate(t *testing.T) {
	p := auth.NewBasicAuthProvider(auth.Config{ApplicationKey: "k", ApplicationSecret: "s"})
	p.Invalidate()
	tok, err := p.GetToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error after Invalidate: %v", err)
	}
	if tok == "" {
		t.Error("expected non-empty token after Invalidate")
	}
}

func TestBasicAuthProviderEmptyCredentials(t *testing.T) {
	p := auth.NewBasicAuthProvider(auth.Config{})
	// Empty creds produce a structurally valid (but wrong) Basic token.
	// Auth failure happens server-side, not in GetToken itself.
	if _, err := p.GetToken(context.Background()); err != nil {
		t.Fatalf("unexpected error from GetToken with empty creds: %v", err)
	}
}

func TestOAuth2ProviderAliasedToBasic(t *testing.T) {
	p := auth.NewOAuth2Provider(auth.Config{ApplicationKey: "k", ApplicationSecret: "s"})
	tok, err := p.GetToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tok) < 6 || tok[:6] != "Basic " {
		t.Errorf("expected Basic prefix, got: %s", tok)
	}
}

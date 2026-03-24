package tests

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/iamkanishka/yapily-client-go/middleware"
)

func TestVerifyWebhookSignatureValid(t *testing.T) {
	secret := "my-webhook-secret"
	payload := []byte(`{"event":"payment.completed","id":"pay-001"}`)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if err := middleware.VerifyWebhookSignature(payload, secret, sig); err != nil {
		t.Fatalf("expected valid signature, got: %v", err)
	}
}

func TestVerifyWebhookSignatureInvalid(t *testing.T) {
	payload := []byte(`{"event":"payment.completed"}`)

	tests := []struct {
		name string
		sig  string
	}{
		{"wrong secret", "sha256=" + hex.EncodeToString([]byte("not-the-real-sig"))},
		{"missing prefix", hex.EncodeToString([]byte("abc"))},
		{"empty", ""},
		{"bad hex", "sha256=gggg"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := middleware.VerifyWebhookSignature(payload, "secret", tc.sig)
			if err == nil {
				t.Fatalf("expected error for invalid signature")
			}
		})
	}
}

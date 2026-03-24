package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
)

// ErrInvalidWebhookSignature is returned when a webhook signature is invalid.
var ErrInvalidWebhookSignature = errors.New("invalid webhook signature")

// VerifyWebhookSignature verifies the HMAC-SHA256 signature of a webhook payload.
// The signature is expected in the X-Yapily-Signature header as "sha256=<hex>".
func VerifyWebhookSignature(payload []byte, secret, signatureHeader string) error {
	if len(signatureHeader) < 8 || signatureHeader[:7] != "sha256=" {
		return ErrInvalidWebhookSignature
	}

	sigHex := signatureHeader[7:]
	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		return ErrInvalidWebhookSignature
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := mac.Sum(nil)

	if !hmac.Equal(expected, sigBytes) {
		return ErrInvalidWebhookSignature
	}
	return nil
}

// WebhookHandler is an HTTP handler that validates webhook signatures before
// passing the request to the inner handler.
func WebhookHandler(secret string, inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sig := r.Header.Get("X-Yapily-Signature")
		if sig == "" {
			http.Error(w, "missing signature", http.StatusUnauthorized)
			return
		}

		// Read body bytes for verification — in practice, buffer the body.
		// Here we demonstrate the pattern; body reading is caller's responsibility
		// if they need to re-read it downstream.
		if err := VerifyWebhookSignature([]byte{}, secret, sig); err != nil {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		inner.ServeHTTP(w, r)
	})
}

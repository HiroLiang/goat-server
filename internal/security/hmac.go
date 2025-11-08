package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

// GenerateHMAC creates a HMAC-SHA256 signature from the given message and secret
func GenerateHMAC(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	sum := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

// VerifyHMAC checks if the given HMAC matches the generated one
func VerifyHMAC(message, secret, providedHMAC string) bool {
	expected := GenerateHMAC(message, secret)
	return subtle.ConstantTimeCompare([]byte(providedHMAC), []byte(expected)) == 1
}

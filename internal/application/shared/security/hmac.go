package security

type HMACer interface {
	Sign(message string) string
	Verify(message, signature string) bool
}

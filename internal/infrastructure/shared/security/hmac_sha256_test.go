package security

import "testing"

func TestHMAC(t *testing.T) {
	HMACer := NewSHA256HMACer("secret")

	input := "file_string"
	got := HMACer.Sign(input)

	if isValid := HMACer.Verify(input, got); !isValid {
		t.Errorf("VerifyHMAC() = %v, want %v", isValid, true)
	}

	if isValid := HMACer.Verify("wrong", got); isValid {
		t.Errorf("VerifyHMAC() = %v, want %v", isValid, false)
	}
}

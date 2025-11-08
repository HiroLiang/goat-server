package security

import (
	"testing"
)

func TestHMAC(t *testing.T) {
	input := "file_string"
	secret := "secret"
	got := GenerateHMAC(input, secret)

	if isValid := VerifyHMAC(input, secret, got); !isValid {
		t.Errorf("VerifyHMAC() = %v, want %v", isValid, true)
	}

	if isValid := VerifyHMAC("wrong", secret, got); isValid {
		t.Errorf("VerifyHMAC() = %v, want %v", isValid, false)
	}
}

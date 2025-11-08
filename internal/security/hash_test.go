package security

import (
	"testing"
)

func TestArgon2(t *testing.T) {
	input := "password"
	got, err := HashArgon2Base64(input)
	if err != nil {
		t.Errorf("HashArgon2Base64() error = %v", err)
		return
	}

	if isValid := VerifyArgon2Base64(input, got); !isValid {
		t.Errorf("VerifyArgon2Base64() = %v, want %v", isValid, true)
	}

	if isValid := VerifyArgon2Base64("wrong", got); isValid {
		t.Errorf("VerifyArgon2Base64() = %v, want %v", isValid, false)
	}
}

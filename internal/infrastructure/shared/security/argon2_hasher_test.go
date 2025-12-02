package security

import "testing"

func TestArgon2(t *testing.T) {
	hasher := NewArgon2Hasher()

	input := "password"
	got, err := hasher.Hash(input)
	if err != nil {
		t.Errorf("HashArgon2Base64() error = %v", err)
		return
	}

	if isValid := hasher.Verify(input, got); !isValid {
		t.Errorf("VerifyArgon2Base64() = %v, want %v", isValid, true)
	}

	if isValid := hasher.Verify("wrong", got); isValid {
		t.Errorf("VerifyArgon2Base64() = %v, want %v", isValid, false)
	}
}

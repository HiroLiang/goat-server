package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/argon2"
)

// HashArgon2Base64 generates a Base64 hash from a plaintext password
func HashArgon2Base64(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := getArgon2Hash(password, salt)
	encoded := base64.RawStdEncoding.EncodeToString(hash)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	return encodedSalt + ":" + encoded, nil
}

// VerifyArgon2Base64 verifies a Base64 hash against a plaintext password
func VerifyArgon2Base64(password, stored string) bool {
	parts := strings.Split(stored, ":")
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	newHash := getArgon2Hash(password, salt)
	return subtle.ConstantTimeCompare(hash, newHash) == 1
}

func getArgon2Hash(str string, salt []byte) []byte {
	return argon2.IDKey([]byte(str), salt, 1, 64*1024, 4, 32)
}

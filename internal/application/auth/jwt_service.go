package auth

import (
	"context"

	"github.com/HiroLiang/goat-server/internal/domain/auth"
)

type JWTService interface { // TODO Implement JWTService

	// GenerateAuthorizationJWT issues a short-lived JWT after approval
	GenerateAuthorizationJWT(ctx context.Context, payload auth.AuthorizationJWTPayload) (token string, err error)

	// ValidateJWT verifies signature, expiration, and returns payload
	ValidateJWT(ctx context.Context, token string) (*auth.AuthorizationJWTPayload, error)

	// RevokeJWT stores jti in blacklist (Redis)
	RevokeJWT(ctx context.Context, jti string, exp int64) error

	// IsJWTRevoked checks whether token id is in the blacklist
	IsJWTRevoked(ctx context.Context, jti string) (bool, error)
}

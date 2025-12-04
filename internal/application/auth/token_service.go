package auth

import (
	"context"

	"github.com/HiroLiang/goat-server/internal/domain/auth"
)

type TokenService interface {
	Generate(ctx context.Context, params auth.CreateSessionParams) (string, error)
	Refresh(ctx context.Context, token string) (string, error)
	Validate(ctx context.Context, token string) (*auth.Session, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID string) error
	ListUserSessions(ctx context.Context, userID string) ([]*auth.Session, error)
}

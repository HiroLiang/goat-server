package session

import (
	"context"
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/auth"
)

type Store interface {
	Get(ctx context.Context, token string) (*auth.Session, error)
	Set(ctx context.Context, token string, session *auth.Session, ttl time.Duration) error
	Delete(ctx context.Context, token string) error
	Refresh(ctx context.Context, token string, ttl time.Duration) error

	AddUserSession(ctx context.Context, userID, token string, ttl time.Duration) error
	RemoveUserSession(ctx context.Context, userID, token string) error
	ListUserSessions(ctx context.Context, userID string) ([]string, error)
	DeleteAllUserSessions(ctx context.Context, userID string) error
	CleanExpiredUserSessions(ctx context.Context, userID string) error
}

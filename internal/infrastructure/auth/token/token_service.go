package token

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	iAuth "github.com/HiroLiang/goat-server/internal/application/shared/auth"
	"github.com/HiroLiang/goat-server/internal/domain/auth"
	"github.com/HiroLiang/goat-server/internal/infrastructure/auth/session"
)

type AuthTokenService struct {
	store    session.Store
	tokenTTL time.Duration
}

var _ iAuth.TokenService = (*AuthTokenService)(nil)

func NewAuthTokenService(store session.Store, ttl time.Duration) *AuthTokenService {
	return &AuthTokenService{
		store:    store,
		tokenTTL: ttl,
	}
}

// Generate creates a new session token and stores it.
func (s *AuthTokenService) Generate(ctx context.Context, params auth.CreateSessionParams) (string, error) {
	token, err := s.generateSecureToken(32)
	if err != nil {
		return "", auth.ErrGenerateToken
	}

	sess := &auth.Session{
		UserID:    params.UserID,
		IP:        params.IP,
		UserAgent: params.UserAgent,
		CreatedAt: time.Now(),
	}

	if err := s.store.Set(ctx, token, sess, s.tokenTTL); err != nil {
		return "", err
	}

	return token, nil
}

// Validate checks if the session exists and auto-refreshes TTL (sliding session).
func (s *AuthTokenService) Validate(ctx context.Context, token string) (*auth.Session, error) {
	sess, err := s.store.Get(ctx, token)
	if err != nil {
		return nil, err
	}

	if err := s.store.Refresh(ctx, token, s.tokenTTL); err != nil {
		return nil, auth.ErrRefreshToken
	}

	// clean expired sessions in the background
	go func() {
		_ = s.store.CleanExpiredUserSessions(context.Background(), sess.UserID)
	}()

	return sess, nil
}

// Refresh creates a brand-new session token for the same user.
func (s *AuthTokenService) Refresh(ctx context.Context, token string) (string, error) {
	sess, err := s.store.Get(ctx, token)
	if err != nil {
		return "", err
	}

	if err := s.store.Delete(ctx, token); err != nil {
		return "", err
	}

	return s.Generate(ctx, auth.CreateSessionParams{
		UserID:    sess.UserID,
		IP:        sess.IP,
		UserAgent: sess.UserAgent,
	})
}

// Revoke deletes a session token.
func (s *AuthTokenService) Revoke(ctx context.Context, token string) error {
	return s.store.Delete(ctx, token)
}

// ListUserSessions returns all active sessions of the user.
func (s *AuthTokenService) ListUserSessions(ctx context.Context, userID string) ([]*auth.Session, error) {
	tokens, err := s.store.ListUserSessions(ctx, userID)
	if err != nil {
		return nil, err
	}

	sessions := make([]*auth.Session, 0, len(tokens))
	for _, token := range tokens {
		sess, err := s.store.Get(ctx, token)
		if errors.Is(err, auth.ErrSessionNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, sess)
	}

	return sessions, nil
}

// RevokeAllForUser deletes all sessions for a user.
func (s *AuthTokenService) RevokeAllForUser(ctx context.Context, userID string) error {
	return s.store.DeleteAllUserSessions(ctx, userID)
}

// --- Helpers ---

func (s *AuthTokenService) generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

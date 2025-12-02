package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/auth"
	"github.com/redis/go-redis/v9"
)

type AuthTokenService struct {
	redis    *redis.Client
	tokenTTL time.Duration
}

func NewAuthTokenService(redis *redis.Client, ttl time.Duration) *AuthTokenService {
	return &AuthTokenService{
		redis:    redis,
		tokenTTL: ttl,
	}
}

// Generate creates a new session token and stores it in Redis.
func (s *AuthTokenService) Generate(ctx context.Context, params auth.CreateSessionParams) (string, error) {

	// generate secure random token
	token, err := s.generateSecureToken(32)
	if err != nil {
		return "", auth.ErrGenerateToken
	}

	// create session data
	session := &auth.Session{
		UserID:    params.UserID,
		IP:        params.IP,
		UserAgent: params.UserAgent,
		CreatedAt: time.Now(),
	}

	// store session data
	if err := s.storeSession(ctx, token, session); err != nil {
		return "", err
	}

	return token, nil
}

// Validate checks if the session exists and auto-refreshes TTL (sliding session).
func (s *AuthTokenService) Validate(ctx context.Context, token string) (*auth.Session, error) {
	key := "session:" + token

	// get session data
	session, err := s.getSession(ctx, token)
	if err != nil {
		return nil, auth.ErrSessionNotFound
	}

	// sliding session TTL
	if err := s.redis.Expire(ctx, key, s.tokenTTL).Err(); err != nil {
		return nil, auth.ErrRefreshToken
	}

	// clean unused session tokens
	if err := s.cleanUnusedUserSessions(ctx, session.UserID); err != nil {
		return nil, err
	}

	return session, nil
}

// Refresh creates a brand-new session token for the same user.
func (s *AuthTokenService) Refresh(ctx context.Context, token string) (string, error) {

	// get session data
	session, err := s.getSession(ctx, token)
	if err != nil {
		return "", err
	}

	// revoke old token
	if err := s.revokeSession(ctx, token); err != nil {
		return "", err
	}

	// create a new token
	return s.Generate(ctx, auth.CreateSessionParams{
		UserID:    session.UserID,
		IP:        session.IP,
		UserAgent: session.UserAgent,
	})
}

// Revoke deletes a session token.
func (s *AuthTokenService) Revoke(ctx context.Context, token string) error {
	return s.revokeSession(ctx, token)
}

// ListUserSessions returns all active sessions of the user.
func (s *AuthTokenService) ListUserSessions(ctx context.Context, userID string) ([]*auth.Session, error) {
	indexKey := "user_sessions:" + userID

	// get all tokens for the user
	tokens, err := s.redis.SMembers(ctx, indexKey).Result()
	if err != nil {
		return nil, fmt.Errorf("redis smembers: %w", err)
	}

	// get session data for each token
	sessions := make([]*auth.Session, 0, len(tokens))

	// ignore errors for tokens that don't exist'
	for _, token := range tokens {
		session, err := s.getSession(ctx, token)
		if errors.Is(err, auth.ErrSessionNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *AuthTokenService) RevokeAllForUser(ctx context.Context, userID string) error {
	indexKey := "user_sessions:" + userID

	// get all tokens for the user
	tokens, err := s.redis.SMembers(ctx, indexKey).Result()
	if err != nil {
		return fmt.Errorf("redis smembers: %w", err)
	}

	// delete all session tokens
	for _, token := range tokens {
		if err := s.revokeSession(ctx, token); err != nil {
			return err
		}
	}

	// delete the index
	if err := s.redis.Del(ctx, indexKey).Err(); err != nil {
		return err
	}

	return nil
}

// --- Helpers ---

func (s *AuthTokenService) generateSecureToken(length int) (string, error) {

	// Generate a secure random token
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("secure token generation failed: %w", err)
	}

	// Encode the token in base64
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (s *AuthTokenService) storeSession(ctx context.Context, token string, session *auth.Session) error {
	key := "session:" + token

	b, _ := json.Marshal(session)

	// store session data in Redis
	s.redis.Set(ctx, key, b, s.tokenTTL)

	// add token to user index
	indexKey := "user_sessions:" + session.UserID
	s.redis.SAdd(ctx, indexKey, token)

	return nil
}

func (s *AuthTokenService) getSession(ctx context.Context, token string) (*auth.Session, error) {
	key := "session:" + token

	// get session data from Redis
	b, err := s.redis.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, auth.ErrSessionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("redis get session: %w", err)
	}

	// unmarshal session data
	var session auth.Session
	if err := json.Unmarshal(b, &session); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	return &session, nil
}

func (s *AuthTokenService) revokeSession(ctx context.Context, token string) error {
	session, _ := s.getSession(ctx, token)

	// delete session data
	s.redis.Del(ctx, "session:"+token)

	// remove token from user index (idempotent)
	if session != nil {
		indexKey := "user_sessions:" + session.UserID
		s.redis.SRem(ctx, indexKey, token)
	}

	return nil
}

func (s *AuthTokenService) cleanUnusedUserSessions(ctx context.Context, userID string) error {
	indexKey := "user_sessions:" + userID

	// get all tokens for the user
	tokens, err := s.redis.SMembers(ctx, indexKey).Result()
	if err != nil {
		return fmt.Errorf("redis smembers: %w", err)
	}

	// delete unused session tokens
	for _, token := range tokens {
		_, err := s.getSession(ctx, token)
		if err != nil {
			s.redis.SRem(ctx, indexKey, token)
		}
	}

	return nil
}

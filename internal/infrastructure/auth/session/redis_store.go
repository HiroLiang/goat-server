package session

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/auth"
	"github.com/HiroLiang/goat-server/internal/infrastructure/cache"
	"github.com/redis/go-redis/v9"
)

type RedisSessionStore struct {
	cache cache.Cache
	redis *redis.Client
}

var _ Store = (*RedisSessionStore)(nil)

func NewRedisSessionStore(cache cache.Cache, redis *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{
		cache: cache,
		redis: redis,
	}
}

func (s *RedisSessionStore) Get(ctx context.Context, token string) (*auth.Session, error) {
	key := "session:" + token

	b, ok, err := s.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, auth.ErrSessionNotFound
	}

	var session auth.Session
	if err := json.Unmarshal(b, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *RedisSessionStore) Set(ctx context.Context, token string, session *auth.Session, ttl time.Duration) error {
	key := "session:" + token

	b, err := json.Marshal(session)
	if err != nil {
		return err
	}

	if err := s.cache.Set(ctx, key, b, ttl); err != nil {
		return err
	}

	return s.AddUserSession(ctx, session.UserID, token, ttl)
}

func (s *RedisSessionStore) Delete(ctx context.Context, token string) error {
	session, err := s.Get(ctx, token)
	if err != nil && !errors.Is(err, auth.ErrSessionNotFound) {
		return err
	}

	if err := s.cache.Delete(ctx, "session:"+token); err != nil {
		return err
	}

	if session != nil {
		return s.RemoveUserSession(ctx, session.UserID, token)
	}

	return nil
}

func (s *RedisSessionStore) Refresh(ctx context.Context, token string, ttl time.Duration) error {
	session, err := s.Get(ctx, token)
	if err != nil {
		return err
	}

	pipe := s.redis.Pipeline()
	pipe.Expire(ctx, "session:"+token, ttl)
	pipe.Expire(ctx, "user_sessions:"+session.UserID, ttl)
	_, err = pipe.Exec(ctx)

	return err
}

func (s *RedisSessionStore) AddUserSession(ctx context.Context, userID, token string, ttl time.Duration) error {
	key := "user_sessions:" + userID

	pipe := s.redis.Pipeline()
	pipe.SAdd(ctx, key, token)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)

	return err
}

func (s *RedisSessionStore) RemoveUserSession(ctx context.Context, userID, token string) error {
	key := "user_sessions:" + userID
	return s.redis.SRem(ctx, key, token).Err()
}

func (s *RedisSessionStore) ListUserSessions(ctx context.Context, userID string) ([]string, error) {
	key := "user_sessions:" + userID
	return s.redis.SMembers(ctx, key).Result()
}

func (s *RedisSessionStore) DeleteAllUserSessions(ctx context.Context, userID string) error {
	tokens, err := s.ListUserSessions(ctx, userID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	// batch delete all sessions
	keys := make([]string, len(tokens))
	for i, token := range tokens {
		keys[i] = "session:" + token
	}

	pipe := s.redis.Pipeline()
	pipe.Del(ctx, keys...)
	pipe.Del(ctx, "user_sessions:"+userID)
	_, err = pipe.Exec(ctx)

	return err
}

func (s *RedisSessionStore) CleanExpiredUserSessions(ctx context.Context, userID string) error {
	tokens, err := s.ListUserSessions(ctx, userID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		return nil
	}

	pipe := s.redis.Pipeline()
	cmd := make([]*redis.IntCmd, len(tokens))
	for i, token := range tokens {
		cmd[i] = pipe.Exists(ctx, "session:"+token)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	// collect expired tokens
	var expiredTokens []interface{}
	for i, cmd := range cmd {
		if cmd.Val() == 0 {
			expiredTokens = append(expiredTokens, tokens[i])
		}
	}

	// batch delete expired tokens
	if len(expiredTokens) > 0 {
		return s.redis.SRem(ctx, "user_sessions:"+userID, expiredTokens...).Err()
	}

	return nil
}

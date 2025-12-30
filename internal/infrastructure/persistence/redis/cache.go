package redis

import (
	"context"
	"errors"
	"time"

	"github.com/HiroLiang/goat-server/internal/infrastructure/cache"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
}

var _ cache.Cache = (*RedisCache)(nil)

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client}
}

func (r RedisCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	b, err := r.Client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return b, true, nil
}

func (r RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r RedisCache) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

package security

import (
	"context"
	"time"
)

type RateLimitRepository interface {

	// Increment increments the rate limit counter for a key
	Increment(ctx context.Context, key string, window time.Duration) (count int64, err error)

	// Get returns the rate limit counter for a key
	Get(ctx context.Context, key string) (count int64, err error)

	// Reset resets the rate limit counter for a key
	Reset(ctx context.Context, key string) error

	// IncrementSliding increments the rate limit counter for a key with the sliding window
	IncrementSliding(ctx context.Context, key string, window time.Duration, now time.Time) (count int64, err error)
}

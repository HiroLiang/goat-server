package security

import "context"

// RateLimiter Generic rate limiter
type RateLimiter interface {
	CheckIP(ctx context.Context, ip string) error
	CheckGlobal(ctx context.Context) error
}

// LoginRateLimiter THe rate limiter for login attempts
type LoginRateLimiter interface {
	CheckLoginAttempt(ctx context.Context, ip, email string) error
	RecordLoginAttempt(ctx context.Context, ip, email string, success bool) error
	ReleaseLock(ctx context.Context, email string) error
}

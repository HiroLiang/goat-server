package mock

import (
	"context"

	"github.com/HiroLiang/goat-server/internal/application/shared/security"
)

type RateLimiter struct{}

func MockRateLimiter() security.RateLimiter {
	return RateLimiter{}
}

var _ security.RateLimiter = (*RateLimiter)(nil)

func (m RateLimiter) CheckGlobal(_ context.Context) error {
	return nil
}

func (m RateLimiter) CheckIP(_ context.Context, _ string) error {
	return nil
}

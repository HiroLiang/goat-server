package security

import "time"

type RateLimitPolicy struct {
	Limit  int64
	Window time.Duration
}

package auth

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/user"
)

type TokenRecord struct {
	UserID    user.ID
	ExpiresAt time.Time
}

package role

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/user"
)

type Role struct {
	ID        ID
	Type      Type
	Creator   user.ID
	CreateAt  time.Time
	UpdatedAt time.Time
}

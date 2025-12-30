package role

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/role"
	"github.com/HiroLiang/goat-server/internal/domain/user"
)

type RoleRecord struct {
	ID        role.ID   `db:"id"`
	Type      role.Type `db:"type"`
	Creator   user.ID   `db:"creator"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

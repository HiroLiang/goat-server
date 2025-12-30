package userrole

import (
	"github.com/HiroLiang/goat-server/internal/domain/role"
	"github.com/HiroLiang/goat-server/internal/domain/user"
)

type UserRole struct {
	ID   user.ID
	Role role.Role
}

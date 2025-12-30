package auth

import "github.com/HiroLiang/goat-server/internal/domain/user"

type RoleService interface {
	HasRole(userID user.ID, roleID string) bool
}

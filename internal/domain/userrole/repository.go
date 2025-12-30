package userrole

import (
	"context"

	"github.com/HiroLiang/goat-server/internal/domain/role"
	"github.com/HiroLiang/goat-server/internal/domain/user"
)

type Repository interface {
	FindRolesByUser(ctx context.Context, userID user.ID) ([]*role.Role, error)
	Exists(ctx context.Context, userID user.ID, role role.Type) bool
	Assign(ctx context.Context, userID user.ID, role role.Type) error
	Revoke(ctx context.Context, userID user.ID, role role.Type) error
}

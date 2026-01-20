package mock

import (
	"context"

	"github.com/HiroLiang/goat-server/internal/domain/role"
	"github.com/HiroLiang/goat-server/internal/domain/user"
	"github.com/HiroLiang/goat-server/internal/domain/userrole"
)

type UserRoleRepo struct{}

func MockUserRoleRepo() *UserRoleRepo {
	return &UserRoleRepo{}
}

var _ userrole.Repository = (*UserRoleRepo)(nil)

func (u UserRoleRepo) FindRolesByUser(ctx context.Context, userID user.ID) ([]*role.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRoleRepo) Exists(ctx context.Context, userID user.ID, role role.Type) bool {
	//TODO implement me
	panic("implement me")
}

func (u UserRoleRepo) Assign(ctx context.Context, userID user.ID, role role.Type) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRoleRepo) Revoke(ctx context.Context, userID user.ID, role role.Type) error {
	//TODO implement me
	panic("implement me")
}

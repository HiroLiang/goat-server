package user

import (
	"context"

	"github.com/HiroLiang/goat-server/internal/domain/role"
	"github.com/HiroLiang/goat-server/internal/domain/user"
	"github.com/HiroLiang/goat-server/internal/domain/userrole"
	"github.com/stretchr/testify/mock"
)

type MockUserRoleRepo struct {
	mock.Mock
}

var _ userrole.Repository = (*MockUserRoleRepo)(nil)

func (m *MockUserRoleRepo) FindRolesByUser(ctx context.Context, userID user.ID) ([]*role.Role, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*role.Role), args.Error(1)
}

func (m *MockUserRoleRepo) Exists(ctx context.Context, userID user.ID, role role.Type) bool {
	args := m.Called(ctx, userID, role)
	return args.Bool(0)
}

func (m *MockUserRoleRepo) Assign(ctx context.Context, userID user.ID, role role.Type) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *MockUserRoleRepo) Revoke(ctx context.Context, userID user.ID, role role.Type) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

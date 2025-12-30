package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HiroLiang/goat-server/internal/application/shared"
	"github.com/HiroLiang/goat-server/internal/domain/role"
	"github.com/HiroLiang/goat-server/internal/domain/user"
	"github.com/HiroLiang/goat-server/internal/domain/userrole"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAssignRoleToUser_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	repo := new(MockUserRoleRepo)
	repo.
		On("Assign", mock.Anything, user.ID(1), role.Admin).
		Return(nil)

	uc := &UseCase{userRoleRepo: repo}

	input := shared.UseCaseInput[AssignRoleInput]{
		Data: AssignRoleInput{
			UserID: user.ID(1),
			Role:   role.Admin,
		},
	}

	err := uc.AssignRoleToUser(ctx, input)

	assert.NoError(t, err)
}

func TestAssignRoleToUser_Idempotent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	repo := new(MockUserRoleRepo)
	repo.
		On("Assign", mock.Anything, user.ID(1), role.Admin).
		Return(userrole.ErrUserRoleAlreadyAssigned)

	uc := UseCase{userRoleRepo: repo}

	input := shared.UseCaseInput[AssignRoleInput]{
		Data: AssignRoleInput{
			UserID: user.ID(1),
			Role:   role.Admin,
		},
	}

	err := uc.AssignRoleToUser(ctx, input)

	assert.ErrorIs(t, err, userrole.ErrUserRoleAlreadyAssigned)
}

func TestAssignRoleToUser_RepoError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	repo := new(MockUserRoleRepo)
	repo.
		On("Assign", mock.Anything, user.ID(1), role.Admin).
		Return(errors.New("db error"))

	uc := &UseCase{userRoleRepo: repo}

	input := shared.UseCaseInput[AssignRoleInput]{
		Data: AssignRoleInput{
			UserID: user.ID(1),
			Role:   role.Admin,
		},
	}

	err := uc.AssignRoleToUser(ctx, input)

	assert.ErrorIs(t, err, userrole.ErrAssignFailed)
}

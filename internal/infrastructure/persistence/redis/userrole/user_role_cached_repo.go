package userrole

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/role"
	"github.com/HiroLiang/goat-server/internal/domain/user"
	"github.com/HiroLiang/goat-server/internal/domain/userrole"
	"github.com/HiroLiang/goat-server/internal/infrastructure/cache"
)

func userRolesCacheKey(userID user.ID) string {
	return fmt.Sprintf("user_roles:%d:v1", userID)
}

type UserRoleCachedRepo struct {
	client       cache.Cache
	userRoleRepo userrole.Repository
}

var _ userrole.Repository = (*UserRoleCachedRepo)(nil)

func NewUserRoleCachedRepo(client cache.Cache, userRoleRepo userrole.Repository) *UserRoleCachedRepo {
	return &UserRoleCachedRepo{client: client, userRoleRepo: userRoleRepo}
}

func (u UserRoleCachedRepo) FindRolesByUser(ctx context.Context, userID user.ID) ([]*role.Role, error) {
	key := userRolesCacheKey(userID)

	// 1. try cache
	if b, ok, err := u.client.Get(ctx, key); err == nil && ok {
		roles, err := decodeRoles(b)
		if err == nil {
			return roles, nil
		}
		_ = u.client.Delete(ctx, key)
	}

	// 2. fallback DB
	roles, err := u.userRoleRepo.FindRolesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. set cache
	_ = u.client.Set(ctx, key, encodeRoles(roles), time.Minute)

	return roles, nil
}

func (u UserRoleCachedRepo) Exists(ctx context.Context, userID user.ID, role role.Type) bool {
	return u.userRoleRepo.Exists(ctx, userID, role)
}

func (u UserRoleCachedRepo) Assign(ctx context.Context, userID user.ID, role role.Type) error {
	return u.userRoleRepo.Assign(ctx, userID, role)
}

func (u UserRoleCachedRepo) Revoke(ctx context.Context, userID user.ID, role role.Type) error {
	return u.userRoleRepo.Revoke(ctx, userID, role)
}

func encodeRoles(roles []*role.Role) []byte {
	b, _ := json.Marshal(roles)
	return b
}

func decodeRoles(b []byte) ([]*role.Role, error) {
	var roles []*role.Role
	err := json.Unmarshal(b, &roles)
	return roles, err
}

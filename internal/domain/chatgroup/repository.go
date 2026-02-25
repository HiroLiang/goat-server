package chatgroup

import (
	"context"

	"github.com/HiroLiang/goat-server/internal/domain/user"
)

type Repository interface {
	FindByID(ctx context.Context, id ID) (*ChatGroup, error)
	FindByCreator(ctx context.Context, creatorID user.ID) ([]*ChatGroup, error)
	Create(ctx context.Context, group *ChatGroup) error
	Update(ctx context.Context, group *ChatGroup) error
	SoftDelete(ctx context.Context, id ID) error
}

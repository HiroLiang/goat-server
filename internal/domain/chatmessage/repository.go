package chatmessage

import (
	"context"

	"github.com/HiroLiang/goat-server/internal/domain/chatgroup"
	"github.com/HiroLiang/goat-server/internal/domain/participant"
)

type Repository interface {
	FindByID(ctx context.Context, id ID) (*ChatMessage, error)
	FindByGroup(ctx context.Context, groupID chatgroup.ID, limit, offset uint64) ([]*ChatMessage, error)
	FindBySender(ctx context.Context, senderID participant.ID) ([]*ChatMessage, error)
	Create(ctx context.Context, message *ChatMessage) error
	Update(ctx context.Context, message *ChatMessage) error
	SoftDelete(ctx context.Context, id ID) error
}

package chatmember

import (
	"context"

	"github.com/HiroLiang/goat-server/internal/domain/chatgroup"
	"github.com/HiroLiang/goat-server/internal/domain/participant"
)

type Repository interface {
	FindByID(ctx context.Context, id ID) (*ChatMember, error)
	FindByGroupAndParticipant(ctx context.Context, groupID chatgroup.ID, participantID participant.ID) (*ChatMember, error)
	FindByGroup(ctx context.Context, groupID chatgroup.ID) ([]*ChatMember, error)
	FindByParticipant(ctx context.Context, participantID participant.ID) ([]*ChatMember, error)
	Add(ctx context.Context, member *ChatMember) error
	Update(ctx context.Context, member *ChatMember) error
	Remove(ctx context.Context, groupID chatgroup.ID, participantID participant.ID) error
}

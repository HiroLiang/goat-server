package participant

import (
	"context"

	"github.com/HiroLiang/goat-server/internal/domain/agent"
	"github.com/HiroLiang/goat-server/internal/domain/user"
)

type Repository interface {
	FindByID(ctx context.Context, id ID) (*Participant, error)
	FindByUserID(ctx context.Context, userID user.ID) (*Participant, error)
	FindByAgentID(ctx context.Context, agentID agent.ID) (*Participant, error)
	FindSystem(ctx context.Context) (*Participant, error)
	Create(ctx context.Context, p *Participant) error
}

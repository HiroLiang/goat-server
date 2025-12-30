package agent

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/agent"
	"github.com/HiroLiang/goat-server/internal/domain/user"
)

type AgentRecord struct {
	ID        agent.ID     `db:"id"`
	Name      string       `db:"name"`
	Type      agent.Type   `db:"type"`
	Status    agent.Status `db:"status"`
	Engine    agent.Engine `db:"engine"`
	CreatedAt time.Time    `db:"created_at"`
	CreatedBy user.ID      `db:"created_by"`
	UpdatedAt time.Time    `db:"updated_at"`
	UpdatedBy user.ID      `db:"updated_by"`
}

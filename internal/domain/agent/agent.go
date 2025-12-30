package agent

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/user"
)

type Agent struct {
	ID        ID
	Name      string
	Type      Type
	Status    Status
	Engine    Engine
	CreatedAt time.Time
	CreatedBy user.ID
	UpdatedAt time.Time
	UpdatedBy user.ID
}

func (a Agent) InUse() bool {
	return a.Status == Available
}

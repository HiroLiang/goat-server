package agent

import "github.com/HiroLiang/goat-server/internal/domain/agent"

type QueryAvailableAgentsOutput struct {
	Name     string
	Provider string
	Status   agent.Status
}

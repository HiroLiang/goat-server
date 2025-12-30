package agent

import (
	"context"

	"github.com/HiroLiang/goat-server/internal/application/shared"
	"github.com/HiroLiang/goat-server/internal/domain/agent"
)

type UseCase struct {
	agentRepo agent.Repository
}

func NewUseCase(agentRepo agent.Repository) *UseCase {
	return &UseCase{agentRepo}
}

func (u UseCase) FindAvailableAgents(
	ctx context.Context,
	_ shared.UseCaseInput[QueryAvailableAgentsInput],
) ([]QueryAvailableAgentsOutput, error) {
	domains, err := u.agentRepo.FindAllByStatus(ctx, agent.Available)
	if err != nil {
		return nil, err
	}

	outputs := make([]QueryAvailableAgentsOutput, len(domains))
	for i, domain := range domains {
		outputs[i] = QueryAvailableAgentsOutput{
			Name:     domain.Name,
			Provider: "",
			Status:   domain.Status,
		}
	}

	return outputs, nil
}

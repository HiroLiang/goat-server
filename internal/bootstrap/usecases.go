package bootstrap

import (
	"github.com/HiroLiang/goat-server/internal/application/agent"
	"github.com/HiroLiang/goat-server/internal/application/user"
)

type UseCases struct {
	UserUseCase  *user.UseCase
	AgentUseCase *agent.UseCase
}

func BuildUseCases(deps *Dependencies) *UseCases {
	return &UseCases{
		UserUseCase:  user.NewUseCase(deps.UserRepo, deps.UserRoleRepo, deps.Hasher, deps.TokenService),
		AgentUseCase: agent.NewUseCase(deps.AgentRepo),
	}
}

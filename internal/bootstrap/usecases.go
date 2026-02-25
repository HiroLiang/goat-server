package bootstrap

import (
	"github.com/HiroLiang/goat-server/internal/application/agent"
	"github.com/HiroLiang/goat-server/internal/application/chat"
	"github.com/HiroLiang/goat-server/internal/application/user"
)

type UseCases struct {
	UserUseCase  *user.UseCase
	AgentUseCase *agent.UseCase
	ChatUseCase  *chat.UseCase
}

func BuildUseCases(deps *Dependencies) *UseCases {
	return &UseCases{
		UserUseCase:  user.NewUseCase(deps.UserRepo, deps.UserRoleRepo, deps.Hasher, deps.TokenService),
		AgentUseCase: agent.NewUseCase(deps.AgentRepo),
		ChatUseCase: chat.NewUseCase(
			deps.ParticipantRepo,
			deps.ChatGroupRepo,
			deps.ChatMemberRepo,
			deps.ChatMessageRepo,
		),
	}
}

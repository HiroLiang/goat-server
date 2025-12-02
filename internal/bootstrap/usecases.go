package bootstrap

import "github.com/HiroLiang/goat-server/internal/application/user"

type UseCases struct {
	UserUseCase *user.UseCase
}

func BuildUseCases(deps *Dependencies) *UseCases {
	return &UseCases{
		UserUseCase: user.NewUseCase(deps.UserRepo, deps.Hasher, deps.TokenService),
	}
}

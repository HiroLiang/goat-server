package user

import (
	"context"
	"strconv"

	"github.com/HiroLiang/goat-server/internal/application/shared"
	"github.com/HiroLiang/goat-server/internal/application/shared/auth"
	"github.com/HiroLiang/goat-server/internal/application/shared/security"
	session "github.com/HiroLiang/goat-server/internal/domain/auth"
	"github.com/HiroLiang/goat-server/internal/domain/user"
	"github.com/HiroLiang/goat-server/internal/shared/timeutil"
)

type UseCase struct {
	userRepo     user.Repository
	hasher       security.Hasher
	tokenService auth.TokenService
}

func NewUseCase(repo user.Repository, hasher security.Hasher, tokenService auth.TokenService) *UseCase {
	return &UseCase{
		userRepo:     repo,
		hasher:       hasher,
		tokenService: tokenService,
	}
}

// Register User register
func (u *UseCase) Register(ctx context.Context, input shared.UseCaseInput[RegisterInput]) error {
	hash, err := u.hasher.Hash(input.Data.Password)
	if err != nil {
		return user.ErrInvalidPassword
	}

	email, err := user.NewEmail(input.Data.Email)
	if err != nil {
		return user.ErrInvalidEmail
	}

	newUser := user.NewUser(
		input.Data.Name,
		email,
		hash,
		input.Base.Request.IP,
	)

	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return err
	}

	return nil
}

// Login User login
func (u *UseCase) Login(ctx context.Context, input shared.UseCaseInput[LoginInput]) (LoginOutput, error) {

	// Check email and build vo
	email, err := user.NewEmail(input.Data.Email)
	if err != nil {
		return LoginOutput{}, user.ErrInvalidEmail
	}

	// Check is user exists
	currentUser, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return LoginOutput{}, user.ErrUserNotFound
	} else if !currentUser.IsValid() {
		return LoginOutput{}, user.ErrInvalidUser
	}

	// Check password
	if !u.hasher.Verify(input.Data.Password, currentUser.Password) {
		return LoginOutput{}, user.ErrInvalidPassword
	}

	// Generate auth token and store in redis
	authToken, err := u.tokenService.Generate(ctx, session.CreateSessionParams{
		UserID:    strconv.FormatInt(int64(currentUser.ID), 10),
		IP:        input.Base.Request.IP,
		UserAgent: "",
	})
	if err != nil {
		return LoginOutput{}, user.ErrGenerateToken
	}

	return LoginOutput{Token: authToken}, nil
}

// Logout User logout
func (u *UseCase) Logout(ctx context.Context, input shared.UseCaseInput[struct{}]) error {
	return u.tokenService.Revoke(ctx, input.Base.Auth.Token)
}

// CurrentUserInfo Get current user info
func (u *UseCase) CurrentUserInfo(ctx context.Context, input shared.UseCaseInput[struct{}]) (CurrentUserOutput, error) {
	id, err := user.ToID(input.Base.Auth.UserID)
	if err != nil {
		return CurrentUserOutput{}, user.ErrInvalidUser
	}

	domainUser, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return CurrentUserOutput{}, user.ErrUserNotFound
	}

	return CurrentUserOutput{
		Name:     domainUser.Name,
		Email:    string(domainUser.Email),
		CreateAt: timeutil.Format(domainUser.CreatedAt, "2006/01/02 15:04:05"),
	}, nil
}

package user

import "context"

type Repository interface {
	FindByID(ctx context.Context, id ID) (*User, error)
	FindByEmail(ctx context.Context, email Email) (*User, error)
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
}

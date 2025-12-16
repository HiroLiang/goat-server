package user

import "errors"

var (
	ErrInvalidID       = errors.New("invalid id format")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidUser     = errors.New("invalid user")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrGenerateToken   = errors.New("generate token error")
)

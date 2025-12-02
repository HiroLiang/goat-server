package user

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidEmail = errors.New("invalid email format")
	ErrInvalidID    = errors.New("invalid id format")
)

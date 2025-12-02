package auth

import "errors"

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrGenerateToken   = errors.New("generate token error")
	ErrRefreshToken    = errors.New("refresh token error")
)

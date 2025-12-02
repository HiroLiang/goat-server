package auth

import "time"

type Session struct {
	ID        string
	UserID    string
	IP        string
	UserAgent string
	CreatedAt time.Time
}

type CreateSessionParams struct {
	UserID    string
	IP        string
	UserAgent string
}

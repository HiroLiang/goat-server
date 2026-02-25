package chatmessage

import "errors"

var (
	ErrNotFound  = errors.New("chat message not found")
	ErrDeleted   = errors.New("chat message has been deleted")
	ErrForbidden = errors.New("operation not permitted for this chat message")
)

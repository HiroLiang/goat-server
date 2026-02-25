package chatgroup

import "errors"

var (
	ErrNotFound  = errors.New("chat group not found")
	ErrDeleted   = errors.New("chat group has been deleted")
	ErrFull      = errors.New("chat group has reached its member limit")
	ErrForbidden = errors.New("operation not permitted for this chat group type")
)

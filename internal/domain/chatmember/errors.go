package chatmember

import "errors"

var (
	ErrNotFound          = errors.New("chat member not found")
	ErrAlreadyMember     = errors.New("participant is already a member of this group")
	ErrForbidden         = errors.New("insufficient role to perform this action")
	ErrCannotRemoveOwner = errors.New("cannot remove the group owner")
)

package userrole

import "errors"

var (
	ErrAssignFailed            = errors.New("assign failed")
	ErrRevokeFailed            = errors.New("revoke failed")
	ErrUserRoleAlreadyAssigned = errors.New("user role already assigned")
)

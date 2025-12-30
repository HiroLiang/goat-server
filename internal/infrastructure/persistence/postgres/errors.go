package postgres

import "errors"

var (
	ErrNotFound = errors.New("record not found")
)

package postgres

import "errors"

var (
	ErrNotFound = errors.New("record not found")
	ErrExec     = errors.New("db exec error")
)

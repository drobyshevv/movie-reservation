package repository

import "errors"

var (
	ErrNotFound error = errors.New("not found")
	ErrConflict error = errors.New("conflict")
)

package errs

import "errors"

var (
	ErrSessionAlreadyExists      = errors.New("session already exists")
	ErrSessionNotFound           = errors.New("session not found")
	ErrSessionNotCancelled       = errors.New("session not cancelled")
	ErrConflictingSessionUpdates = errors.New("cannot update status and period simultaneously")
)

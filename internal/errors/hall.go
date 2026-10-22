package errs

import "errors"

var (
	ErrHallAlreadyExists     = errors.New("hall already exists")
	ErrHallNotFound          = errors.New("hall not found")
	ErrHallHasActiveSessions = errors.New("hall has active sessions")
)

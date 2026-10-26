package errs

import "errors"

var (
	ErrHallAlreadyExists           = errors.New("hall already exists")
	ErrHallNotFound                = errors.New("hall not found")
	ErrHallHasNonCancelledSessions = errors.New("hall has non-cancelled sessions")
)

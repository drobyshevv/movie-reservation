package errs

import "errors"

var (
	ErrHallAlreadyExists = errors.New("hall already exists")
	ErrHallNotFound      = errors.New("hall not found")
	ErrHallHasSessions   = errors.New("hall has sessions")
)

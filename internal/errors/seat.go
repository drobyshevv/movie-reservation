package errs

import "errors"

var (
	ErrSeatAlreadyExists       = errors.New("seats already exists")
	ErrSeatNotFound            = errors.New("seats not found")
	ErrDuplicateSeatsInRequest = errors.New("duplicate seats found in request payload")
)

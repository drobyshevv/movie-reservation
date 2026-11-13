package errs

import "errors"

var (
	ErrBookingNotFound      = errors.New("booking not found")
	ErrBookingAlreadyExists = errors.New("booking already exists")
)

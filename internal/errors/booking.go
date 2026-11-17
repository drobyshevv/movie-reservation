package errs

import "errors"

var (
	ErrBookingNotFound          = errors.New("booking not found")
	ErrBookingAlreadyExists     = errors.New("booking already exists")
	ErrSeatAlreadyBooked        = errors.New("seat already booked")
	ErrInvalidBookingConditions = errors.New("invalid booking conditions: session is not scheduled or seat belongs to a different hall")
)

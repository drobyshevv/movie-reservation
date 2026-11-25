package errs

import "errors"

var (
	ErrGenreAlreadyExists = errors.New("genre already exists")
	ErrGenreNotFound      = errors.New("genre not found")
)

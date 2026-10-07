package repository

import "errors"

var (
	ErrNotFound error = errors.New("not found")
	ErrConflict error = errors.New("conflict")

	ErrGenreNotFound error = errors.New("genre  not found")
	ErrMovieNotFound error = errors.New("movie not found")
)

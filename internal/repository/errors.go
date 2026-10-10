package repository

import "errors"

var (
	ErrNotFound error = errors.New("not found")
	ErrConflict error = errors.New("conflict")

	ErrGenreNotFound error = errors.New("genre  not found")
	ErrMovieNotFound error = errors.New("movie not found")
	ErrSeatNotFound  error = errors.New("seats not found")
	ErrHallNotFound  error = errors.New("hall not found")

	ErrSessionExists error = errors.New("session exists")

	ErrUnexpectedSeatsCount = errors.New("unexpected number of created seats")
)

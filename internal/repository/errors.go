package repository

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")

	ErrGenreNotFound = errors.New("genre  not found")
	ErrMovieNotFound = errors.New("movie not found")
	ErrSeatNotFound  = errors.New("seats not found")
	ErrHallNotFound  = errors.New("hall not found")

	ErrSessionExists   = errors.New("session exists")
	ErrSessionNotFound = errors.New("session not found")

	ErrUnexpectedSeatsCount = errors.New("unexpected number of created seats")
)

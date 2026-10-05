package service

import "errors"

var (
	ErrMovieAlreadyExists = errors.New("movie already exists")
	ErrMovieNotFound      = errors.New("movie not found")

	ErrHallAlreadyExists = errors.New("hall already exists")
	ErrHallNotFound      = errors.New("hall not found")
)

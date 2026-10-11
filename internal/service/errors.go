package service

import "errors"

var (
	ErrSeatAlreadyExists = errors.New("seats already exists")
	ErrSeatNotFound      = errors.New("seats not found")

	ErrMovieAlreadyExists = errors.New("movie already exists")
	ErrMovieNotFound      = errors.New("movie not found")

	ErrHallAlreadyExists = errors.New("hall already exists")
	ErrHallNotFound      = errors.New("hall not found")

	ErrGenreAlreadyExists = errors.New("genre already exists")
	ErrGenreNotFound      = errors.New("genre not found")

	ErrMovieGenreAlreadyExists = errors.New("movie_genre already exists")
	ErrMovieGenreNotFound      = errors.New("movie_genre already not found")

	ErrHallHasSessions error = errors.New("hall has sessions")
)

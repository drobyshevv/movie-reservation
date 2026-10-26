package errs

import "errors"

var (
	ErrMovieAlreadyExists           = errors.New("movie already exists")
	ErrMovieNotFound                = errors.New("movie not found")
	ErrMovieGenreAlreadyExists      = errors.New("movie_genre already exists")
	ErrMovieGenreNotFound           = errors.New("movie_genre already not found")
	ErrMovieHasNonCancelledSessions = errors.New("movie has non-cancelled sessions")
)

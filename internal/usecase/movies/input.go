package movies

import "time"

type CreateMovieInput struct {
	Title       string
	Description *string
	Duration    time.Duration
}

type UpdateMovieInput struct {
	Title       *string
	Description *string
	Duration    *time.Duration
}

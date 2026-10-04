package models

import (
	"time"
)

type Movie struct {
	ID          int64
	Title       string
	Description string
	Duration    time.Duration
	Image       []byte
	CreatedAt   time.Time
}

type CreateMovieParams struct {
	Title       string
	Description *string
	Duration    time.Duration
	Image       []byte
}

type UpdateMovieParams struct {
	Title       *string
	Description *string
	Duration    *time.Duration
	Image       []byte
}

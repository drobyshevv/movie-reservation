package dto

import "time"

type UpdateMovieRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,min=5,max=500"`
	Duration    *int    `json:"duration" validate:"omitempty,gt=0,lte=1000"`
}

type CreateMovieRequest struct {
	Title       string  `json:"title" validate:"required,max=100"`
	Description *string `json:"description" validate:"omitempty,min=5,max=500"`
	Duration    int     `json:"duration" validate:"required,numeric,gt=0,lte=1000"`
}

type MovieResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
}

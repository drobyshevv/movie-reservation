package dto

import "time"

type UpdateMovieRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Duration    *int    `json:"duration"`
	Image       []byte  `json:"image"`
}

type CreateMovieRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Duration    int     `json:"duration"`
	Image       []byte  `json:"image"`
}

type MovieResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	Image       []byte    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
}

package dto

import (
	"time"
)

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
	Description *string   `json:"description"`
	Duration    int       `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
}

/*
func (m *MovieResponse) FromModel(movie *models.Movie) {
	m.ID = movie.ID
	m.Title = movie.Title
	m.Description = movie.Description
	m.Duration = int(movie.Duration / time.Minute)
	m.CreatedAt = movie.CreatedAt
}

func (m *CreateMovieRequest) ToModel() *models.CreateMovieParams {
	return &models.CreateMovieParams{
		Title:       m.Title,
		Description: m.Description,
		Duration:    time.Duration(m.Duration) * time.Minute,
	}
}

func (m *UpdateMovieRequest) ToModel() *models.UpdateMovieParams {
	return &models.UpdateMovieParams{
		Title:       m.Title,
		Description: m.Description,
		Duration: func() *time.Duration {
			if m.Duration != nil {
				d := time.Duration(*m.Duration) * time.Minute
				return &d
			}
			return nil
		}(),
	}
}
*/

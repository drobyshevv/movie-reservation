package dto

import (
	"time"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	movieuc "github.com/drobyshevv/movie-reservation/internal/usecase/movie"
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

type MovieWithGenresResponse struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Description *string         `json:"description"`
	Duration    int             `json:"duration"`
	Genres      []GenreResponse `json:"genres"`
	CreatedAt   time.Time       `json:"created_at"`
}

func (m *MovieResponse) FromModel(movie *domain.Movie) {
	m.ID = movie.ID()
	m.Title = movie.Title()
	m.Description = movie.Description()
	m.Duration = int(movie.Duration() / time.Minute)
	m.CreatedAt = movie.CreatedAt()
}

func (m *CreateMovieRequest) ToModel() movieuc.CreateMovieInput {
	return movieuc.CreateMovieInput{
		Title:       m.Title,
		Description: m.Description,
		Duration:    time.Duration(m.Duration) * time.Minute,
	}
}

func (m *UpdateMovieRequest) ToModel() movieuc.UpdateMovieInput {
	return movieuc.UpdateMovieInput{
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

func (m *MovieWithGenresResponse) FromModel(withGenres *movieuc.MovieWithGenres) {
	m.ID = withGenres.Movie.ID()
	m.Title = withGenres.Movie.Title()
	m.Description = withGenres.Movie.Description()
	m.Duration = int(withGenres.Movie.Duration() / time.Minute)
	m.CreatedAt = withGenres.Movie.CreatedAt()
	for _, g := range withGenres.Genres {
		m.Genres = append(m.Genres, GenreResponse{
			ID:        g.ID(),
			TypeGenre: g.TypeGenre(),
		})
	}
}

func FromModelList(movies []domain.Movie) []MovieResponse {
	resp := []MovieResponse{}
	for _, m := range movies {
		resp = append(resp, MovieResponse{
			ID:          m.ID(),
			Title:       m.Title(),
			Description: m.Description(),
			Duration:    int(m.Duration() / time.Minute),
			CreatedAt:   m.CreatedAt(),
		})
	}
	return resp
}

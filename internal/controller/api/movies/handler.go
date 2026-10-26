package movies

import (
	"context"
	"log/slog"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	movieuc "github.com/drobyshevv/movie-reservation/internal/usecase/movies"
	"github.com/go-playground/validator/v10"
)

type MovieHandler struct {
	log       *slog.Logger
	usecase   MovieUseCase
	validator *validator.Validate
}

func NewMovieHandler(log *slog.Logger, usecase MovieUseCase, validator *validator.Validate) *MovieHandler {
	return &MovieHandler{
		log:       log,
		usecase:   usecase,
		validator: validator,
	}
}

type MovieUseCase interface {
	List(ctx context.Context) ([]domain.Movie, error)
	Get(ctx context.Context, id int64) (*domain.Movie, error)
	GetWithGenres(ctx context.Context, id int64) (*movieuc.MovieWithGenres, error)
	Create(ctx context.Context, input movieuc.CreateMovieInput) (*domain.Movie, error)
	Update(ctx context.Context, id int64, input movieuc.UpdateMovieInput) (*domain.Movie, error)
	Delete(ctx context.Context, id int64) error
	CreateMovieGenre(ctx context.Context, movieID int64, genreID int64) error
	DeleteMovieGenre(ctx context.Context, movieID int64, genreID int64) error
}

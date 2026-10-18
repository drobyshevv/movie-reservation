package genres

import (
	"context"
	"log/slog"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	genreuc "github.com/drobyshevv/movie-reservation/internal/usecase/genres"
	"github.com/go-playground/validator/v10"
)

type GenreHandler struct {
	log       *slog.Logger
	usecase   GenreUseCase
	validator *validator.Validate
}

func NewGenreHandler(log *slog.Logger, usecase GenreUseCase, validator *validator.Validate) *GenreHandler {
	return &GenreHandler{
		log:       log,
		usecase:   usecase,
		validator: validator,
	}
}

type GenreUseCase interface {
	List(ctx context.Context) ([]domain.Genre, error)
	Get(ctx context.Context, id int64) (*domain.Genre, error)
	Create(ctx context.Context, typeGenre string) (*domain.Genre, error)
	Update(ctx context.Context, id int64, typeGenre string) (*domain.Genre, error)
	Delete(ctx context.Context, id int64) error
	GetGenreMovies(ctx context.Context, id int64) ([]genreuc.GenreMovie, error)
}

package movies

import (
	"context"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

type MovieUseCase struct {
	repository MovieRepository
}

func NewMovieUseCase(repository MovieRepository) *MovieUseCase {
	return &MovieUseCase{
		repository: repository,
	}
}

type MovieRepository interface {
	List(ctx context.Context) ([]domain.Movie, error)
	Get(ctx context.Context, id int64) (*domain.Movie, error)
	GetWithGenres(ctx context.Context, id int64) (*MovieWithGenres, error)
	Create(ctx context.Context, movie domain.Movie) (*domain.Movie, error)
	Update(ctx context.Context, id int64, updateFn func(*domain.Movie) error) (*domain.Movie, error)
	Delete(ctx context.Context, id int64) error
	CreateMovieGenre(ctx context.Context, movieID int64, genreID int64) error
	DeleteMovieGenre(ctx context.Context, movieID int64, genreID int64) error
}

package genres

import (
	"context"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

type GenreUseCase struct {
	repository GenreRepository
}

func NewGenreUseCase(repository GenreRepository) *GenreUseCase {
	return &GenreUseCase{
		repository: repository,
	}
}

type GenreRepository interface {
	List(ctx context.Context) ([]domain.Genre, error)
	Get(ctx context.Context, id int64) (*domain.Genre, error)
	Create(ctx context.Context, genre domain.Genre) (*domain.Genre, error)
	Update(ctx context.Context, id int64, genre domain.Genre) (*domain.Genre, error)
	Delete(ctx context.Context, id int64) error
	GetGenreMovies(ctx context.Context, id int64) ([]GenreMovie, error)
}

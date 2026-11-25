package movies

import (
	"context"
	"io"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	"github.com/google/uuid"
)

type MovieUseCase struct {
	repository MovieRepository
	storage    MovieStorage
}

func NewMovieUseCase(repository MovieRepository, storage MovieStorage) *MovieUseCase {
	return &MovieUseCase{
		repository: repository,
		storage:    storage,
	}
}

type MovieRepository interface {
	List(ctx context.Context) ([]domain.Movie, error)
	Get(ctx context.Context, id int64) (*domain.Movie, error)
	Create(ctx context.Context, movie domain.Movie) (*domain.Movie, error)
	Update(ctx context.Context, id int64, updateFn func(*domain.Movie) error) (*domain.Movie, error)
	Delete(ctx context.Context, id int64) error

	GetWithGenres(ctx context.Context, id int64) (*MovieWithGenres, error)
	CreateMovieGenre(ctx context.Context, movieID int64, genreID int64) error
	DeleteMovieGenre(ctx context.Context, movieID int64, genreID int64) error

	GetPosterID(ctx context.Context, id int64) (uuid.UUID, error)
	PutPosterID(ctx context.Context, id int64, posterID uuid.UUID) error
	ClearPosterID(ctx context.Context, id int64, posterID uuid.UUID) error
}

type MovieStorage interface {
	PutPoster(ctx context.Context, file io.ReadSeeker, size int64, contentType string) (uuid.UUID, error)
	GetPoster(ctx context.Context, id uuid.UUID) (*File, error)
	DeletePoster(ctx context.Context, id uuid.UUID) error
}

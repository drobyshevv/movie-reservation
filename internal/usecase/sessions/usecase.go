package sessions

import (
	"context"
	"time"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

type UseCase struct {
	repository SessionRepository
	movies     MovieReader
}

func New(repository SessionRepository, movies MovieReader) *UseCase {
	return &UseCase{
		repository: repository,
		movies:     movies,
	}
}

type SessionRepository interface {
	List(ctx context.Context) ([]SessionWithDetails, error)
	Get(ctx context.Context, id int64) (*SessionWithDetails, error)
	Create(ctx context.Context, session domain.Session) (*domain.Session, error)
	Update(ctx context.Context, id int64, updateFn func(session *domain.Session) error) (*domain.Session, error)
	Delete(ctx context.Context, id int64, deleteFn func(session *domain.Session) error) error
}

type MovieReader interface {
	GetDuration(ctx context.Context, movieID int64) (time.Duration, error)
}

package halls

import (
	"context"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

type HallUseCase struct {
	repository HallRepository
}

func NewHallUseCase(repository HallRepository) *HallUseCase {
	return &HallUseCase{
		repository: repository,
	}
}

type HallRepository interface {
	List(ctx context.Context) ([]domain.Hall, error)
	Get(ctx context.Context, id int64) (*domain.Hall, error)
	GetWithSeats(ctx context.Context, id int64) (*HallWithSeats, error)
	GetSummary(ctx context.Context, id int64) (*HallSummary, error)
	CreateWithSeats(ctx context.Context, input CreateHallWithSeatsInput) (*HallWithSeats, error)
	Update(ctx context.Context, id int64, name string) (*domain.Hall, error)
	Delete(ctx context.Context, id int64) error
	ReplaceSeats(ctx context.Context, id int64, seats []domain.Seat) ([]domain.Seat, error)
	DeleteSeats(ctx context.Context, id int64) error
}

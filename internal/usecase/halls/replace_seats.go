package halls

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *HallUseCase) ReplaceSeats(ctx context.Context, id int64, input []HallSeatsInput) ([]domain.Seat, error) {
	const op = "usecase.halls.ReplaceSeats"

	domainSeats := make([]domain.Seat, 0, len(input))
	for _, in := range input {
		seat, err := domain.NewSeat(id, in.Row, in.Number)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid seat data: %w", op, err)
		}
		domainSeats = append(domainSeats, *seat)
	}

	seats, err := u.repository.ReplaceSeats(ctx, id, domainSeats)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return seats, nil
}

package halls

import (
	"context"
	"fmt"
)

func (u *HallUseCase) CreateWithSeats(ctx context.Context, input CreateHallWithSeatsInput) (*HallWithSeats, error) {
	const op = "usecase.halls.CreateWithSeats"

	hall, err := u.repository.CreateWithSeats(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return hall, nil
}

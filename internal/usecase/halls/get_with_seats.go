package halls

import (
	"context"
	"fmt"
)

func (u *HallUseCase) GetWithSeats(ctx context.Context, id int64) (*HallWithSeats, error) {
	const op = "usecase.halls.GetWithSeats"

	halls, err := u.repository.GetWithSeats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return halls, nil
}

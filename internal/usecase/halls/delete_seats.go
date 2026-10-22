package halls

import (
	"context"
	"fmt"
)

func (u *HallUseCase) DeleteSeats(ctx context.Context, id int64) error {
	const op = "usecase.halls.DeleteSeats"

	err := u.repository.DeleteSeats(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

package halls

import (
	"context"
	"fmt"
)

func (u *HallUseCase) DeleteHall(ctx context.Context, id int64) error {
	const op = "usecase.halls.Delete"

	err := u.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

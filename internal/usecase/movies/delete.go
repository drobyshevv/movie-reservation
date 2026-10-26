package movies

import (
	"context"
	"fmt"
)

func (u *MovieUseCase) Delete(ctx context.Context, id int64) error {
	const op = "usecase.movie.Delete"

	err := u.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

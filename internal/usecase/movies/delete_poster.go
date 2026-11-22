package movies

import (
	"context"
	"fmt"
)

func (u *MovieUseCase) DeletePoster(ctx context.Context, id int64) error {
	const op = "usecase.movies.DeletePoster"

	posterID, err := u.repository.GetPosterID(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.storage.DeletePoster(ctx, posterID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.repository.ClearPosterID(ctx, id, posterID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

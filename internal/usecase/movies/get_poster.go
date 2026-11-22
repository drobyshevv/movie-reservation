package movies

import (
	"context"
	"fmt"
)

func (u *MovieUseCase) GetPoster(ctx context.Context, id int64) (*File, error) {
	const op = "usecase.movies.GetPoster"

	posterID, err := u.repository.GetPosterID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	file, err := u.storage.GetPoster(ctx, posterID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return file, nil
}

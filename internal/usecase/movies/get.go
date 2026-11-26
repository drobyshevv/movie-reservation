package movies

import (
	"context"
	"fmt"
)

func (u *MovieUseCase) Get(ctx context.Context, id int64) (*MovieWithGenres, error) {
	const op = "usecase.movie.Get"

	movieWithGenres, err := u.repository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movieWithGenres, nil
}

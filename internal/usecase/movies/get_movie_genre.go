package movies

import (
	"context"
	"fmt"
)

func (u *MovieUseCase) GetWithGenres(ctx context.Context, id int64) (*MovieWithGenres, error) {
	const op = "usecase.movie.GetMovieGenres"

	withGenres, err := u.repository.GetWithGenres(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return withGenres, nil

}

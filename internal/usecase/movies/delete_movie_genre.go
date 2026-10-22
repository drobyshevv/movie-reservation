package movies

import (
	"context"
	"fmt"
)

func (u *MovieUseCase) DeleteMovieGenre(ctx context.Context, movieID int64, genreID int64) error {
	const op = "usecase.movie.DeleteMovieGenre"

	err := u.repository.DeleteMovieGenre(ctx, movieID, genreID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

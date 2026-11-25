package movies

import (
	"context"
	"fmt"
)

func (u *MovieUseCase) CreateMovieGenre(ctx context.Context, movieID int64, genreID int64) error {
	const op = "usecase.movie.CreateMovieGenre"

	err := u.repository.CreateMovieGenre(ctx, movieID, genreID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

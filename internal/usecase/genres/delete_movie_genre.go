package genres

import (
	"context"
	"fmt"
)

func (u *GenreUseCase) DeleteMovieGenre(ctx context.Context, movieID int64, genreID int64) error {
	const op = "usecase.genres.DeleteMovieGenre"

	err := u.repository.DeleteMovieGenre(ctx, movieID, genreID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

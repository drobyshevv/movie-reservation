package genres

import (
	"context"
	"fmt"
)

func (u *GenreUseCase) CreateMovieGenre(ctx context.Context, movieID int64, genreID int64) error {
	const op = "usecase.genres.CreateMovieGenre"

	err := u.repository.CreateMovieGenre(ctx, movieID, genreID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

package genres

import (
	"context"
	"fmt"
)

func (u *GenreUseCase) GetGenreMovies(ctx context.Context, id int64) ([]GenreMovie, error) {
	const op = "usecase.genres.GetGenreMovies"

	movies, err := u.repository.GetGenreMovies(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

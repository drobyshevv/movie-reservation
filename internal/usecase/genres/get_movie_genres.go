package genres

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *GenreUseCase) GetMovieGenres(ctx context.Context, movieID int64) ([]domain.Genre, error) {
	const op = "usecase.genres.GetMovieGenres"

	genres, err := u.repository.GetMovieGenres(ctx, movieID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return genres, nil
}

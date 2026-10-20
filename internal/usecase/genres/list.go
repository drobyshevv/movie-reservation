package genres

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *GenreUseCase) List(ctx context.Context) ([]domain.Genre, error) {
	const op = "usecase.genres.GetGenres"

	genres, err := u.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return genres, nil
}

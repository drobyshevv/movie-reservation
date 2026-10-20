package movie

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *MovieUseCase) List(ctx context.Context) ([]domain.Movie, error) {
	const op = "usecase.movie.List"

	movies, err := u.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

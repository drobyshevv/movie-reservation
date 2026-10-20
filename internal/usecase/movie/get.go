package movie

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *MovieUseCase) Get(ctx context.Context, id int64) (*domain.Movie, error) {
	const op = "usecase.movie.Get"

	movie, err := u.repository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movie, nil
}

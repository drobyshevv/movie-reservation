package genres

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *GenreUseCase) Get(ctx context.Context, id int64) (*domain.Genre, error) {
	const op = "usecase.genres.Get"

	genre, err := u.repository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return genre, nil
}

package halls

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *HallUseCase) Get(ctx context.Context, id int64) (*domain.Hall, error) {
	const op = "usecase.halls.Get"

	hall, err := u.repository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return hall, nil
}

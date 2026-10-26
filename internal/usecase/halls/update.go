package halls

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *HallUseCase) Update(ctx context.Context, id int64, name string) (*domain.Hall, error) {
	const op = "usecase.halls.Update"

	hall, err := u.repository.Update(ctx, id, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return hall, nil
}

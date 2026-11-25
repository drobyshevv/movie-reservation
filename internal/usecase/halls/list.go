package halls

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *HallUseCase) List(ctx context.Context) ([]domain.Hall, error) {
	const op = "usecase.halls.List"

	halls, err := u.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return halls, nil
}

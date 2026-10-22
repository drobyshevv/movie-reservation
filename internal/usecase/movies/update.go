package movies

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *MovieUseCase) Update(ctx context.Context, id int64, input UpdateMovieInput) (*domain.Movie, error) {
	const op = "usecase.movie.Update"

	movieOutput, err := u.repository.Update(ctx, id, func(m *domain.Movie) error {
		if input.Title != nil {
			if err := m.ChangeTitle(*input.Title); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
		if input.Description != nil {
			m.ChangeDescription(input.Description)
		}
		if input.Duration != nil {
			if err := m.ChangeDuration(*input.Duration); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movieOutput, nil
}

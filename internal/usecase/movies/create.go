package movies

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *MovieUseCase) Create(ctx context.Context, input CreateMovieInput) (*domain.Movie, error) {
	const op = "usecase.movie.Create"

	movieInput, err := domain.NewMovie(input.Title, input.Description, input.Duration)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	movieOutput, err := u.repository.Create(ctx, movieInput)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movieOutput, nil
}

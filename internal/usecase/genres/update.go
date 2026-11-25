package genres

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *GenreUseCase) Update(ctx context.Context, id int64, typeGenre string) (*domain.Genre, error) {
	const op = "usecase.genres.Update"

	genre, err := domain.NewGenre(typeGenre)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	genreOutput, err := u.repository.Update(ctx, id, genre)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return genreOutput, nil
}

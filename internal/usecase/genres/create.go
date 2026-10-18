package genres

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (s *GenreUseCase) Create(ctx context.Context, typeGenre string) (*domain.Genre, error) {
	const op = "usecase.genres.Create"

	genre, err := domain.NewGenre(typeGenre)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	genreOutput, err := s.repository.Create(ctx, genre)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return genreOutput, nil
}

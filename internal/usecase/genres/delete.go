package genres

import (
	"context"
	"fmt"
)

func (s *GenreUseCase) Delete(ctx context.Context, id int64) error {
	const op = "usecase.genres.DeleteGenre"

	err := s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

package sessions

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (u *UseCase) Delete(ctx context.Context, id int64) error {
	const op = "usecase.sessions.Delete"

	err := u.repository.Delete(ctx, id, func(session *domain.Session) error {
		if session.Status() != domain.StatusCancelled {
			return fmt.Errorf("%s: %w", op, errs.ErrSessionNotCancelled)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

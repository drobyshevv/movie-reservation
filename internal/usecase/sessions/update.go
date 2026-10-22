package sessions

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (u *SessionUseCase) Update(ctx context.Context, id int64, input UpdateSessionInput) (*domain.Session, error) {
	const op = "usecase.sessions.Update"

	session, err := u.repository.Update(ctx, id, func(session *domain.Session) error {
		if input.Status != nil && (input.StartAt != nil || input.EndAt != nil) {
			return fmt.Errorf("%s: %w", op, errs.ErrConflictingSessionUpdates)
		}

		if input.Status != nil {
			err := session.ChangeStatus(*input.Status)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

		startAt := session.Period().Start()
		endAt := session.Period().End()

		if input.StartAt != nil {
			startAt = *input.StartAt
		}
		if input.EndAt != nil {
			endAt = *input.EndAt
		}

		period, err := domain.NewTimeRange(startAt, endAt)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = session.ChangePeriod(period)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return session, nil
}

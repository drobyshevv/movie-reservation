package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *UseCase) Create(ctx context.Context, input CreateSessionInput) (*domain.Session, error) {
	const op = "usecase.sessions.Create"

	var endAt time.Time

	if input.EndAt != nil {
		endAt = *input.EndAt
	} else {
		duraton, err := u.movies.GetDuration(ctx, input.MovieID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		endAt = input.StartAt.Add(duraton)
	}

	timeRange, err := domain.NewTimeRange(input.StartAt, endAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	session, err := domain.NewSession(
		input.HallID,
		input.MovieID,
		timeRange,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sessionOutput, err := u.repository.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sessionOutput, nil
}

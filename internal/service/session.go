package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/drobyshevv/movie-reservation/internal/repository"
)

type SessionService struct {
	repository SessionRepository
	movies     MovieReader
}

func NewSessionService(repository SessionRepository, movies MovieReader) *SessionService {
	return &SessionService{
		repository: repository,
		movies:     movies,
	}
}

type SessionRepository interface {
	GetSessions(ctx context.Context) ([]models.SessionWithDetails, error)
	GetSession(ctx context.Context, id int64) (*models.SessionWithDetails, error)
	CreateSession(ctx context.Context, session domain.Session) (*domain.Session, error)
	UpdateSession(ctx context.Context, id int64, updateFn func(session *domain.Session) error) (*domain.Session, error)
	DeleteSession(ctx context.Context, id int64, deleteFn func(session *domain.Session) error) error
}

type MovieReader interface {
	GetDuration(ctx context.Context, movieID int64) (time.Duration, error)
}

func (s *SessionService) GetSessions(ctx context.Context) ([]models.SessionWithDetails, error) {
	const op = "service.GetSessions"

	sessions, err := s.repository.GetSessions(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sessions, nil
}

func (s *SessionService) GetSession(ctx context.Context, id int64) (*models.SessionWithDetails, error) {
	const op = "service.GetSession"

	session, err := s.repository.GetSession(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrSessionNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (s *SessionService) CreateSession(ctx context.Context, params models.CreateSessionParams) (*domain.Session, error) {
	const op = "service.CreateSession"

	var endAt time.Time

	if params.EndAt != nil {
		endAt = *params.EndAt
	} else {
		duraton, err := s.movies.GetDuration(ctx, params.MovieID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		endAt = params.StartAt.Add(duraton)
	}

	timeRange, err := domain.NewTimeRange(params.StartAt, endAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	session, err := domain.NewSession(
		params.HallID,
		params.MovieID,
		timeRange,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sessionOutput, err := s.repository.CreateSession(ctx, session)
	if err != nil {
		if errors.Is(err, repository.ErrSessionExists) {
			return nil, fmt.Errorf("%s: %w", op, ErrSessionAlreadyExists)
		}
		if errors.Is(err, repository.ErrHallNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrHallNotFound)
		}
		if errors.Is(err, repository.ErrMovieNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrMovieNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sessionOutput, nil

}

func (s *SessionService) UpdateSession(ctx context.Context, id int64, params models.UpdateSessionParams) (*domain.Session, error) {
	const op = "service.UpdateSession"

	session, err := s.repository.UpdateSession(ctx, id, func(session *domain.Session) error {
		if params.Status != nil && (params.StartAt != nil || params.EndAt != nil) {
			return ErrConflictingSessionUpdates
		}

		if params.Status != nil {
			err := session.ChangeStatus(*params.Status)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

		startAt := session.Period().Start()
		endAt := session.Period().End()

		if params.StartAt != nil {
			startAt = *params.StartAt
		}
		if params.EndAt != nil {
			endAt = *params.EndAt
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

func (s *SessionService) DeleteSession(ctx context.Context, id int64) error {
	const op = "service.DeleteSession"

	err := s.repository.DeleteSession(ctx, id, func(session *domain.Session) error {
		if session.Status() != domain.StatusCancelled {
			return fmt.Errorf("%s: %w", op, ErrSessionNotCancelled)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

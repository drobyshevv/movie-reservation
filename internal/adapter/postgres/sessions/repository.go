package sessions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/drobyshevv/movie-reservation/internal/adapter/postgres"
	"github.com/drobyshevv/movie-reservation/internal/domain"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
	sessionuc "github.com/drobyshevv/movie-reservation/internal/usecase/sessions"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		pool: pool,
	}
}

func (r *SessionRepository) List(ctx context.Context) ([]sessionuc.SessionWithDetails, error) {
	const op = "postgres.sessions.List"

	query := `
	SELECT 
		s.id, 
		s.status,
		s.start_at,
		s.end_at,
		s.created_at,
		h.id,
		h.name,
		m.id,
		m.title
	FROM sessions AS s
	JOIN halls as h
	ON s.hall_id = h.id
	JOIN movies AS m
	ON s.movie_id = m.id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var sessions []sessionuc.SessionWithDetails

	for rows.Next() {
		var s sessionuc.SessionWithDetails

		err := rows.Scan(
			&s.ID, &s.Status, &s.StartAt, &s.EndAt, &s.CreatedAt,
			&s.Hall.ID, &s.Hall.Name,
			&s.Movie.ID, &s.Movie.Title,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: iterate rows: %w", op, err)
	}

	return sessions, nil
}

func (r *SessionRepository) Get(ctx context.Context, id int64) (*sessionuc.SessionWithDetails, error) {
	const op = "postgres.sessions.Get"

	query := `
	SELECT 
		s.id, 
		s.status,
		s.start_at,
		s.end_at,
		s.created_at,
		h.id,
		h.name,
		m.id,
		m.title
	FROM sessions AS s
	JOIN halls as h
	ON s.hall_id = h.id
	JOIN movies AS m
	ON s.movie_id = m.id
	WHERE s.id = $1
	`

	var session sessionuc.SessionWithDetails

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.Status,
		&session.StartAt,
		&session.EndAt,
		&session.CreatedAt,
		&session.Hall.ID,
		&session.Hall.Name,
		&session.Movie.ID,
		&session.Movie.Title,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &session, nil
}

func (r *SessionRepository) Create(ctx context.Context, session domain.Session) (*domain.Session, error) {
	const op = "postgres.sessions.Create"

	query := `
	INSERT INTO sessions (hall_id, movie_id, status, start_at, end_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, hall_id, movie_id, status, start_at, end_at, created_at
	`

	var (
		id        int64
		hallID    int64
		movieID   int64
		status    domain.StatusType
		startAt   time.Time
		endAt     time.Time
		createdAt time.Time
	)

	err := r.pool.QueryRow(
		ctx,
		query,
		session.HallID(),
		session.MovieID(),
		session.Status(),
		session.Period().Start(),
		session.Period().End(),
	).Scan(
		&id,
		&hallID,
		&movieID,
		&status,
		&startAt,
		&endAt,
		&createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	period, err := domain.NewTimeRange(startAt, endAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sessionOutput, err := domain.RestoreSession(
		id,
		hallID,
		movieID,
		status,
		period,
		createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &sessionOutput, nil
}

func (r *SessionRepository) Update(ctx context.Context, id int64, updateFn func(session *domain.Session) error) (*domain.Session, error) {
	return postgres.RunInTxWithResult(ctx, r.pool, func(tx pgx.Tx) (*domain.Session, error) {
		const op = "postgres.sessions.Update"

		var (
			sessionID int64
			hallID    int64
			movieID   int64
			status    domain.StatusType
			startAt   time.Time
			endAt     time.Time
			createdAt time.Time
		)

		err := tx.QueryRow(
			ctx, `
		SELECT 
			id, 
			hall_id,
			movie_id,
			status,
			start_at,
			end_at,
			created_at
		FROM sessions
		WHERE id = $1
		FOR UPDATE
		`, id).Scan(
			&sessionID,
			&hallID,
			&movieID,
			&status,
			&startAt,
			&endAt,
			&createdAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, errs.ErrSessionNotFound)
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		period, err := domain.NewTimeRange(startAt, endAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		session, err := domain.RestoreSession(
			sessionID,
			hallID,
			movieID,
			status,
			period,
			createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		err = updateFn(&session)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		tag, err := tx.Exec(
			ctx, `
			UPDATE sessions
			SET 
				status = $1,
				start_at = $2,
				end_at = $3
			WHERE id = $4
			`,
			session.Status(),
			session.Period().Start(),
			session.Period().End(),
			session.ID(),
		)
		if err != nil {
			return nil, fmt.Errorf("%s: update session: %w", op, err)
		}

		if tag.RowsAffected() == 0 {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrSessionNotFound)
		}

		return &session, nil
	})
}

func (r *SessionRepository) Delete(ctx context.Context, id int64, deleteFn func(session *domain.Session) error) error {
	return postgres.RunInTx(ctx, r.pool, func(tx pgx.Tx) error {
		const op = "postgres.sessions.Delete"

		var (
			sessionID int64
			hallID    int64
			movieID   int64
			status    domain.StatusType
			startAt   time.Time
			endAt     time.Time
			createdAt time.Time
		)

		err := tx.QueryRow(
			ctx, `
		SELECT 
			id, 
			hall_id,
			movie_id,
			status,
			start_at,
			end_at,
			created_at
		FROM sessions
		WHERE id = $1
		FOR UPDATE
		`, id).Scan(
			&sessionID,
			&hallID,
			&movieID,
			&status,
			&startAt,
			&endAt,
			&createdAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%s: %w", op, errs.ErrSessionNotFound)
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		period, err := domain.NewTimeRange(startAt, endAt)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		session, err := domain.RestoreSession(
			sessionID,
			hallID,
			movieID,
			status,
			period,
			createdAt,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = deleteFn(&session)
		if err != nil {
			return fmt.Errorf("%s: delete session: %w", op, err)
		}

		_, err = tx.Exec(
			ctx, `
		DELETE FROM sessions
		WHERE id = $1
		`, id)
		if err != nil {
			return fmt.Errorf("%s: delete session: %w", op, err)
		}

		return nil
	})
}

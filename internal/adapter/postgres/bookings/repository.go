package bookings

import (
	"context"
	"errors"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/adapter/postgres"
	"github.com/drobyshevv/movie-reservation/internal/domain"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
	bookinguc "github.com/drobyshevv/movie-reservation/internal/usecase/bookings"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	pool *pgxpool.Pool
}

func NewBookingRepository(pool *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{
		pool: pool,
	}
}

func (r *BookingRepository) Get(ctx context.Context, id int64) (*domain.Booking, error) {
	const op = "postgres.bookings.Get"

	query := `
	SELECT
		id,
		session_id,
		user_id,
		seat_id,
		status,
		created_at
	FROM bookings
	WHERE id = $1;
	`

	bookingRaw := bookingRow{}

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&bookingRaw.id,
		&bookingRaw.sessionID,
		&bookingRaw.userID,
		&bookingRaw.seatID,
		&bookingRaw.status,
		&bookingRaw.createdAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrBookingNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	booking, err := domain.RestoreBooking(
		bookingRaw.id,
		bookingRaw.sessionID,
		bookingRaw.userID,
		bookingRaw.seatID,
		bookingRaw.status,
		bookingRaw.createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s restore model: %w", op, err)
	}

	return booking, nil
}

func (r *BookingRepository) Create(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	const op = "postgres.bookings.Create"

	query := `
	WITH
	verification AS (
		SELECT
			$1 AS session_id,
			$2 AS user_id,
			$3 AS seat_id,
			$4 AS status,
			sn.status AS session_status,
			sn.hall_id AS session_hall_id,
			st.hall_id AS seat_hall_id
		FROM sessions sn
		JOIN seats st ON st.id = $3
		WHERE sn.id = $1
	)
	INSERT INTO bookings (session_id, user_id, seat_id, status)
	SELECT session_id, user_id, seat_id, status FROM verification v
	WHERE v.session_hall_id = v.seat_hall_id
	AND v.session_status = 'scheduled'
	RETURNING id, session_id, user_id, seat_id, status, created_at;
	`

	bookingRaw := bookingRow{}

	err := r.pool.QueryRow(
		ctx,
		query,
		booking.SessionID(),
		booking.UserID(),
		booking.SeatID(),
		booking.Status(),
	).Scan(
		&bookingRaw.id,
		&bookingRaw.sessionID,
		&bookingRaw.userID,
		&bookingRaw.seatID,
		&bookingRaw.status,
		&bookingRaw.createdAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrInvalidBookingConditions)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrSeatAlreadyBooked)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	restored, err := domain.RestoreBooking(
		bookingRaw.id,
		bookingRaw.sessionID,
		bookingRaw.userID,
		bookingRaw.seatID,
		bookingRaw.status,
		bookingRaw.createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s restore model: %w", op, err)
	}

	return restored, nil
}

func (r *BookingRepository) ListBySession(ctx context.Context, sessionID int64) ([]bookinguc.SeatStatusOutput, error) {
	const op = "postgres.bookings.ListBySession"

	query := `
	SELECT seat_id, status
	FROM bookings
	WHERE session_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	seatStatus := []bookinguc.SeatStatusOutput{}

	for rows.Next() {
		current := bookinguc.SeatStatusOutput{}
		err := rows.Scan(&current.SeatID, &current.Status)
		if err != nil {
			return nil, fmt.Errorf("%s scan rows: %w", op, err)
		}

		seatStatus = append(seatStatus, current)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return seatStatus, nil
}

// TODO filter by status
func (r *BookingRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]bookinguc.UserBookingOutput, error) {
	const op = "postgres.bookings.ListByUser"

	query := `
	SELECT 
		b.id, 
		m.title, 
		h.name, 
		sn.start_at, 
		st.row, 
		st.number, 
		b.status, 
		sn.status 
	FROM bookings b
	JOIN seats st ON
		b.seat_id = st.id
	JOIN sessions sn ON
		b.session_id = sn.id
	JOIN halls h ON
		h.id = sn.hall_id
	JOIN movies m ON
		m.id = sn.movie_id
	WHERE b.user_id = $1
	ORDER BY sn.start_at ASC;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	userBookings := []bookinguc.UserBookingOutput{}

	for rows.Next() {
		current := bookinguc.UserBookingOutput{}
		var rowNum string
		var seatNum int16
		err := rows.Scan(
			&current.ID,
			&current.MovieTitle,
			&current.HallName,
			&current.StartAt,
			&rowNum,
			&seatNum,
			&current.BookingStatus,
			&current.SessionStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("%s scan rows: %w", op, err)
		}
		current.SeatNumber = fmt.Sprintf("%s%d", rowNum, seatNum)

		userBookings = append(userBookings, current)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userBookings, nil
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, id int64, status domain.BookingStatus, updateFn func(booking *domain.Booking, status domain.BookingStatus) error) error {
	return postgres.RunInTx(ctx, r.pool, func(tx pgx.Tx) error {
		const op = "postgres.bookings.UpdateStatus"

		bookingRaw := bookingRow{}
		err := tx.QueryRow(ctx, `
		SELECT
			id,
			session_id,
			user_id,
			seat_id,
			status,
			created_at
		FROM bookings
		WHERE id = $1
		FOR UPDATE; 	
		`).Scan(
			&bookingRaw.id,
			&bookingRaw.sessionID,
			&bookingRaw.userID,
			&bookingRaw.seatID,
			&bookingRaw.status,
			&bookingRaw.createdAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%s: %w", op, errs.ErrBookingNotFound)
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		restored, err := domain.RestoreBooking(
			bookingRaw.id,
			bookingRaw.sessionID,
			bookingRaw.userID,
			bookingRaw.seatID,
			bookingRaw.status,
			bookingRaw.createdAt,
		)
		if err != nil {
			return fmt.Errorf("%s restore model: %w", op, err)
		}

		err = updateFn(restored, status)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		tag, err := tx.Exec(ctx, `
			UPDATE bookings
			SET status = $1
			WHERE id = $2;	
		`, status, id)
		if err != nil {
			return fmt.Errorf("%s update status: %w", op, err)
		}

		if tag.RowsAffected() == 0 {
			return fmt.Errorf("%s: %w", op, errs.ErrBookingNotFound)
		}

		return nil
	})
}

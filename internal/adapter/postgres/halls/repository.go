package halls

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/drobyshevv/movie-reservation/internal/adapter/postgres"
	"github.com/drobyshevv/movie-reservation/internal/domain"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
	halluc "github.com/drobyshevv/movie-reservation/internal/usecase/halls"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HallRepository struct {
	pool *pgxpool.Pool
}

func NewHallRepository(pool *pgxpool.Pool) *HallRepository {
	return &HallRepository{
		pool: pool,
	}
}

func (r *HallRepository) List(ctx context.Context) ([]domain.Hall, error) {
	const op = "postgres.halls.List"

	query := `
	SELECT 
		id, 
		name 
	FROM halls`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	halls := []domain.Hall{}
	for rows.Next() {
		var hallID int64
		var name string

		if err := rows.Scan(
			&hallID,
			&name,
		); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		hall, err := domain.RestoreHall(
			hallID,
			name,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: restore: %w", op, err)
		}

		halls = append(halls, hall)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: iterate rows: %w", op, err)
	}

	return halls, nil
}

func (r *HallRepository) Get(ctx context.Context, id int64) (*domain.Hall, error) {
	const op = "postgres.halls.Get"

	query := `
	SELECT id, name FROM halls
	WHERE id = $1
	`

	var hallID int64
	var name string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&hallID,
		&name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrHallNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	hall, err := domain.RestoreHall(hallID, name)
	if err != nil {
		return nil, fmt.Errorf("%s: restore: %w", op, err)
	}

	return &hall, nil
}

func (r *HallRepository) GetWithSeats(ctx context.Context, id int64) (*halluc.HallWithSeats, error) {
	const op = "postgres.halls.GetWithSeats"

	var hallID int64
	var hallName string
	err := r.pool.QueryRow(ctx, `SELECT id, name FROM halls WHERE id = $1`, id).Scan(&hallID, &hallName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrHallNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	hall, err := domain.RestoreHall(hallID, hallName)
	if err != nil {
		return nil, fmt.Errorf("%s: restore hall: %w", op, err)
	}

	query := `
	SELECT id, hall_id, row, number 
	FROM seats
	WHERE hall_id = $1
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	seats := []domain.Seat{}
	for rows.Next() {
		var seatID int64
		var seatHallID int64
		var row string
		var number int16

		if err := rows.Scan(&seatID, &seatHallID, &row, &number); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		seat, err := domain.RestoreSeat(seatID, seatHallID, row, number)
		if err != nil {
			return nil, fmt.Errorf("%s: restore seat: %w", op, err)
		}

		seats = append(seats, *seat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: iterate rows: %w", op, err)
	}

	return &halluc.HallWithSeats{
		Hall:  hall,
		Seats: seats,
	}, nil
}

func (r *HallRepository) GetSummary(ctx context.Context, id int64) (*halluc.HallSummary, error) {
	const op = "postgres.halls.GetSummary"

	query := `
	SELECT 
		h.id, 
		h.name, 
		COUNT(s.id) as total_seats
	FROM halls h
	LEFT JOIN seats s ON h.id = s.hall_id
	WHERE h.id = $1
	GROUP BY h.id, h.name
	`

	summary := &halluc.HallSummary{}

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&summary.ID,
		&summary.Name,
		&summary.TotalSeats,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrHallNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return summary, nil
}

func (r *HallRepository) CreateWithSeats(ctx context.Context, input halluc.CreateHallWithSeatsInput) (*halluc.HallWithSeats, error) {
	return postgres.RunInTxWithResult(ctx, r.pool, func(tx pgx.Tx) (*halluc.HallWithSeats, error) {
		const op = "postgres.halls.CreateWithSeats"

		var hallID int64
		err := tx.QueryRow(ctx, `
			INSERT INTO halls (name) 
			VALUES ($1) 
			RETURNING id
		`, input.Name).Scan(&hallID)
		if err != nil {
			return nil, fmt.Errorf("%s: insert hall: %w", op, err)
		}

		domainSeats := make([]domain.Seat, 0, len(input.Seats))
		for _, s := range input.Seats {
			seat, err := domain.NewSeat(hallID, s.Row, s.Number)
			if err != nil {
				return nil, fmt.Errorf("%s: create domain seat: %w", op, err)
			}
			domainSeats = append(domainSeats, *seat)
		}

		query, args := makeQueryHallSeats(hallID, domainSeats)

		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("%s: insert seats: %w", op, err)
		}
		defer rows.Close()

		seats := []domain.Seat{}
		for rows.Next() {
			var seatID int64
			var seatHallID int64
			var row string
			var number int16

			if err := rows.Scan(&seatID, &seatHallID, &row, &number); err != nil {
				return nil, fmt.Errorf("%s: scan seat: %w", op, err)
			}

			seat, err := domain.RestoreSeat(seatID, seatHallID, row, number)
			if err != nil {
				return nil, fmt.Errorf("%s: restore seat: %w", op, err)
			}

			seats = append(seats, *seat)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("%s: iterate rows: %w", op, err)
		}

		hall, err := domain.RestoreHall(hallID, input.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: restore hall: %w", op, err)
		}

		return &halluc.HallWithSeats{
			Hall:  hall,
			Seats: seats,
		}, nil
	})
}

func (r *HallRepository) Update(ctx context.Context, id int64, name string) (*domain.Hall, error) {
	const op = "postgres.halls.Update"

	query := `
		UPDATE halls 
		SET name = $1 
		WHERE id = $2 
		RETURNING id, name
	`

	var hallID int64
	var updatedName string

	err := r.pool.QueryRow(ctx, query, name, id).Scan(&hallID, &updatedName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrHallNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	hall, err := domain.RestoreHall(hallID, updatedName)
	if err != nil {
		return nil, fmt.Errorf("%s: restore: %w", op, err)
	}

	return &hall, nil
}

func (r *HallRepository) Delete(ctx context.Context, id int64) error {
	return postgres.RunInTx(ctx, r.pool, func(tx pgx.Tx) error {
		const op = "postgres.halls.Delete"

		var hallID int64
		err := tx.QueryRow(ctx, `SELECT id FROM halls WHERE id = $1 FOR UPDATE`, id).Scan(&hallID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%s: %w", op, errs.ErrHallNotFound)
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		var hasActiveSessions bool
		err = tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM sessions 
				WHERE hall_id = $1 AND status != 'cancelled'
			)
		`, id).Scan(&hasActiveSessions)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if hasActiveSessions {
			return fmt.Errorf("%s: %w", op, errs.ErrHallHasNonCancelledSessions)
		}

		var deletedID int64
		err = tx.QueryRow(ctx, `DELETE FROM halls WHERE id = $1 RETURNING id`, id).Scan(&deletedID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})
}

func (r *HallRepository) ReplaceSeats(ctx context.Context, id int64, seats []domain.Seat) ([]domain.Seat, error) {
	return postgres.RunInTxWithResult(ctx, r.pool, func(tx pgx.Tx) ([]domain.Seat, error) {
		const op = "postgres.halls.ReplaceSeats"

		var hallID int64
		err := tx.QueryRow(ctx, `SELECT id FROM halls WHERE id = $1 FOR UPDATE`, id).Scan(&hallID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, errs.ErrHallNotFound)
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		var hasActiveSessions bool
		err = tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM sessions 
				WHERE hall_id = $1 AND status != 'cancelled'
			)
		`, id).Scan(&hasActiveSessions)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if hasActiveSessions {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrHallHasNonCancelledSessions)
		}

		_, err = tx.Exec(ctx, `DELETE FROM seats WHERE hall_id = $1`, id)
		if err != nil {
			return nil, fmt.Errorf("%s: delete old seats: %w", op, err)
		}

		query, args := makeQueryHallSeats(id, seats)
		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("%s: insert new seats: %w", op, err)
		}
		defer rows.Close()

		createdSeats := make([]domain.Seat, 0, len(seats))
		for rows.Next() {
			var seatID int64
			var seatHallID int64
			var row string
			var number int16

			if err := rows.Scan(&seatID, &seatHallID, &row, &number); err != nil {
				return nil, fmt.Errorf("%s: scan seat: %w", op, err)
			}

			seat, err := domain.RestoreSeat(seatID, seatHallID, row, number)
			if err != nil {
				return nil, fmt.Errorf("%s: restore seat: %w", op, err)
			}

			createdSeats = append(createdSeats, *seat)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("%s: iterate rows: %w", op, err)
		}

		return createdSeats, nil
	})
}

func (r *HallRepository) DeleteSeats(ctx context.Context, id int64) error {
	return postgres.RunInTx(ctx, r.pool, func(tx pgx.Tx) error {
		const op = "postgres.halls.DeleteSeats"

		var hallID int64
		err := tx.QueryRow(ctx, `SELECT id FROM halls WHERE id = $1 FOR UPDATE`, id).Scan(&hallID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%s: %w", op, errs.ErrHallNotFound)
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		var hasActiveSessions bool
		err = tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM sessions 
				WHERE hall_id = $1 AND status != 'cancelled'
			)
		`, id).Scan(&hasActiveSessions)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if hasActiveSessions {
			return fmt.Errorf("%s: %w", op, errs.ErrHallHasNonCancelledSessions)
		}

		_, err = tx.Exec(ctx, `DELETE FROM seats WHERE hall_id = $1`, id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})
}

func makeQueryHallSeats(id int64, seat []domain.Seat) (string, []any) {
	query := `
	INSERT INTO seats (hall_id, row, number) VALUES
	`
	var placeholders []string
	var args []any
	n := 1

	for _, p := range seat {
		args = append(args, id)
		args = append(args, p.Row())
		args = append(args, p.Number())
		phRow := fmt.Sprintf("($%d, $%d, $%d)", n, n+1, n+2)
		n += 3
		placeholders = append(placeholders, phRow)
	}
	query += strings.Join(placeholders, ", ")
	query += " RETURNING id, hall_id, row, number;"

	return query, args
}

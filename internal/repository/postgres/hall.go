package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/drobyshevv/movie-reservation/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *HallRepository) GetHalls(ctx context.Context) ([]models.Hall, error) {
	const op = "postgres.GetHalls"

	query := `SELECT id, name FROM halls`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	halls, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Hall])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return halls, nil
}

func (r *HallRepository) GetHall(ctx context.Context, id int64) (*models.Hall, error) {
	const op = "postgres.GetHall"
	query := `
	SELECT id, name FROM halls
	WHERE id = $1
	`

	hall := &models.Hall{}

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&hall.ID,
		&hall.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return hall, nil
}

func (r *HallRepository) CreateHall(ctx context.Context, params models.CreateHallParams) (*models.Hall, error) {
	const op = "postgres.CreateHall"

	query := `
	INSERT INTO halls (name)
	VALUES ($1)
	RETURNING id, name
	`

	hall := &models.Hall{}
	err := r.pool.QueryRow(ctx, query, params.Name).Scan(
		&hall.ID,
		&hall.Name,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrConflict)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return hall, nil
}

func (r *HallRepository) UpdateHall(ctx context.Context, id int64, params models.UpdateHallParams) (*models.Hall, error) {
	const op = "postgres.UpdateHall"

	query := `
	UPDATE halls
	SET 
	    name = COALESCE($1, name)
	WHERE id = $2
	RETURNING id, name
	`

	hall := &models.Hall{}
	err := r.pool.QueryRow(ctx, query, params.Name, id).Scan(
		&hall.ID,
		&hall.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return hall, nil
}

func (r *HallRepository) DeleteHall(ctx context.Context, id int64) error {
	const op = "postgres.DeleteHall"
	query := `
	DELETE FROM halls
	WHERE id = $1
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrNotFound)
	}

	return nil
}

func (r *HallRepository) GetHallSeats(ctx context.Context, id int64) ([]models.Seat, error) {
	const op = "repository.GetHallSeats"

	query := `
	SELECT id, hall_id, row, number FROM seats
	WHERE hall_id = $1
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	seats, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Seat])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return seats, nil
}

func (r *HallRepository) CreateHallSeats(ctx context.Context, id int64, params []models.CreateHallSeatsParams) ([]models.Seat, error) {
	return runInTxWithResult(ctx, r.pool, func(tx pgx.Tx) ([]models.Seat, error) {
		const op = "repository.CreateHallSeats"
		var hallID int64
		err := tx.QueryRow(ctx, `
		SELECT id 
		FROM halls
		WHERE id = $1
		FOR UPDATE
		`).Scan(&hallID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, repository.ErrHallNotFound)
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		var hasSession bool
		err = tx.QueryRow(ctx, `
		SELECT EXISTS (
		SELECT * FROM sessions WHERE halls_id = $1
		)
		`, id).Scan(&hasSession)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if hasSession != false {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrSessionExists)
		}

		var hasSeats bool
		err = tx.QueryRow(ctx, `
		SELECT EXISTS (
		SELECT * FROM seats WHERE halls_id = $1
		)
		`, id).Scan(&hasSeats)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if hasSeats != false {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrConflict)
		}

		query, args := makeQueryHallSeats(id, params)
		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()

		seats, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Seat])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if len(seats) != len(params) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrUnexpectedSeatsCount)
		}

		return seats, nil
	})
}

func (r *HallRepository) DeleteHallSeats(ctx context.Context, id int64) error {
	return runInTx(ctx, r.pool, func(tx pgx.Tx) error {
		const op = "repository.DeleteHallSeats"
		var hallID int64
		err := tx.QueryRow(ctx, `
		SELECT id 
		FROM halls 
		WHERE id = $1
		FOR UPDATE
		`, id).Scan(&hallID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%s: %w", op, repository.ErrHallNotFound)
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		var hasSession bool
		err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT * FROM sessions WHERE hall_id = $1
		)`, id).Scan(&hasSession)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if hasSession != false {
			return fmt.Errorf("%s: %w", op, repository.ErrSessionExists)
		}

		tag, err := tx.Exec(ctx, `
		DELETE FROM seats
		WHERE hall_id = $1
		`, id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("%s: %w", op, repository.ErrSeatNotFound)
		}

		return nil
	})
}

func makeQueryHallSeats(id int64, params []models.CreateHallSeatsParams) (string, []any) {
	query := `
	INSERT INTO seats (hall_id, row, number) VALUES
	`
	var placeholders []string
	var args []any
	n := 1

	for _, p := range params {
		args = append(args, id)
		args = append(args, p.Row)
		args = append(args, p.Number)
		phRow := fmt.Sprintf("($%d, $%d, $%d)", n, n+1, n+2)
		n += 3
		placeholders = append(placeholders, phRow)
	}
	query += strings.Join(placeholders, ", ")
	query += "RETURNING id, hall_id, row, number;"

	return query, args
}

package postgres

import (
	"context"
	"errors"
	"fmt"

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

func (r *HallRepository) GetHalls() ([]models.Hall, error) {
	const op = "postgres.GetHalls"

	query := `SELECT id, name FROM halls`

	rows, err := r.pool.Query(context.Background(), query)
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

func (r *HallRepository) GetHall(id int64) (*models.Hall, error) {
	const op = "postgres.GetHall"
	query := `
	SELECT id, name FROM halls
	WHERE id = $1
	`

	hall := &models.Hall{}

	err := r.pool.QueryRow(context.Background(), query, id).Scan(
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

func (r *HallRepository) CreateHall(params models.CreateHallParams) (*models.Hall, error) {
	const op = "postgres.CreateHall"

	query := `
	INSERT INTO halls (name)
	VALUES ($1)
	RETURNING id, name
	`

	hall := &models.Hall{}
	err := r.pool.QueryRow(context.Background(), query, params.Name).Scan(
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

func (r *HallRepository) UpdateHall(id int64, params models.UpdateHallParams) (*models.Hall, error) {
	const op = "postgres.UpdateHall"

	query := `
	UPDATE halls
	SET 
	    name = COALESCE($1, name)
	WHERE id = $2
	RETURNING id, name
	`

	hall := &models.Hall{}
	err := r.pool.QueryRow(context.Background(), query, params.Name, id).Scan(
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

func (r *HallRepository) DeleteHall(id int64) error {
	const op = "postgres.DeleteHall"
	query := `
	DELETE FROM halls
	WHERE id = $1
	`

	tag, err := r.pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrNotFound)
	}

	return nil
}

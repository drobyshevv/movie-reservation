package postgres

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatRepository struct {
	pool *pgxpool.Pool
}

func NewSeatRepository(pool *pgxpool.Pool) *SeatRepository {
	return &SeatRepository{
		pool: pool,
	}
}

func (r *SeatRepository) DeleteSeat(ctx context.Context, id int64) error {
	const op = "postgres.DeleteSeat"

	qwery := `DELETE FROM seats WHERE id = $1`
	tag, err := r.pool.Exec(ctx, qwery, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrNotFound)
	}

	return nil
}

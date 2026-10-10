package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func runInTx(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		if err := tx.Commit(ctx); err != nil {
			return err
		}
		return nil
	}

	rollbackErr := tx.Rollback(ctx)
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}

func runInTxWithResult[T any](ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) (T, error)) (T, error) {
	var zero T
	tx, err := pool.Begin(ctx)
	if err != nil {
		return zero, err
	}

	item, err := fn(tx)
	if err == nil {
		if err := tx.Commit(ctx); err != nil {
			return zero, err
		}

		return item, nil
	}

	rollbackErr := tx.Rollback(ctx)
	if rollbackErr != nil {
		return zero, errors.Join(err, rollbackErr)
	}

	return zero, err
}

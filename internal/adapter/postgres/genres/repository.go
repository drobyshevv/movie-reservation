package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
	genreuc "github.com/drobyshevv/movie-reservation/internal/usecase/genres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GenreRepository struct {
	pool *pgxpool.Pool
}

func NewGenreRepository(pool *pgxpool.Pool) *GenreRepository {
	return &GenreRepository{
		pool: pool,
	}
}

func (r *GenreRepository) List(ctx context.Context) ([]domain.Genre, error) {
	const op = "postgres.genres.List"

	query := `SELECT id, type FROM genres`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var genres []domain.Genre

	for rows.Next() {
		var id int64
		var typeGenre string

		if err := rows.Scan(&id, &typeGenre); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		genre, err := domain.RestoreGenre(id, typeGenre)
		if err != nil {
			return nil, fmt.Errorf("%s: restore genre: %w", op, err)
		}

		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: iterate rows: %w", op, err)
	}

	return genres, nil
}

func (r *GenreRepository) Get(ctx context.Context, id int64) (*domain.Genre, error) {
	const op = "postgres.genres.Get"

	query := `
		SELECT id, type FROM genres
		WHERE id = $1
	`

	var (
		genreID   int64
		typeGenre string
	)

	err := r.pool.QueryRow(ctx, query, id).Scan(&genreID, &typeGenre)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrGenreNotFound)
		}
		return nil, fmt.Errorf("%s: query row: %w", op, err)
	}

	genre, err := domain.RestoreGenre(genreID, typeGenre)
	if err != nil {
		return nil, fmt.Errorf("%s: restore genre: %w", op, err)
	}

	return &genre, nil
}

func (r *GenreRepository) Create(
	ctx context.Context,
	genre domain.Genre,
) (*domain.Genre, error) {
	const op = "postgres.genres.Create"

	query := `
		INSERT INTO genres (type)
		VALUES ($1)
		RETURNING id, type
	`

	var id int64
	var typeGenre string

	err := r.pool.QueryRow(ctx, query, genre.Name()).Scan(&id, &typeGenre)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrGenreAlreadyExists)
		}
		return nil, fmt.Errorf("%s: insert genre: %w", op, err)
	}

	createdGenre, err := domain.RestoreGenre(id, typeGenre)
	if err != nil {
		return nil, fmt.Errorf("%s: restore genre: %w", op, err)
	}

	return &createdGenre, nil
}

func (r *GenreRepository) Update(
	ctx context.Context,
	id int64,
	genre domain.Genre,
) (*domain.Genre, error) {
	const op = "postgres.genres.Update"

	query := `
		UPDATE genres
		SET type = $1
		WHERE id = $2
		RETURNING id, type
	`

	var updatedID int64
	var typeGenre string

	err := r.pool.QueryRow(ctx, query, genre.Name(), id).Scan(&updatedID, &typeGenre)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrGenreNotFound)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrGenreAlreadyExists)
		}

		return nil, fmt.Errorf("%s: update genre: %w", op, err)
	}

	updatedGenre, err := domain.RestoreGenre(updatedID, typeGenre)
	if err != nil {
		return nil, fmt.Errorf("%s: restore genre: %w", op, err)
	}

	return &updatedGenre, nil
}

func (r *GenreRepository) Delete(ctx context.Context, id int64) error {
	const op = "postgres.genres.Delete"

	query := `
		DELETE FROM genres
		WHERE id = $1
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: delete genre: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errs.ErrGenreNotFound)
	}

	return nil
}

func (r *GenreRepository) GetGenreMovies(
	ctx context.Context,
	id int64,
) ([]genreuc.GenreMovie, error) {
	const op = "postgres.genres.GetGenreMovies"

	query := `
		SELECT
			movies.id,
			movies.title,
			movies.description,
			movies.duration,
			movies.created_at
		FROM movies
		JOIN movie_genre ON movie_genre.movie_id = movies.id
		WHERE movie_genre.genre_id = $1
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: query movies: %w", op, err)
	}
	defer rows.Close()

	movies, err := pgx.CollectRows(rows, pgx.RowToStructByName[genreuc.GenreMovie])
	if err != nil {
		return nil, fmt.Errorf("%s: collect movies: %w", op, err)
	}

	return movies, nil
}

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

type GenreRepository struct {
	pool *pgxpool.Pool
}

func NewGenreRepository(pool *pgxpool.Pool) *GenreRepository {
	return &GenreRepository{
		pool: pool,
	}
}

func (r *GenreRepository) GetGenres() ([]models.Genre, error) {
	const op = "postgres.GetGenres"

	query := `SELECT id, type_genre FROM genres`

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	genres, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Genre])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return genres, nil

}

func (r *GenreRepository) GetGenre(id int64) (*models.Genre, error) {
	const op = "postgres.GetGenre"

	query := `
	SELECT id, type_genre FROM genres
	WHERE id = $1
	`

	genre := &models.Genre{}

	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&genre.ID,
		&genre.TypeGenre,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return genre, nil
}

func (r *GenreRepository) CreateGenre(typeGenre string) (*models.Genre, error) {
	const op = "postgres.CreateGenre"

	query := `
	INSERT INTO genres (type_genre)
	VALUES ($1)
	RETURNING id, type_genre
	`

	genre := &models.Genre{}

	err := r.pool.QueryRow(context.Background(), query, typeGenre).Scan(
		&genre.ID,
		&genre.TypeGenre,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrConflict)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return genre, nil
}

func (r *GenreRepository) UpdateGenre(id int64, typeGenre string) (*models.Genre, error) {
	const op = "postgres.UpdateGenre"

	query := `
	UPDATE genres
	SET type_genre = $1
	WHERE id = $2
	RETURNING id, type_genre
	`

	genre := &models.Genre{}

	err := r.pool.QueryRow(context.Background(), query, typeGenre, id).Scan(
		&genre.ID,
		&genre.TypeGenre,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return genre, nil
}

func (r *GenreRepository) DeleteGenre(id int64) error {
	const op = "postgres.DeleteGenre"

	query := `
	DELETE FROM genres
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

func (r *GenreRepository) GetGenreMovies(id int64) ([]models.Movie, error) {
	const op = "postgres.GetGenreMovies"

	query := `
	SELECT movies.id, movies.title, movies.description, movies.duration, movies.created_at
	FROM movies
	JOIN movie_genre
		ON movie_genre.movie_id = movies.id
	WHERE movie_genre.genre_id = $1
	`

	rows, err := r.pool.Query(context.Background(), query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	movies, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Movie])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

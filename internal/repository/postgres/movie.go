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

type MovieRepository struct {
	pool *pgxpool.Pool
}

func NewMovieRepository(pool *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{
		pool: pool,
	}
}

func (r *MovieRepository) GetAll(ctx context.Context) ([]models.Movie, error) {
	const op = "postgres.GetAll"

	query := `SELECT id, title, description, duration, image, created_at FROM movies`

	rows, err := r.pool.Query(ctx, query)
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

func (r *MovieRepository) Get(ctx context.Context, id int64) (*models.Movie, error) {
	const op = "postgres.Get"

	query := `
	SELECT id, title, description, duration, image, created_at FROM movies
	WHERE id = $1
	`

	movie := &models.Movie{}

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Duration,
		&movie.Image,
		&movie.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movie, nil
}

func (r *MovieRepository) Create(ctx context.Context, params models.CreateMovieParams) (*models.Movie, error) {
	const op = "postgres.Create"

	query := `
	INSERT INTO movies (title, description, duration, image)
	VALUES ($1, $2, $3, $4)	
	RETURNING id, title, description, duration, image, created_at
	`

	movie := models.Movie{}

	err := r.pool.QueryRow(
		ctx,
		query,
		params.Title,
		params.Description,
		params.Duration,
		params.Image,
	).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Duration,
		&movie.Image,
		&movie.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &movie, nil
}

func (r *MovieRepository) Update(ctx context.Context, id int64, params models.UpdateMovieParams) (*models.Movie, error) {
	const op = "postgres.Update"

	query := `
	UPDATE movies
	SET
		title = COALESCE($1, title),
		description = COALESCE($2, description),
		duration = COALESCE($3, duration),
		image = COALESCE($4, image)
	WHERE id = $5
	RETURNING id, title, description, duration, image, created_at;
	`

	movie := models.Movie{}

	err := r.pool.QueryRow(
		ctx,
		query,
		params.Title,
		params.Description,
		params.Duration,
		params.Image,
		id,
	).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Duration,
		&movie.Image,
		&movie.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &movie, err
}

func (r *MovieRepository) Delete(ctx context.Context, id int64) error {
	const op = "postgres.Delete"

	query := `
	DELETE FROM movies
	WHERE id = $1
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *MovieRepository) GetImage(ctx context.Context, id int64) ([]byte, error) {
	const op = "repository.GetImage"

	query := `
	SELECT image FROM movies
	WHERE id = $1
	`

	var image []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(&image)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if image == nil {
		return nil, repository.ErrNotFound
	}

	return image, nil
}

func (r *MovieRepository) PutImage(ctx context.Context, id int64, image []byte) error {
	const op = "repository.PutImage"

	query := `
	UPDATE movies
	SET image = $1
	WHERE id = $2`

	tag, err := r.pool.Exec(ctx, query, image, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *MovieRepository) DeleteImage(ctx context.Context, id int64) error {
	const op = "repository.DeleteImage"

	query := `
	UPDATE movies
	SET image = NULL
	WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *MovieRepository) GetMovieGenres(ctx context.Context, id int64) ([]models.Genre, error) {
	const op = "repository.GetMovieGenres"

	query := `
	SELECT genres.id, genres.type AS type_genre FROM movies
	JOIN movie_genre
	ON movie_genre.movie_id = movies.id
	JOIN genres 
	ON movie_genre.genre_id = genres.id 
	WHERE movies.id = $1
	`

	rows, err := r.pool.Query(context.Background(), query, id)
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

func (r *MovieRepository) PostMovieGenre(ctx context.Context, movieID int64, genreID int64) error {
	const op = "repository.PostMovieGenre"

	query := `
	INSERT INTO movie_genre (movie_id, genre_id)	
	VALUES ($1, $2)
	`

	_, err := r.pool.Exec(ctx, query, movieID, genreID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "movie_genre_movie_id_fkey" {
				return fmt.Errorf("%s: %w", op, repository.ErrMovieNotFound)
			}
			if pgErr.ConstraintName == "movie_genre_genre_id_fkey" {
				return fmt.Errorf("%s: %w", op, repository.ErrGenreNotFound)
			}
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, repository.ErrConflict)
			}
			return err
		}
	}

	return nil
}

func (r *MovieRepository) DeleteMovieGenre(ctx context.Context, movieID int64, genreID int64) error {
	const op = "repository.DeleteMovieGenre"

	query := `
	DELETE FROM movie_genre
	WHERE movie_id = $1
	AND genre_id = $2
	`

	tag, err := r.pool.Exec(ctx, query, movieID, genreID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

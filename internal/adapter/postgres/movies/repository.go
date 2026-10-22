package movies

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/drobyshevv/movie-reservation/internal/adapter/postgres"
	"github.com/drobyshevv/movie-reservation/internal/domain"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
	movieuc "github.com/drobyshevv/movie-reservation/internal/usecase/movie"
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

func (r *MovieRepository) List(ctx context.Context) ([]domain.Movie, error) {
	const op = "postgres.movies.List"

	query := `
	SELECT 
		id, 
		title, 
		description, 
		duration,
		created_at 
	FROM movies`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	movies := []domain.Movie{}
	for rows.Next() {
		var movieID int64
		var title string
		var description *string
		var duration int32
		var createdAt time.Time
		if err := rows.Scan(
			&movieID,
			&title,
			&description,
			&duration,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		movie, err := domain.RestoreMovie(
			movieID,
			title,
			description,
			time.Duration(duration)*time.Minute,
			createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: restore: %w", op, err)
		}

		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: iterate rows: %w", op, err)
	}

	return movies, nil
}

func (r *MovieRepository) Get(ctx context.Context, id int64) (*domain.Movie, error) {
	const op = "postgres.movies.Get"

	query := `
	SELECT 
		id, 
		title, 
		description, 
		duration, 
		created_at 
	FROM movies
	WHERE id = $1
	`

	var (
		movieID     int64
		title       string
		description *string
		duration    int32
		createdAt   time.Time
	)

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&movieID,
		&title,
		&description,
		&duration,
		&createdAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrMovieNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	movie, err := domain.RestoreMovie(
		movieID,
		title,
		description,
		time.Duration(duration)*time.Minute,
		createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: restore: %w", op, err)
	}

	return &movie, nil
}

func (r *MovieRepository) Create(ctx context.Context, input movieuc.CreateMovieInput) (*domain.Movie, error) {
	const op = "postgres.movies.Create"

	query := `
	INSERT INTO movies (title, description, duration)
	VALUES ($1, $2, $3)	
	RETURNING id, title, description, duration, created_at
	`

	var (
		movieID     int64
		title       string
		description *string
		duration    int32
		createdAt   time.Time
	)

	err := r.pool.QueryRow(
		ctx,
		query,
		input.Title,
		input.Description,
		input.Duration,
	).Scan(
		&movieID,
		&title,
		&description,
		&duration,
		&createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	movie, err := domain.RestoreMovie(
		movieID,
		title,
		description,
		time.Duration(duration)*time.Minute,
		createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: restore: %w", op, err)
	}

	return &movie, nil
}

func (r *MovieRepository) Update(ctx context.Context, id int64, updateFn func(*domain.Movie) error) (*domain.Movie, error) {
	return postgres.RunInTxWithResult(ctx, r.pool, func(tx pgx.Tx) (*domain.Movie, error) {
		const op = "postgres.movies.Update"

		var (
			movieID     int64
			title       string
			description *string
			duration    int32
			createdAt   time.Time
		)

		err := tx.QueryRow(
			ctx, `
		SELECT 
			id, 
			title, 
			description, 
			duration, 
			created_at 
		FROM movies
		WHERE id = $1
		FOR UPDATE
			`,
			id,
		).Scan(
			&movieID,
			&title,
			&description,
			&duration,
			&createdAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, errs.ErrMovieNotFound)
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		var hasSessions bool

		err = tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM sessions
				WHERE movie_id = $1
				AND status IN ('scheduled', 'active')
			)
		`, id).Scan(&hasSessions)
		if err != nil {
			return nil, fmt.Errorf("%s: check sessions: %w", op, err)
		}
		if hasSessions == true {
			return nil, fmt.Errorf("%s: : %w", op, errs.ErrMovieHasNonCancelledSessions)
		}

		movie, err := domain.RestoreMovie(
			movieID,
			title,
			description,
			time.Duration(duration)*time.Minute,
			createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: restore: %w", op, err)
		}

		err = updateFn(&movie)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		tag, err := tx.Exec(
			ctx, `
			UPDATE movie
			SET
				title = $1,
				description = $2,
				duration = $3
			WHERE id = $4	
			`,
			movie.Title(),
			movie.Description(),
			movie.Duration(),
			movie.ID(),
		)
		if err != nil {
			return nil, fmt.Errorf("%s: update movie: %w", op, err)
		}
		if tag.RowsAffected() == 0 {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrMovieNotFound)
		}

		return &movie, nil
	})
}

func (r *MovieRepository) Delete(ctx context.Context, id int64) error {
	return postgres.RunInTx(ctx, r.pool, func(tx pgx.Tx) error {
		const op = "postgres.movies.Delete"

		var lockedID int64

		err := tx.QueryRow(
			ctx, `
			SELECT id
			FROM movies
			WHERE id = $1
			FOR UPDATE;
		`,
			id).Scan(&lockedID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%s: %w", op, errs.ErrMovieNotFound)
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		var hasSessions bool

		err = tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM sessions
				WHERE movie_id = $1
				AND status IN ('scheduled', 'active')
			)
		`, lockedID).Scan(&hasSessions)
		if err != nil {
			return fmt.Errorf("%s: check sessions: %w", op, err)
		}
		if hasSessions == true {
			return fmt.Errorf("%s: : %w", op, errs.ErrMovieHasNonCancelledSessions)
		}

		_, err = tx.Exec(ctx,
			`
			DELETE FROM movie
			WHERE id = $1
		`,
			lockedID,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (r *MovieRepository) GetWithGenres(
	ctx context.Context,
	movieID int64,
) (*movieuc.MovieWithGenres, error) {
	const op = "postgres.movies.GetWithGenres"

	query := `
	SELECT
		m.id,
		m.title,
		m.description,
		m.duration,
		m.created_at,
		g.id AS genre_id,
		g.type AS genre_type
	FROM movies m
	LEFT JOIN movie_genre mg ON mg.movie_id = m.id
	LEFT JOIN genres g ON g.id = mg.genre_id
	WHERE m.id = $1
	`

	rows, err := r.pool.Query(ctx, query, movieID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrMovieNotFound)
		}
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var result movieuc.MovieWithGenres
	var movieLoaded bool

	for rows.Next() {
		var (
			id          int64
			title       string
			description *string
			duration    int32
			createdAt   time.Time
			genreID     *int64
			genreType   *string
		)

		if err := rows.Scan(
			&id,
			&title,
			&description,
			&duration,
			&createdAt,
			&genreID,
			&genreType,
		); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		if !movieLoaded {
			movie, err := domain.RestoreMovie(
				id,
				title,
				description,
				time.Duration(duration)*time.Minute,
				createdAt,
			)
			if err != nil {
				return nil, fmt.Errorf("%s: restore movie: %w", op, err)
			}

			result.Movie = movie
			result.Genres = make([]domain.Genre, 0)
			movieLoaded = true
		}

		if genreID != nil && genreType != nil {
			genre, err := domain.RestoreGenre(*genreID, *genreType)
			if err != nil {
				return nil, fmt.Errorf("%s: restore genre: %w", op, err)
			}

			result.Genres = append(result.Genres, genre)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows: %w", op, err)
	}

	if !movieLoaded {
		return nil, errs.ErrMovieNotFound
	}

	return &result, nil
}

func (r *MovieRepository) CreateMovieGenre(ctx context.Context, movieID int64, genreID int64) error {
	const op = "postgres.movies.PostMovieGenre"

	query := `
	INSERT INTO movie_genre (movie_id, genre_id)	
	VALUES ($1, $2)
	`

	_, err := r.pool.Exec(ctx, query, movieID, genreID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "movie_genre_movie_id_fkey" {
				return fmt.Errorf("%s: %w", op, errs.ErrMovieNotFound)
			}
			if pgErr.ConstraintName == "movie_genre_genre_id_fkey" {
				return fmt.Errorf("%s: %w", op, errs.ErrGenreNotFound)
			}
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, errs.ErrMovieGenreAlreadyExists)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *MovieRepository) DeleteMovieGenre(ctx context.Context, movieID int64, genreID int64) error {
	const op = "postgres.movies.DeleteMovieGenre"

	query := `
	DELETE FROM movie_genre
	WHERE movie_id = $1
	AND genre_id = $2
	`

	tag, err := r.pool.Exec(ctx, query, movieID, genreID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errs.ErrMovieGenreNotFound)
	}

	return nil
}

func (r *MovieRepository) GetDuration(ctx context.Context, movieID int64) (time.Duration, error) {
	const op = "postgres.movies.GetDuration"

	query := `
		SELECT duration
		FROM movies
		WHERE id = $1
	`

	var minutes int32

	err := r.pool.QueryRow(ctx, query, movieID).Scan(&minutes)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errs.ErrMovieNotFound
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return time.Duration(minutes) * time.Minute, nil
}

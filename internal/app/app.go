package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/adapter/postgres"

	repoGenres "github.com/drobyshevv/movie-reservation/internal/adapter/postgres/genres"
	repoHalls "github.com/drobyshevv/movie-reservation/internal/adapter/postgres/halls"
	repoMovies "github.com/drobyshevv/movie-reservation/internal/adapter/postgres/movies"
	repoSessions "github.com/drobyshevv/movie-reservation/internal/adapter/postgres/sessions"

	"github.com/drobyshevv/movie-reservation/internal/config"

	"github.com/drobyshevv/movie-reservation/internal/controller/api"
	handlerGenres "github.com/drobyshevv/movie-reservation/internal/controller/api/genres"
	handlerHalls "github.com/drobyshevv/movie-reservation/internal/controller/api/halls"
	handlerMovies "github.com/drobyshevv/movie-reservation/internal/controller/api/movies"
	handlerSessions "github.com/drobyshevv/movie-reservation/internal/controller/api/sessions"

	ucGenres "github.com/drobyshevv/movie-reservation/internal/usecase/genres"
	ucHalls "github.com/drobyshevv/movie-reservation/internal/usecase/halls"
	ucMovies "github.com/drobyshevv/movie-reservation/internal/usecase/movie" // Обрати внимание: папка называется movie, а не movies
	ucSessions "github.com/drobyshevv/movie-reservation/internal/usecase/sessions"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	HTTPServer *http.Server
	Pool       *pgxpool.Pool
}

func NewApp(log *slog.Logger, cfg config.Config) (*App, error) {
	pool, err := postgres.New(context.Background(), config.BuildDSN(cfg.DB))
	if err != nil {
		return nil, err
	}

	hallRepository := repoHalls.NewHallRepository(pool)
	genreRepository := repoGenres.NewGenreRepository(pool)
	movieRepository := repoMovies.NewMovieRepository(pool)
	sessionRepository := repoSessions.NewSessionRepository(pool)

	hallUseCase := ucHalls.NewHallUseCase(hallRepository)
	genreUseCase := ucGenres.NewGenreUseCase(genreRepository)
	movieUseCase := ucMovies.NewMovieUseCase(movieRepository)
	sessionUseCase := ucSessions.NewSessionUseCase(sessionRepository, movieRepository)

	validator := validator.New()

	hallHandler := handlerHalls.NewHallHandler(log, hallUseCase, validator)
	genreHandler := handlerGenres.NewGenreHandler(log, genreUseCase, validator)
	movieHandler := handlerMovies.NewMovieHandler(log, movieUseCase, validator)
	sessionHandler := handlerSessions.NewSessionHandler(log, sessionUseCase, validator)

	router := api.NewRouter(log, hallHandler, genreHandler, movieHandler, sessionHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	return &App{
		HTTPServer: server,
		Pool:       pool,
	}, nil
}

func (a *App) Run() error {
	const op = "app.Run"
	err := a.HTTPServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	const op = "app.Stop"

	err := a.HTTPServer.Shutdown(ctx)
	if err != nil {
		a.HTTPServer.Close()
		a.Pool.Close()
		return fmt.Errorf("%s: %w", op, err)
	}

	a.Pool.Close()

	return nil
}

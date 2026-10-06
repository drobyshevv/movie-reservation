package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/config"
	"github.com/drobyshevv/movie-reservation/internal/handler"
	"github.com/drobyshevv/movie-reservation/internal/handler/middleware/logger"
	"github.com/drobyshevv/movie-reservation/internal/repository/postgres"
	"github.com/drobyshevv/movie-reservation/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	hallRepo := postgres.NewHallRepository(pool)
	genreRepo := postgres.NewGenreRepository(pool)
	movieRepo := postgres.NewMovieRepository(pool)

	hallServ := service.NewHallService(hallRepo)
	genreServ := service.NewGenreService(genreRepo)
	movieServ := service.NewMovieService(movieRepo)

	validator := validator.New()

	hallHand := handler.NewHallHandler(log, hallServ, validator)
	genreHand := handler.NewGenreHandler(log, genreServ, validator)
	movieHand := handler.NewMovieHandler(log, movieServ, validator)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/halls", func(r chi.Router) {
		r.Get("/", hallHand.GetHalls)
		r.Post("/", hallHand.PostHall)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", hallHand.GetHall)
			r.Patch("/", hallHand.PatchHall)
			r.Delete("/", hallHand.DeleteHall)
		})
	})

	router.Route("/genres", func(r chi.Router) {
		r.Get("/", genreHand.GetGenres)
		r.Post("/", genreHand.PostGenre)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", genreHand.GetGenre)
			r.Put("/", genreHand.PutGenre)
			r.Delete("/", genreHand.DeleteGenre)
			r.Get("/movies", genreHand.GetGenreMovies)
		})

	})

	router.Route("/movies", func(r chi.Router) {
		r.Get("/", movieHand.GetMovies)
		r.Post("/", movieHand.PostMovie)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", movieHand.GetMovie)
			r.Patch("/", movieHand.PatchMovies)
			r.Delete("/", movieHand.DeleteMovies)
			r.Get("/image", movieHand.GetImage)
			r.Put("/image", movieHand.PutImage)
			r.Delete("/image", movieHand.DeleteImage)
		})
	})

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

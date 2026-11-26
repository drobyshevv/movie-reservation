package api

import (
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/bookings"
	"github.com/drobyshevv/movie-reservation/internal/controller/api/genres"
	"github.com/drobyshevv/movie-reservation/internal/controller/api/halls"
	"github.com/drobyshevv/movie-reservation/internal/controller/api/movies"
	"github.com/drobyshevv/movie-reservation/internal/controller/api/mw"
	"github.com/drobyshevv/movie-reservation/internal/controller/api/sessions"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	log *slog.Logger,
	hallHandler *halls.HallHandler,
	genreHandler *genres.GenreHandler,
	movieHandler *movies.MovieHandler,
	sessionHandler *sessions.SessionHandler,
	bookingHandler *bookings.BookingHandler,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mw.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/halls", func(r chi.Router) {
		r.Get("/", hallHandler.List)
		r.Post("/", hallHandler.PostWithSeats)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", hallHandler.Get)
			r.Put("/", hallHandler.Put)
			r.Delete("/", hallHandler.Delete)
			r.Route("/seats", func(r chi.Router) {
				r.Get("/", hallHandler.GetWithSeats)
				r.Put("/", hallHandler.PutSeats)
				r.Delete("/", hallHandler.DeleteSeats)
			})
		})
	})

	router.Route("/genres", func(r chi.Router) {
		r.Get("/", genreHandler.List)
		r.Post("/", genreHandler.Post)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", genreHandler.Get)
			r.Put("/", genreHandler.Put)
			r.Delete("/", genreHandler.Delete)
			r.Get("/movies", genreHandler.GetGenreMovies)
		})

	})

	router.Route("/movies", func(r chi.Router) {
		r.Get("/", movieHandler.List)
		r.Post("/", movieHandler.Post)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", movieHandler.Get)
			r.Patch("/", movieHandler.Patch)
			r.Delete("/", movieHandler.Delete)
			r.Get("/genres", movieHandler.GetMovieGenres)
			r.Post("/genres/{genre_id}", movieHandler.PostMovieGenre)
			r.Delete("/genres/{genre_id}", movieHandler.DeleteMovieGenre)

			r.Get("/poster", movieHandler.GetPoster)
			r.Put("/poster", movieHandler.PutPoster)
			r.Delete("/poster", movieHandler.DeletePoster)
		})
	})

	router.Route("/sessions", func(r chi.Router) {
		r.Get("/", sessionHandler.List)
		r.Post("/", sessionHandler.Post)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", sessionHandler.Get)
			r.Patch("/", sessionHandler.Patch)
			r.Delete("/", sessionHandler.Delete)

			r.Get("/bookings", bookingHandler.ListBySession)
		})

	})

	router.Route("/bookings", func(r chi.Router) {
		r.Get("/{id}", bookingHandler.Get)
		r.Post("/", bookingHandler.Post)
		r.Patch("/{id}", bookingHandler.Patch)
	})

	router.Route("/users", func(r chi.Router) {
		r.Get("/{user_uid}/bookings", bookingHandler.ListByUser)
	})

	return router
}

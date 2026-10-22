package http

import (
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/controller/http/genres"
	"github.com/drobyshevv/movie-reservation/internal/controller/http/halls"
	"github.com/drobyshevv/movie-reservation/internal/controller/http/movies"
	"github.com/drobyshevv/movie-reservation/internal/controller/http/sessions"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	hallHandler *halls.HallHandler,
	genreHandler *genres.GenreHandler,
	movieHandler *movies.MovieHandler,
	sessionHandler *sessions.SessionHandler,
) http.Handler {
	router := chi.NewRouter()

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
			r.Get("/genres", movieHandler.GetMovieGenre)
			r.Post("/genres/{genre_id}", movieHandler.PostMovieGenre)
			r.Delete("/genres/{genre_id}", movieHandler.DeleteMovieGenre)
		})
	})

	router.Route("/sessions", func(r chi.Router) {
		r.Get("/", sessionHandler.List)
		r.Post("/", sessionHandler.Post)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", sessionHandler.Get)
			r.Patch("/", sessionHandler.Patch)
			r.Delete("/", sessionHandler.Delete)
		})
	})

	return router
}

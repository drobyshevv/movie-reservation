package http

import (
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/controller/http/genres"
	"github.com/drobyshevv/movie-reservation/internal/controller/http/sessions"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	sessionHandler *sessions.SessionHandler,
	genreHandler *genres.GenreHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Route("/genres", func(r chi.Router) {
		r.Get("/", genreHandler.List)
		r.Post("/", genreHandler.Post)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", genreHandler.Get)
			r.Put("/", genreHandler.Put)
			r.Delete("/", genreHandler.Delete)
			r.Get("/movies", genreHandler.GetGenreMovies)
		})

	})

	r.Route("/sessions", func(r chi.Router) {
		r.Get("/", sessionHandler.List)
		r.Post("/", sessionHandler.Post)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", sessionHandler.Get)
			r.Patch("/", sessionHandler.Update)
			r.Delete("/", sessionHandler.Delete)
		})
	})

	return r
}

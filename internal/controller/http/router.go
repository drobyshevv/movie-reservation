package http

import (
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/controller/http/sessions"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	sessionHandler *sessions.SessionHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Route("/sessions", func(r chi.Router) {
		r.Get("/", sessionHandler.Get)
		r.Post("/", sessionHandler.Create)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", sessionHandler.Get)
			r.Patch("/", sessionHandler.Update)
			r.Delete("/", sessionHandler.Delete)
		})
	})

	return r
}

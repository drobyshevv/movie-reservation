package movies

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/http/response"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *MovieHandler) DeleteMovieGenre(w http.ResponseWriter, r *http.Request) {
	const op = "http.movies.DeleteMovieGenre"
	log := h.log.With(slog.String("op", op))

	movieIDStr := r.PathValue("id")
	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "movie_id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid movie ID", http.StatusBadRequest, "Movie ID must be a valid integer")
		return
	}

	genreIDStr := r.PathValue("genre_id")
	genreID, err := strconv.ParseInt(genreIDStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "genre_id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid genre ID", http.StatusBadRequest, "Genre ID must be a valid integer")
		return
	}

	err = h.usecase.DeleteMovieGenre(r.Context(), movieID, genreID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrMovieGenreNotFound):
			log.Error("movie genre link not found", "err", err)
			response.WriteError(w, "not_found", "Movie genre link not found", http.StatusNotFound, err.Error())
		default:
			log.Error("failed to delete movie genre", "err", err)
			response.WriteError(w, "internal_error", "Failed to delete movie genre", http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

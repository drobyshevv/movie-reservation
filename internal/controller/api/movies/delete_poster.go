package movies

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *MovieHandler) DeletePoster(w http.ResponseWriter, r *http.Request) {
	const op = "http.movies.DeletePoster"

	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid ID", http.StatusBadRequest, "ID must be a valid integer")
		return
	}

	err = h.usecase.DeletePoster(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrMovieNotFound):
			log.Warn("movie not found", "err", err)
			response.WriteError(w, "not_found", "Movie not found", http.StatusNotFound, "")
		default:
			log.Error("failed to get poster", "err", err)
			response.WriteError(w, "internal_error", "Failed to get poster", http.StatusInternalServerError, "")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

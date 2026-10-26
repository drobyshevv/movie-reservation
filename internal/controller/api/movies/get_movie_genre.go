package movies

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *MovieHandler) GetMovieGenre(w http.ResponseWriter, r *http.Request) {
	const op = "http.movies.GetWithGenres"
	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid ID", http.StatusBadRequest, "ID must be a valid integer")
		return
	}

	result, err := h.usecase.GetWithGenres(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrMovieNotFound):
			log.Warn("movie not found", "err", err)
			response.WriteError(w, "not_found", "Movie not found", http.StatusNotFound, err.Error())
		default:
			log.Error("failed to get movie with genres", "err", err)
			response.WriteError(w, "internal_error", "Failed to get movie with genres", http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := dto.MovieWithGenresResponse{}
	resp.FromModel(*result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

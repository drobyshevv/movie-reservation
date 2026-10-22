package genres

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

func (h *GenreHandler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "http.genres.Get"
	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("failed to parse id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid ID format", http.StatusBadRequest, "ID must be a valid integer")
		return
	}

	result, err := h.usecase.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrGenreNotFound):
			log.Error("genre not found", "err", err)
			response.WriteError(w, "not_found", "Genre not found", http.StatusNotFound, err.Error())
		default:
			log.Error("failed to get genre", "err", err)
			response.WriteError(w, "internal_error", "Failed to get genre", http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := dto.GenreResponse{}
	resp.FromModel(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

package movies

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/http/response"
	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *MovieHandler) Patch(w http.ResponseWriter, r *http.Request) {
	const op = "http.movies.Patch"
	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid ID", http.StatusBadRequest, "ID must be a valid integer")
		return
	}

	var req dto.UpdateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", "err", err)
		response.WriteError(w, "invalid_request", "Invalid request body", http.StatusBadRequest, "Failed to decode JSON")
		return
	}

	if req.Title == nil && req.Description == nil && req.Duration == nil {
		log.Warn("empty patch request")
		response.WriteError(w, "validation_error", "Empty patch request", http.StatusBadRequest, "At least one field must be provided for update")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		log.Error("validation failed", "err", err)
		response.WriteError(w, "validation_error", "Validation failed", http.StatusBadRequest, err.Error())
		return
	}

	input := req.ToInput()

	result, err := h.usecase.Update(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrMovieNotFound):
			log.Warn("movie not found", "err", err)
			response.WriteError(w, "not_found", "Movie not found", http.StatusNotFound, err.Error())
		default:
			log.Error("failed to update movie", "err", err)
			response.WriteError(w, "internal_error", "Failed to update movie", http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := dto.MovieResponse{}
	resp.FromModel(*result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

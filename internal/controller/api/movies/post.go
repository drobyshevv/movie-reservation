package movies

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
	"github.com/drobyshevv/movie-reservation/internal/dto"
)

func (h *MovieHandler) Post(w http.ResponseWriter, r *http.Request) {
	const op = "http.movies.Post"
	log := h.log.With(slog.String("op", op))

	var req dto.CreateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", "err", err)
		response.WriteError(w, "invalid_request", "Invalid request body", http.StatusBadRequest, "Failed to decode JSON")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		log.Error("validation failed", "err", err)
		response.WriteError(w, "validation_error", "Validation failed", http.StatusBadRequest, err.Error())
		return
	}

	input := req.ToInput()

	result, err := h.usecase.Create(r.Context(), input)
	if err != nil {
		log.Error("failed to create movie", "err", err)
		response.WriteError(w, "internal_error", "Failed to create movie", http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.MovieResponse{}
	resp.FromModel(*result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

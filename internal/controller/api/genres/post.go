package genres

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *GenreHandler) Post(w http.ResponseWriter, r *http.Request) {
	const op = "http.genres.Post"
	log := h.log.With(slog.String("op", op))

	var req dto.CreateGenreRequest
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

	result, err := h.usecase.Create(r.Context(), req.TypeGenre)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrGenreAlreadyExists):
			log.Error("genre conflict", "err", err)
			response.WriteError(w, "conflict", "Genre already exists", http.StatusConflict, err.Error())
		default:
			log.Error("failed to create genre", "err", err)
			response.WriteError(w, "internal_error", "Failed to create genre", http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := dto.GenreResponse{}
	resp.FromModel(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

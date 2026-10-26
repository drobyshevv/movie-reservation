package sessions

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *SessionHandler) Post(w http.ResponseWriter, r *http.Request) {
	const op = "http.sessions.Post"
	log := h.log.With(slog.String("op", op))

	var req dto.CreateSessionRequest
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
		switch {
		case errors.Is(err, errs.ErrHallNotFound):
			log.Warn("hall not found", "err", err)
			response.WriteError(w, "not_found", "Hall not found", http.StatusNotFound, err.Error())
		case errors.Is(err, errs.ErrMovieNotFound):
			log.Warn("movie not found", "err", err)
			response.WriteError(w, "not_found", "Movie not found", http.StatusNotFound, err.Error())
		case errors.Is(err, errs.ErrSessionAlreadyExists):
			log.Warn("session already exists", "err", err)
			response.WriteError(w, "conflict", "Session already exists", http.StatusConflict, err.Error())
		default:
			log.Error("failed to create session", "err", err)
			response.WriteError(w, "internal_error", "Failed to create session", http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := dto.SessionResponse{}
	resp.FromModel(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

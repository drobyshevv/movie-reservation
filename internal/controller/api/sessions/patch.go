package sessions

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

func (h *SessionHandler) Patch(w http.ResponseWriter, r *http.Request) {
	const op = "http.sessions.Patch"
	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("failed to parse id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid ID format", http.StatusBadRequest, "ID must be a valid integer")
		return
	}

	var req dto.UpdateSessionRequest
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

	if req.Status == nil && req.StartAt == nil && req.EndAt == nil {
		log.Warn("empty patch request")
		response.WriteError(w, "validation_error", "Empty patch request", http.StatusBadRequest, "At least one field must be provided for update")
		return
	}

	input := req.ToInput()

	result, err := h.usecase.Update(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrSessionNotFound):
			log.Warn("session not found", "err", err)
			response.WriteError(w, "not_found", "Session not found", http.StatusNotFound, err.Error())
		case errors.Is(err, errs.ErrSessionNotCancelled):
			log.Warn("session must be cancelled first", "err", err)
			response.WriteError(w, "conflict", "Session must be cancelled before modification", http.StatusConflict, err.Error())
		case errors.Is(err, errs.ErrConflictingSessionUpdates):
			log.Warn("conflicting session updates", "err", err)
			response.WriteError(w, "conflict", "Conflicting session updates", http.StatusConflict, err.Error())
		default:
			log.Error("failed to update session", "err", err)
			response.WriteError(w, "internal_error", "Failed to update session", http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := dto.SessionResponse{}
	resp.FromModel(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

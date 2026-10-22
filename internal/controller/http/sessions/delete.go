package sessions

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/http/response"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *SessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "http.sessions.Delete"
	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("failed to parse id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid ID format", http.StatusBadRequest, "ID must be a valid integer")
		return
	}

	err = h.usecase.Delete(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrSessionNotFound):
			log.Warn("session not found", "err", err)
			response.WriteError(w, "not_found", "Session not found", http.StatusNotFound, err.Error())
		case errors.Is(err, errs.ErrSessionNotCancelled):
			log.Warn("cannot delete active session", "err", err)
			response.WriteError(w, "conflict", "Session must be cancelled before deletion", http.StatusConflict, err.Error())
		default:
			log.Error("failed to delete session", "err", err)
			response.WriteError(w, "internal_error", "Failed to delete session", http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

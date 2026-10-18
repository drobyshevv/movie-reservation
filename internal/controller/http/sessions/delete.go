package sessions

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *SessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "http.sessions.Delete"

	log := h.log.With(
		slog.String("op", op),
	)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("failed to parse id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usecase.DeleteSession(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrSessionNotFound) {
			log.Error("session not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, errs.ErrSessionNotCancelled) {
			log.Error("cannot delete session that is not cancelled", "err", err)
			http.Error(w, "session must be cancelled before deletion", http.StatusConflict)
			return
		}
		log.Error("failed to update session", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

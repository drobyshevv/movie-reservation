package sessions

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/handler/dto"
)

func (h *SessionHandler) List(w http.ResponseWriter, r *http.Request) {
	const op = "http.sessions.List"

	log := h.log.With(
		slog.String("op", op),
	)

	sessions, err := h.service.GetSessions(r.Context())
	if err != nil {
		log.Error("failed to get halls", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := dto.SessionWithDetailsListFromModels(sessions)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

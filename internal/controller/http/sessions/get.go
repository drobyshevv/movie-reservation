package sessions

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	errs "github.com/drobyshevv/movie-reservation/internal/errors"
	"github.com/drobyshevv/movie-reservation/internal/handler/dto"
)

func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "http.sessions.Get"

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

	session, err := h.usecase.GetSession(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrSessionNotFound) {
			log.Info("session not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("failed to get session", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.SessionWithDetailsResponse{}
	resp.FromModel(session)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

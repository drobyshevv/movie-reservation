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

func (h *SessionHandler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "http.sessions.Update"

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

	req := dto.UpdateSessionRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Info("failed to decode request body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Info("invalid request", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Status == nil &&
		req.StartAt == nil &&
		req.EndAt == nil {
		err := errors.New("at least one field must be provided")
		log.Info("invalid request", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := req.ToModel()

	session, err := h.usecase.UpdateSession(r.Context(), id, input)
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
		if errors.Is(err, errs.ErrConflictingSessionUpdates) {
			log.Error("conflicting session updates", "err", err)
			http.Error(w, "conflicting session updates", http.StatusConflict)
			return
		}
		log.Error("failed to update session", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.SessionResponse{}
	resp.FromModel(session)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

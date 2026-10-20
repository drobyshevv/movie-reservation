package sessions

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errs "github.com/drobyshevv/movie-reservation/internal/errors"
	"github.com/drobyshevv/movie-reservation/internal/handler/dto"
)

func (h *SessionHandler) Post(w http.ResponseWriter, r *http.Request) {
	const op = "http.sessions.Post"

	log := h.log.With(
		slog.String("op", op),
	)

	req := dto.CreateSessionRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
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

	input := req.ToModel()

	session, err := h.usecase.CreateSession(r.Context(), input)
	if err != nil {
		if errors.Is(err, errs.ErrHallNotFound) {
			log.Error("hall not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, errs.ErrMovieNotFound) {
			log.Error("movie not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, errs.ErrSessionAlreadyExists) {
			log.Error("session already exists", "err", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		log.Error("failed to create session", "err", err)
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

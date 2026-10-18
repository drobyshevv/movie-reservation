package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	"github.com/drobyshevv/movie-reservation/internal/handler/dto"
	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/drobyshevv/movie-reservation/internal/service"
	"github.com/go-playground/validator/v10"
)

type SessionHandler struct {
	log       *slog.Logger
	service   SessionService
	validator *validator.Validate
}

func NewSessionHandler(log *slog.Logger, service SessionService, validator *validator.Validate) *SessionHandler {
	return &SessionHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

type SessionService interface {
	GetSessions(ctx context.Context) ([]models.SessionWithDetails, error)
	GetSession(ctx context.Context, id int64) (*models.SessionWithDetails, error)
	CreateSession(ctx context.Context, params models.CreateSessionInput) (*domain.Session, error)
	UpdateSession(ctx context.Context, id int64, params models.UpdateSessionInput) (*domain.Session, error)
	DeleteSession(ctx context.Context, id int64) error
}

func (h *SessionHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetSessions"

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

func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetSession"

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

	session, err := h.service.GetSession(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
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

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateSession"

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

	params := req.ToModel()

	session, err := h.service.CreateSession(r.Context(), params)
	if err != nil {
		if errors.Is(err, service.ErrHallNotFound) {
			log.Error("hall not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrMovieNotFound) {
			log.Error("movie not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrSessionAlreadyExists) {
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

func (h *SessionHandler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	const op = "handler.UpdateSession"

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

	params := req.ToModel()

	session, err := h.service.UpdateSession(r.Context(), id, params)
	if err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
			log.Error("session not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrSessionNotCancelled) {
			log.Error("cannot delete session that is not cancelled", "err", err)
			http.Error(w, "session must be cancelled before deletion", http.StatusConflict)
			return
		}
		if errors.Is(err, service.ErrConflictingSessionUpdates) {
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

func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteSession"

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

	err = h.service.DeleteSession(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
			log.Error("session not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrSessionNotCancelled) {
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

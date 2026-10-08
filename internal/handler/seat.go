package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/service"
	"github.com/go-playground/validator/v10"
)

type SeatHandler struct {
	log       *slog.Logger
	service   SeatService
	validator *validator.Validate
}

func NewSeatHandler(log *slog.Logger, service SeatService, validator *validator.Validate) *SeatHandler {
	return &SeatHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

type SeatService interface {
	DeleteSeat(ctx context.Context, id int64) error
}

func (h *SeatHandler) DeleteSeat(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteSeat"
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

	err = h.service.DeleteSeat(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrSeatNotFound) {
			log.Info("seat not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

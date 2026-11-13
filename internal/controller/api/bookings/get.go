package bookings

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *BookingHandler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "http.bookings.Get"

	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("failed to parse id", "err", err)
		response.WriteError(w, response.CodeInvalidRequest, "Invalid ID format", http.StatusBadRequest, "ID must be a valid integer")
		return
	}

	result, err := h.usecase.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrBookingNotFound):
			log.Warn("hall not found", "err", err)
			response.WriteError(w, response.CodeNotFound, "Hall not found", http.StatusNotFound, "")
		default:
			log.Error("failed to get booking", "err", err)
			response.WriteError(w, response.CodeInternalError, "Internal server error", http.StatusInternalServerError, "")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

package bookings

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
	"github.com/drobyshevv/movie-reservation/internal/dto"
	"github.com/google/uuid"
)

func (h *BookingHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	const op = "http.bookings.ListByUser"

	log := h.log.With(slog.String("op", op))

	userIDStr := r.PathValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Info("invalid user ID format", "raw_user_id", userIDStr, "err", err)
		response.WriteError(w, response.CodeInvalidRequest, "Invalid user ID format", http.StatusBadRequest, "")
		return
	}

	result, err := h.usecase.ListByUser(r.Context(), userID)
	if err != nil {
		log.Error("failed to get list bookings by user", "err", err)
		response.WriteError(w, response.CodeInternalError, "Internal server error", http.StatusInternalServerError, "")
		return
	}

	resp := dto.ToUserBookings(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

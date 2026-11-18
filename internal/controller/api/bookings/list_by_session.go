package bookings

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/dto"
	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
)

func (h *BookingHandler) ListBySession(w http.ResponseWriter, r *http.Request) {
	const op = "http.bookings.ListBySession"

	log := h.log.With(slog.String("op", op))

	sessionIdStr := r.PathValue("session_id")
	sessionId, err := strconv.ParseInt(sessionIdStr, 10, 64)
	if err != nil {
		log.Info("failed to parse session_id", "err", err)
		response.WriteError(w, response.CodeInvalidRequest, "Invalid ID format", http.StatusBadRequest, "")
		return
	}

	result, err := h.usecase.ListBySession(r.Context(), sessionId)
	if err != nil {
		log.Error("failed to get list bookings by session", "err", err)
		response.WriteError(w, response.CodeInternalError, "Internal server error", http.StatusInternalServerError, "")
		return
	}

	resp := dto.ToSeatStatusesResponse(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

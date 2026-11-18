package bookings

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
	"github.com/drobyshevv/movie-reservation/internal/domain"
	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *BookingHandler) Patch(w http.ResponseWriter, r *http.Request) {
	const op = "http.booking.Patch"

	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Info("failed to parse id", "err", err)
		response.WriteError(w, response.CodeInvalidRequest, "Invalid ID format", http.StatusBadRequest, "")
		return
	}

	req := dto.UpdateStatusRequest{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Info("failed to decode update status request", "err", err)
		response.WriteError(w, response.CodeInvalidRequest, "Invalid JSON format", http.StatusBadRequest, "")
		return
	}

	input, err := req.ToInput()
	if err != nil {
		log.Info("received invalid booking status value", "raw_status", req.Status)
		response.WriteError(w, response.CodeValidationError, "Invalid status", http.StatusBadRequest, "")
		return
	}

	err = h.usecase.UpdateStatus(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrBookingNotFound):
			log.Warn("booking not found", "booking_id", id)
			response.WriteError(w, response.CodeNotFound, "Booking not found", http.StatusNotFound, "")

		case errors.Is(err, domain.ErrBookingNotPending):
			log.Warn("cannot confirm: booking is not pending", "raw_status", input)
			response.WriteError(w, response.CodeConflict, "Booking must be pending to confirm", http.StatusConflict, "")

		case errors.Is(err, domain.ErrBookingNotConfirmed):
			log.Warn("cannot complete: booking is not confirmed", "booking_id", id)
			response.WriteError(w, response.CodeConflict, "Booking must be confirmed to complete", http.StatusConflict, "")

		case errors.Is(err, domain.ErrBookingAlreadyCancelled):
			log.Info("duplicate cancel attempt", "booking_id", id)
			response.WriteError(w, response.CodeConflict, "Booking is already cancelled", http.StatusConflict, "")

		case errors.Is(err, domain.ErrBookingAlreadyCompleted):
			log.Info("cannot cancel completed booking", "booking_id", id)
			response.WriteError(w, response.CodeConflict, "Cannot cancel a completed booking", http.StatusConflict, "")

		case errors.Is(err, domain.ErrInvalidBookingStatus):
			log.Info("received invalid booking status value", "raw_status", input)
			response.WriteError(w, response.CodeValidationError, "Invalid status value", http.StatusBadRequest, "")

		default:
			log.Error("failed to update booking status", "err", err)
			response.WriteError(w, response.CodeInternalError, "Internal server error", http.StatusInternalServerError, "")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

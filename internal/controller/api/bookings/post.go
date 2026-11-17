package bookings

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *BookingHandler) Post(w http.ResponseWriter, r *http.Request) {
	const op = "http.bookings.Post"

	log := h.log.With(slog.With("op", op))

	req := dto.CreateBookingRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Info("failed to decode request body", "err", err)
		response.WriteError(w, response.CodeInvalidRequest, "Invalid request body", http.StatusBadRequest, "")
		return
	}

	if err = h.validator.Struct(req); err != nil {
		log.Info("validation failed", "err", err)
		response.WriteError(w, response.CodeValidationError, "Validation failed", http.StatusBadRequest, "")
		return
	}

	input, err := req.ToInput()

	result, err := h.usecase.Create(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrBookingAlreadyExists):
			log.Warn("booking already exists", "err", err)
			response.WriteError(w, response.CodeConflict, "Booking already exists", http.StatusConflict, "")
		case errors.Is(err, errs.ErrInvalidBookingConditions):
			log.Warn("invalid booking conditions", "err", err)
			response.WriteError(w, response.CodeConflict, "Invalid booking conditions", http.StatusConflict, "")
		default:
			log.Error("failed to create booking", "err", err)
			response.WriteError(w, response.CodeInternalError, "Internal server error", http.StatusInternalServerError, "")
		}
		return
	}

	resp := dto.ToBookingResponse(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

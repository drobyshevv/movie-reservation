package halls

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *HallHandler) PostWithSeats(w http.ResponseWriter, r *http.Request) {
	const op = "http.halls.CreateWithSeats"
	log := h.log.With(slog.String("op", op))

	var req dto.CreateHallWithSeatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", "err", err)
		response.WriteError(w, "invalid_request", "Invalid request body", http.StatusBadRequest, "Failed to decode JSON")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		log.Error("validation failed", "err", err)
		response.WriteError(w, "validation_error", "Validation failed", http.StatusBadRequest, err.Error())
		return
	}

	input, err := req.ToInput()
	if err != nil {
		if errors.Is(err, errs.ErrDuplicateSeatsInRequest) {
			log.Error("invalid request: duplicate seats in payload", "err", err)
			response.WriteError(w, "validation_error", "Duplicate seats in request", http.StatusBadRequest, err.Error())
			return
		}
		log.Error("failed to map seats input", "err", err)
		response.WriteError(w, "internal_error", "Failed to map seats input", http.StatusInternalServerError, err.Error())
		return
	}

	result, err := h.usecase.CreateWithSeats(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrSeatAlreadyExists):
			log.Error("seats already exist in this hall", "err", err)
			response.WriteError(w, "conflict", "Seats already exist", http.StatusConflict, err.Error())
		case errors.Is(err, errs.ErrHallHasNonCancelledSessions):
			log.Error("cannot modify hall with active sessions", "err", err)
			response.WriteError(w, "conflict", "Hall has non-cancelled sessions", http.StatusConflict, err.Error())
		default:
			log.Error("failed to create hall with seats", "err", err)
			response.WriteError(w, "internal_error", "Failed to create hall with seats", http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := dto.HallWithSeatsResponse{}
	resp.FromModel(*result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

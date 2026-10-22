package halls

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/http/response"
	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *HallHandler) PutSeats(w http.ResponseWriter, r *http.Request) {
	const op = "http.halls.PutSeats"
	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("failed to parse id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid ID format", http.StatusBadRequest, "ID must be a valid integer")
		return
	}

	var req []dto.CreateHallSeatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", "err", err)
		response.WriteError(w, "invalid_request", "Invalid request body", http.StatusBadRequest, "Failed to decode JSON")
		return
	}

	if err := h.validator.Var(req, "required,dive"); err != nil {
		log.Error("validation failed", "err", err)
		response.WriteError(w, "validation_error", "Validation failed", http.StatusBadRequest, err.Error())
		return
	}

	input, err := dto.ToSeatsInput(req)
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

	result, err := h.usecase.ReplaceSeats(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrHallNotFound):
			log.Error("hall not found", "err", err)
			response.WriteError(w, "not_found", "Hall not found", http.StatusNotFound, err.Error())
		case errors.Is(err, errs.ErrSeatAlreadyExists):
			log.Error("conflict: seats already exist in database", "err", err)
			response.WriteError(w, "conflict", "Seats already exist in database", http.StatusConflict, err.Error())
		case errors.Is(err, errs.ErrHallHasActiveSessions):
			log.Error("conflict: hall has active sessions", "err", err)
			response.WriteError(w, "conflict", "Hall has active sessions", http.StatusConflict, err.Error())
		default:
			log.Error("failed to replace hall seats", "err", err)
			response.WriteError(w, "internal_error", "Failed to replace hall seats", http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := dto.ToSeatListResponse(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

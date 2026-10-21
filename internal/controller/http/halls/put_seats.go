package halls

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *HallHandler) PutSeats(w http.ResponseWriter, r *http.Request) {
	const op = "http.halls.PutSeats"
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

	var req []dto.CreateHallSeatsRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("failed to decode request body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, item := range req {
		err = h.validator.Struct(item)
		if err != nil {
			log.Error("validation failed", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	input, err := dto.ToSeatsInput(req)
	if err != nil {
		if errors.Is(err, errs.ErrDuplicateSeatsInRequest) {
			log.Error("invalid request: duplicate seats in payload", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Error("failed to map seats input", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := h.usecase.ReplaceSeats(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrHallNotFound):
			log.Error("hall not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		case errors.Is(err, errs.ErrSeatAlreadyExists):
			log.Error("conflict: seats already exist in database", "err", err)
			w.WriteHeader(http.StatusConflict)
			return
		case errors.Is(err, errs.ErrHallHasSessions):
			log.Error("conflict: hall has active sessions", "err", err)
			w.WriteHeader(http.StatusConflict)
			return
		default:
			log.Error("failed to replace hall seats", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	resp := dto.ToSeatListResponse(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

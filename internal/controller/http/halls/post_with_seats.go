package halls

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *HallHandler) CreateWithSeats(w http.ResponseWriter, r *http.Request) {
	const op = "http.halls.CreateWithSeats"
	log := h.log.With(
		slog.String("op", op),
	)

	var req dto.CreateHallWithSeatsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("failed to decode request body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Error("validation failed", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input, err := req.ToInput()
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

	result, err := h.usecase.CreateWithSeats(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrSeatAlreadyExists):
			log.Error("seats already exist in this hall", "err", err)
			w.WriteHeader(http.StatusConflict)
			return
		case errors.Is(err, errs.ErrHallHasSessions):
			log.Error("cannot modify hall with active sessions", "err", err)
			w.WriteHeader(http.StatusConflict)
			return
		default:
			log.Error("failed to create hall with seats", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	resp := dto.HallWithSeatsResponse{}
	resp.FromModel(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

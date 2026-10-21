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

func (h *HallHandler) GetWithSeats(w http.ResponseWriter, r *http.Request) {
	const op = "http.halls.GetWithSeats"
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

	result, err := h.usecase.GetWithSeats(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrHallNotFound):
			log.Error("hall not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			log.Error("failed to get hall with seats", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	resp := dto.HallWithSeatsResponse{}
	resp.FromModel(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

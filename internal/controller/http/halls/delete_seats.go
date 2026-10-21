package halls

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/service"
)

func (h *HallHandler) DeleteSeats(w http.ResponseWriter, r *http.Request) {
	const op = "http.halls.DeleteSeats"
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

	err = h.usecase.DeleteSeats(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrHallNotFound):
			log.Error("hall not found", "err", err)
			http.Error(w, "hall not found", http.StatusNotFound)
			return

		case errors.Is(err, service.ErrSeatNotFound):
			log.Error("seats not found", "err", err)
			http.Error(w, "seats not found", http.StatusNotFound)
			return

		default:
			log.Error("failed to delete hall", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

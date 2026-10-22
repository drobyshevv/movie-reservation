package halls

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/http/response"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *HallHandler) DeleteSeats(w http.ResponseWriter, r *http.Request) {
	const op = "http.halls.DeleteSeats"
	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("failed to parse id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid ID format", http.StatusBadRequest, "ID must be a valid integer")
		return
	}

	err = h.usecase.DeleteSeats(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrHallNotFound):
			log.Error("hall not found", "err", err)
			response.WriteError(w, "not_found", "Hall not found", http.StatusNotFound, err.Error())
		case errors.Is(err, errs.ErrSeatNotFound):
			log.Error("seats not found", "err", err)
			response.WriteError(w, "not_found", "Seats not found", http.StatusNotFound, err.Error())
		default:
			log.Error("failed to delete hall seats", "err", err)
			response.WriteError(w, "internal_error", "Failed to delete hall seats", http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

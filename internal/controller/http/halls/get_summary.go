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

func (h *HallHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	const op = "http.halls.GetSummary"
	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("failed to parse id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid ID format", http.StatusBadRequest, "ID must be a valid integer")
		return
	}

	result, err := h.usecase.GetSummary(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrHallNotFound):
			log.Error("hall not found", "err", err)
			response.WriteError(w, "not_found", "Hall not found", http.StatusNotFound, err.Error())
		default:
			log.Error("failed to get hall summary", "err", err)
			response.WriteError(w, "internal_error", "Failed to get hall summary", http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := dto.HallSummaryResponse{}
	resp.FromModel(*result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

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

func (h *HallHandler) Put(w http.ResponseWriter, r *http.Request) {
	const op = "http.halls.Put"
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

	var req dto.UpdateHallRequest
	err = json.NewDecoder(r.Body).Decode(&req)
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

	result, err := h.usecase.Update(r.Context(), id, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrHallNotFound):
			log.Error("hall not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			log.Error("failed to update hall", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	resp := dto.HallResponse{}
	resp.FromModel(*result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

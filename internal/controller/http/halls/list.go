package halls

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/controller/http/response"
	"github.com/drobyshevv/movie-reservation/internal/dto"
)

func (h *HallHandler) List(w http.ResponseWriter, r *http.Request) {
	const op = "http.halls.List"
	log := h.log.With(slog.String("op", op))

	result, err := h.usecase.List(r.Context())
	if err != nil {
		log.Error("failed to get halls", "err", err)
		response.WriteError(w, "internal_error", "Failed to get halls", http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.ToHallListResponse(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

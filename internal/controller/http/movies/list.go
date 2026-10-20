package movies

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/dto"
)

func (h *MovieHandler) List(w http.ResponseWriter, r *http.Request) {
	const op = "http.movie.List"
	log := h.log.With(
		slog.String("op", op),
	)

	movies, err := h.usecase.List(r.Context())
	if err != nil {
		log.Error("failed to get movies", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := dto.FromModelList(movies)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

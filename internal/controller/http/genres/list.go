package genres

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/handler/dto"
)

func (h *GenreHandler) List(w http.ResponseWriter, r *http.Request) {
	const op = "http.genres.List"
	log := h.log.With(
		slog.String("op", op),
	)

	genres, err := h.usecase.List(r.Context())
	if err != nil {
		log.Error("failed to get genres", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := []dto.GenreResponse{}

	for _, m := range genres {
		resp = append(resp, dto.GenreResponse{
			ID:        m.ID(),
			TypeGenre: m.Name(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

package genres

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *GenreHandler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "http.genres.Get"
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

	genre, err := h.usecase.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrGenreNotFound) {
			log.Error("genre not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("failed to get genre", "err", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := dto.GenreResponse{}
	resp.FromModel(genre)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

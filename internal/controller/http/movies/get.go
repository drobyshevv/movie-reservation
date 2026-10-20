package movies

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *MovieHandler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "http.movie.Get"
	log := h.log.With(
		slog.String("op", op),
	)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Info("invalid path parameter", "param", "id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movie, err := h.usecase.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrMovieNotFound) {
			log.Info("movie not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("failed to get movie", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.MovieResponse{}
	resp.FromModel(movie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed encode response", "err", err)
	}
}

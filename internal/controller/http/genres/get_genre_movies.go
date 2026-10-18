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

func (h *GenreHandler) GetGenreMovies(w http.ResponseWriter, r *http.Request) {
	const op = "http.genre.GetGenreMovies"
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

	movies, err := h.usecase.GetGenreMovies(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrGenreNotFound) {
			log.Error("genre not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("failed to get genre movies", "err", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := []dto.GenreMovieResponse{}

	for _, m := range movies {
		movieResp := dto.GenreMovieResponse{}
		movieResp.FromModel(&m)
		resp = append(resp, movieResp)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

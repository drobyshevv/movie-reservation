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

func (h *MovieHandler) GetWithGenres(w http.ResponseWriter, r *http.Request) {
	const op = "http.movie.GetMovieGenres"

	log := h.log.With(
		slog.String("op", op),
	)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	withGenres, err := h.usecase.GetWithGenres(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrMovieNotFound) {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.MovieWithGenresResponse{}
	resp.FromModel(withGenres)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

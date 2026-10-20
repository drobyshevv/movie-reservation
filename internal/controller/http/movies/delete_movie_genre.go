package movies

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *MovieHandler) DeleteMovieGenre(w http.ResponseWriter, r *http.Request) {
	const op = "http.movie.DeleteMovieGenre"

	log := h.log.With(
		slog.String("op", op),
	)

	movieIDStr := r.PathValue("id")
	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "movie_id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	genreIDStr := r.PathValue("genre_id")
	genreID, err := strconv.ParseInt(genreIDStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "genre_id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usecase.DeleteMovieGenre(r.Context(), movieID, genreID)
	if err != nil {
		if errors.Is(err, errs.ErrMovieGenreNotFound) {
			http.Error(w, "movie_genre not found", http.StatusNotFound)
			return
		}
		log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

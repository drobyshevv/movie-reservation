package movies

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "http.movie.Delete"

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

	err = h.usecase.Delete(r.Context(), int64(id))
	if err != nil {
		if errors.Is(err, errs.ErrMovieNotFound) {
			log.Info("movie not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

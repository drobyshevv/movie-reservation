package genres

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *GenreHandler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "http.genres.Delete"

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

	err = h.usecase.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrGenreNotFound) {
			log.Error("genre not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("failed to delete genre", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

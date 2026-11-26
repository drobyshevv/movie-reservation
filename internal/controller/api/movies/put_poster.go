package movies

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/api/response"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *MovieHandler) PutPoster(w http.ResponseWriter, r *http.Request) {
	const op = "http.movies.PutPoster"

	log := h.log.With(slog.String("op", op))

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "id", "err", err)
		response.WriteError(w, response.CodeInvalidRequest, "Invalid ID", http.StatusBadRequest, "ID must be a valid integer")
		return
	}

	const maxBufferRequest = 10 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxBufferRequest)

	if err := r.ParseMultipartForm(maxBufferRequest); err != nil {
		log.Warn("file too large or invalid form", "err", err)
		response.WriteError(w, response.CodeInvalidRequest, "File too large or invalid form", http.StatusBadRequest, "")
		return
	}
	defer func() {
		if r.MultipartForm != nil {
			_ = r.MultipartForm.RemoveAll()
		}
	}()

	file, header, err := r.FormFile("poster")
	if err != nil {
		log.Error("failed to read file", "err", err)
		response.WriteError(w, response.CodeInvalidRequest, "Failed to read file", http.StatusBadRequest, "")
		return
	}

	err = h.movieUsecase.PutPoster(r.Context(), id, file, header.Size)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrMovieNotFound):
			log.Warn("movie not found", "err", err)
			response.WriteError(w, "not_found", "Movie not found", http.StatusNotFound, "")
		default:
			log.Error("failed to get movies", "err", err)
			response.WriteError(w, "internal_error", "Failed to get movies", http.StatusInternalServerError, "")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

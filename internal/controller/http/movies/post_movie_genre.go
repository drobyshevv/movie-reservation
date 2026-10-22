package movies

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/controller/http/response"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *MovieHandler) PostMovieGenre(w http.ResponseWriter, r *http.Request) {
	const op = "http.movies.PostMovieGenre"
	log := h.log.With(slog.String("op", op))

	movieIDStr := r.PathValue("id")
	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "movie_id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid movie ID", http.StatusBadRequest, "Movie ID must be a valid integer")
		return
	}

	genreIDStr := r.PathValue("genre_id")
	genreID, err := strconv.ParseInt(genreIDStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "genre_id", "err", err)
		response.WriteError(w, "invalid_request", "Invalid genre ID", http.StatusBadRequest, "Genre ID must be a valid integer")
		return
	}

	err = h.usecase.PostMovieGenre(r.Context(), movieID, genreID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrMovieNotFound):
			log.Warn("movie not found", "err", err)
			response.WriteError(w, "not_found", "Movie not found", http.StatusNotFound, err.Error())
		case errors.Is(err, errs.ErrGenreNotFound):
			log.Warn("genre not found", "err", err)
			response.WriteError(w, "not_found", "Genre not found", http.StatusNotFound, err.Error())
		case errors.Is(err, errs.ErrMovieGenreAlreadyExists):
			log.Warn("movie genre already exists", "err", err)
			response.WriteError(w, "conflict", "Movie genre already exists", http.StatusConflict, err.Error())
		default:
			log.Error("failed to add movie genre", "err", err)
			response.WriteError(w, "internal_error", "Failed to add movie genre", http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

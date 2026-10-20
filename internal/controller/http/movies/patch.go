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

func (h *MovieHandler) Patch(w http.ResponseWriter, r *http.Request) {
	const op = "http.movie.Patch"
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

	req := dto.UpdateMovieRequest{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Info("failed to convert query id to int", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Title == nil && req.Description == nil && req.Duration == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Info("invalid request", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input := req.ToModel()

	movie, err := h.usecase.Update(r.Context(), id, input)
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

	resp := dto.MovieResponse{}
	resp.FromModel(movie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

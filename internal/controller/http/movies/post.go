package movies

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/dto"
)

func (h *MovieHandler) Post(w http.ResponseWriter, r *http.Request) {
	const op = "http.movie.Post"

	log := h.log.With(
		slog.String("op", op),
	)

	req := dto.CreateMovieRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Info("failed to decode request body", "err", err)
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

	movie, err := h.usecase.Create(r.Context(), input)
	if err != nil {
		log.Error("failed to create movie", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.MovieResponse{}
	resp.FromModel(movie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed encode response", "err", err)
	}
}

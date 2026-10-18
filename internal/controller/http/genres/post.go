package genres

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/drobyshevv/movie-reservation/internal/dto"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

func (h *GenreHandler) Post(w http.ResponseWriter, r *http.Request) {
	const op = "http.genres.Post"
	log := h.log.With(
		slog.String("op", op),
	)

	var req dto.CreateGenreRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("failed to decode request body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Error("validation failed", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	genre, err := h.usecase.Create(r.Context(), req.TypeGenre)
	if err != nil {
		if errors.Is(err, errs.ErrGenreNotFound) {
			log.Error("genre not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, errs.ErrGenreAlreadyExists) {
			log.Error("genre conflict", "err", err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		log.Error("failed to create genre", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.GenreResponse{}
	resp.FromModel(genre)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

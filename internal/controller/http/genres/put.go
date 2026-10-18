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

func (h *GenreHandler) Put(w http.ResponseWriter, r *http.Request) {
	const op = "http.genres.Put"
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

	var req dto.UpdateGenreRequest
	err = json.NewDecoder(r.Body).Decode(&req)
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

	genre, err := h.usecase.Update(r.Context(), id, req.TypeGenre)
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
		log.Error("failed to update genre", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.GenreResponse{}
	resp.FromModel(genre)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

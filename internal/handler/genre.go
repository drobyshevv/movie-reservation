package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/handler/dto"
	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/go-playground/validator/v10"
)

type GenreHandler struct {
	log       *slog.Logger
	service   GenreService
	validator *validator.Validate
}

func NewGenreHandler(log *slog.Logger, service GenreService, validator *validator.Validate) *GenreHandler {
	return &GenreHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

type GenreService interface {
	GetGenres() ([]models.Genre, error)
	GetGenre(id int64) (*models.Genre, error)
	CreateGenre(typeGenre string) (*models.Genre, error)
	UpdateGenre(id int64, typeGenre string) (*models.Genre, error)
	DeleteGenre(id int64) error
	GetGenreMovies(id int64) ([]models.Movie, error)
}

func (h *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetGenres"
	log := h.log.With(
		slog.String("op", op),
	)

	genres, err := h.service.GetGenres()
	if err != nil {
		log.Error("failed to get genres", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := []dto.GenreResponse{}

	for _, m := range genres {
		resp = append(resp, dto.GenreResponse{
			ID:        m.ID,
			TypeGenre: m.TypeGenre,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

func (h *GenreHandler) GetGenre(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetGenre"
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

	genre, err := h.service.GetGenre(id)
	if err != nil {
		log.Error("failed to get genre", "err", err)
		w.WriteHeader(http.StatusNotFound)
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

func (h *GenreHandler) PostGenre(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateGenre"
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

	genre, err := h.service.CreateGenre(req.TypeGenre)
	if err != nil {
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

func (h *GenreHandler) PutGenre(w http.ResponseWriter, r *http.Request) {
	const op = "handler.PutGenre"
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

	genre, err := h.service.UpdateGenre(id, req.TypeGenre)
	if err != nil {
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

func (h *GenreHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteGenre"
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

	err = h.service.DeleteGenre(id)
	if err != nil {
		log.Error("failed to delete genre", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GenreHandler) GetGenreMovies(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetGenreMovies"
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

	movies, err := h.service.GetGenreMovies(id)
	if err != nil {
		log.Error("failed to get genre movies", "err", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := []dto.MovieResponse{}

	for _, m := range movies {
		movieResp := dto.MovieResponse{}
		movieResp.FromModel(&m)
		resp = append(resp, movieResp)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

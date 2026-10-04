package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/drobyshevv/movie-reservation/internal/handler/dto"
	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/drobyshevv/movie-reservation/internal/service"
)

type MovieHandler struct {
	log     *slog.Logger
	Service MovieService
}

func NewMovieHandler(service MovieService, log *slog.Logger) *MovieHandler {
	return &MovieHandler{
		log:     log,
		Service: service,
	}
}

type MovieService interface {
	GetMovies(ctx context.Context) ([]models.Movie, error)
	GetMovie(ctx context.Context, id int64) (*models.Movie, error)
	CreateMovie(ctx context.Context, params models.CreateMovieParams) (*models.Movie, error)
	UpdateMovie(ctx context.Context, id int64, params models.UpdateMovieParams) (*models.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
}

func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetMovies"
	h.log = h.log.With(
		slog.String("op", op),
	)

	movies, err := h.Service.GetMovies(r.Context())
	if err != nil {
		h.log.Error("failed to get movies", "err", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	resp := []dto.MovieResponse{}

	for _, m := range movies {
		resp = append(resp, dto.MovieResponse{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Duration:    int(m.Duration / time.Minute),
			Image:       m.Image,
			CreatedAt:   m.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.log.Error("failed to encode response", "err", err)
	}
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetMovie"
	h.log = h.log.With(
		slog.String("op", op),
	)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Info("invalid path parameter", "param", "id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movie, err := h.Service.GetMovie(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			h.log.Info("movie not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.log.Error("failed to get movie", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int(movie.Duration / time.Minute),
		Image:       movie.Image,
		CreatedAt:   movie.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.log.Error("failed encode response", "err", err)
	}
}

func (h *MovieHandler) PostMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.PostMovie"
	h.log = h.log.With(
		slog.String("op", op),
	)

	req := dto.CreateMovieRequest{}

	req.Title = r.FormValue("title")
	description := r.FormValue("description")

	if description == "" {
		req.Description = nil
	} else {
		req.Description = &description
	}

	durationStr := r.FormValue("duration")
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		h.log.Info("invalid form value parameter", "param", "duration", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req.Duration = duration
	file, _, err := r.FormFile("image")
	if err != nil {
		h.log.Info("failed to parse form file", "file", "image", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		h.log.Error("failed to read file", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Image = imageBytes

	params := models.CreateMovieParams{
		Title:       req.Title,
		Description: req.Description,
		Duration:    time.Duration(req.Duration) * time.Minute,
		Image:       req.Image,
	}

	movie, err := h.Service.CreateMovie(r.Context(), params)
	if err != nil {
		h.log.Error("failed to create movie", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int(movie.Duration / time.Second),
		Image:       movie.Image,
		CreatedAt:   movie.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.log.Error("failed encode response", "err", err)
	}
}

func (h *MovieHandler) PatchMovies(w http.ResponseWriter, r *http.Request) {
	const op = "handler.PatchMovies"
	h.log = h.log.With(
		slog.String("op", op),
	)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Info("invalid path parameter", "param", "id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := dto.UpdateMovieRequest{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.log.Info("failed to convert query id to int", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := models.UpdateMovieParams{
		Title:       req.Title,
		Description: req.Description,
		//Duration:    &duration,
		Image: req.Image,
	}

	if req.Duration != nil && *req.Duration != 0 {
		duration := time.Duration(*req.Duration) * time.Minute
		params.Duration = &duration
	}

	movie, err := h.Service.UpdateMovie(r.Context(), id, params)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			h.log.Info("movie not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int(movie.Duration / time.Second),
		Image:       movie.Image,
		CreatedAt:   movie.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.log.Info("failed to convert query id to int", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *MovieHandler) DeleteMovies(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteMovies"
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Error("invalid path parameter", "param", "id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteMovie(r.Context(), int64(id))
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			h.log.Info("movie not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

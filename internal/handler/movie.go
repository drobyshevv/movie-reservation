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
	"github.com/go-playground/validator/v10"
)

type MovieHandler struct {
	log       *slog.Logger
	service   MovieService
	validator *validator.Validate
}

func NewMovieHandler(log *slog.Logger, service MovieService, validator *validator.Validate) *MovieHandler {
	return &MovieHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

type MovieService interface {
	GetMovies(ctx context.Context) ([]models.Movie, error)
	GetMovie(ctx context.Context, id int64) (*models.Movie, error)
	CreateMovie(ctx context.Context, params models.CreateMovieParams) (*models.Movie, error)
	UpdateMovie(ctx context.Context, id int64, params models.UpdateMovieParams) (*models.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
	GetImage(ctx context.Context, id int64) ([]byte, error)
	PutImage(ctx context.Context, id int64, image []byte) error
	DeleteImage(ctx context.Context, id int64) error
	GetMovieGenres(ctx context.Context, id int64) ([]models.Genre, error)
	PostMovieGenre(ctx context.Context, movieID int64, genreID int64) error
	DeleteMovieGenre(ctx context.Context, movieID int64, genreID int64) error
}

func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetMovies"
	log := h.log.With(
		slog.String("op", op),
	)

	movies, err := h.service.GetMovies(r.Context())
	if err != nil {
		log.Error("failed to get movies", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := []dto.MovieResponse{}

	for _, m := range movies {
		resp = append(resp, dto.MovieResponse{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Duration:    int(m.Duration / time.Minute),
			CreatedAt:   m.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetMovie"
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

	movie, err := h.service.GetMovie(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			log.Info("movie not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("failed to get movie", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int(movie.Duration / time.Minute),
		CreatedAt:   movie.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed encode response", "err", err)
	}
}

func (h *MovieHandler) PostMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.PostMovie"

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

	params := models.CreateMovieParams{
		Title:       req.Title,
		Description: req.Description,
		Duration:    time.Duration(req.Duration) * time.Minute,
	}

	movie, err := h.service.CreateMovie(r.Context(), params)
	if err != nil {
		log.Error("failed to create movie", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int(movie.Duration / time.Minute),
		CreatedAt:   movie.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed encode response", "err", err)
	}
}

func (h *MovieHandler) PatchMovies(w http.ResponseWriter, r *http.Request) {
	const op = "handler.PatchMovies"
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

	params := models.UpdateMovieParams{
		Title:       req.Title,
		Description: req.Description,
	}

	if req.Duration != nil {
		duration := time.Duration(*req.Duration) * time.Minute
		params.Duration = &duration
	}

	movie, err := h.service.UpdateMovie(r.Context(), id, params)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			log.Info("movie not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := dto.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int(movie.Duration / time.Minute),
		CreatedAt:   movie.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", "err", err)
	}
}

func (h *MovieHandler) DeleteMovies(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteMovies"

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

	err = h.service.DeleteMovie(r.Context(), int64(id))
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
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

func (h *MovieHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetImage"

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

	image, err := h.service.GetImage(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			log.Info("movie not found", "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(image)

	w.Header().Set("Content-Type", contentType)
	_, err = w.Write(image)
	if err != nil {
		log.Error("failed to write image", "err", err)
	}
}

func (h *MovieHandler) PutImage(w http.ResponseWriter, r *http.Request) {
	const op = "handler.PutImage"

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

	file, header, err := r.FormFile("image")
	if err != nil {
		log.Info("failed to parse form file", "file", "image", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > 5*1024*1024 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		log.Error("failed to read file", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(imageBytes)

	switch contentType {
	case "image/jpeg", "image/png", "image/webp":
		// OK
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.service.PutImage(r.Context(), id, imageBytes)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
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

func (h *MovieHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteImage"

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

	err = h.service.DeleteImage(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
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

func (h *MovieHandler) GetMovieGenres(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetMovieGenres"

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

	genres, err := h.service.GetMovieGenres(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
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

func (h *MovieHandler) PostMovieGenre(w http.ResponseWriter, r *http.Request) {
	const op = "handler.PostMovieGenre"

	log := h.log.With(
		slog.String("op", op),
	)

	movieIDStr := r.PathValue("id")
	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "movie_id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	genreIDStr := r.PathValue("genre_id")
	genreID, err := strconv.ParseInt(genreIDStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "genre_id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.service.PostMovieGenre(r.Context(), movieID, genreID)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrGenreNotFound) {
			http.Error(w, "genre not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrMovieGenreAlreadyExists) {
			http.Error(w, "movie_genre already exists", http.StatusConflict)
			return
		}
		log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *MovieHandler) DeleteMovieGenre(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteMovieGenre"

	log := h.log.With(
		slog.String("op", op),
	)

	movieIDStr := r.PathValue("id")
	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "movie_id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	genreIDStr := r.PathValue("genre_id")
	genreID, err := strconv.ParseInt(genreIDStr, 10, 64)
	if err != nil {
		log.Error("invalid path parameter", "param", "genre_id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.service.DeleteMovieGenre(r.Context(), movieID, genreID)
	if err != nil {
		if errors.Is(err, service.ErrMovieGenreNotFound) {
			http.Error(w, "movie_genre not found", http.StatusNotFound)
			return
		}
		log.Error("internal server error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

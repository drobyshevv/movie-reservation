package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/drobyshevv/movie-reservation/internal/repository"
)

type MovieService struct {
	Repository MovieRepository
}

func NewMovieService(repository MovieRepository) *MovieService {
	return &MovieService{
		Repository: repository,
	}
}

type MovieRepository interface {
	GetAll(ctx context.Context) ([]models.Movie, error)
	Get(ctx context.Context, id int64) (*models.Movie, error)
	Create(ctx context.Context, params models.CreateMovieParams) (*models.Movie, error)
	Update(ctx context.Context, id int64, params models.UpdateMovieParams) (*models.Movie, error)
	Delete(ctx context.Context, id int64) error
}

func (s *MovieService) GetMovies(ctx context.Context) ([]models.Movie, error) {
	const op = "service.GetMovies"

	movies, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

func (s *MovieService) GetMovie(ctx context.Context, id int64) (*models.Movie, error) {
	const op = "service.GetMovie"

	movie, err := s.Repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrMovieNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movie, nil
}

func (s *MovieService) CreateMovie(ctx context.Context, params models.CreateMovieParams) (*models.Movie, error) {
	const op = "service.CreateMovie"

	movie, err := s.Repository.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movie, nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, id int64, params models.UpdateMovieParams) (*models.Movie, error) {
	const op = "service.UpdateMovie"

	movie, err := s.Repository.Update(ctx, id, params)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrMovieNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movie, nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, id int64) error {
	const op = "service.DeleteMovie"

	err := s.Repository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrMovieNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

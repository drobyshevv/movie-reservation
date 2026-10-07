package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/drobyshevv/movie-reservation/internal/repository"
)

type MovieService struct {
	repository MovieRepository
}

func NewMovieService(repository MovieRepository) *MovieService {
	return &MovieService{
		repository: repository,
	}
}

type MovieRepository interface {
	GetAll(ctx context.Context) ([]models.Movie, error)
	Get(ctx context.Context, id int64) (*models.Movie, error)
	Create(ctx context.Context, params models.CreateMovieParams) (*models.Movie, error)
	Update(ctx context.Context, id int64, params models.UpdateMovieParams) (*models.Movie, error)
	Delete(ctx context.Context, id int64) error
	GetImage(ctx context.Context, id int64) ([]byte, error)
	PutImage(ctx context.Context, id int64, image []byte) error
	DeleteImage(ctx context.Context, id int64) error
	GetMovieGenres(ctx context.Context, id int64) ([]models.Genre, error)
	PostMovieGenre(ctx context.Context, movieID int64, genreID int64) error
	DeleteMovieGenre(ctx context.Context, movieID int64, genreID int64) error
}

func (s *MovieService) GetMovies(ctx context.Context) ([]models.Movie, error) {
	const op = "service.GetMovies"

	movies, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

func (s *MovieService) GetMovie(ctx context.Context, id int64) (*models.Movie, error) {
	const op = "service.GetMovie"

	movie, err := s.repository.Get(ctx, id)
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

	movie, err := s.repository.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movie, nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, id int64, params models.UpdateMovieParams) (*models.Movie, error) {
	const op = "service.UpdateMovie"

	movie, err := s.repository.Update(ctx, id, params)
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

	err := s.repository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrMovieNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *MovieService) GetImage(ctx context.Context, id int64) ([]byte, error) {
	const op = "service.GetImage"

	image, err := s.repository.GetImage(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrMovieNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return image, nil
}

func (s *MovieService) PutImage(ctx context.Context, id int64, image []byte) error {
	const op = "service.PutImage"

	err := s.repository.PutImage(ctx, id, image)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrMovieNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *MovieService) DeleteImage(ctx context.Context, id int64) error {
	const op = "service.DeleteImage"

	err := s.repository.DeleteImage(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrMovieNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *MovieService) GetMovieGenres(ctx context.Context, id int64) ([]models.Genre, error) {
	const op = "service.GetMovieGenres"

	genres, err := s.repository.GetMovieGenres(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrMovieNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return genres, nil

}

func (s *MovieService) PostMovieGenre(ctx context.Context, movieID int64, genreID int64) error {
	const op = "service.PostMovieGenre"

	err := s.repository.PostMovieGenre(ctx, movieID, genreID)
	if err != nil {
		if errors.Is(err, repository.ErrMovieNotFound) {
			return fmt.Errorf("%s: %w", op, ErrMovieNotFound)
		}
		if errors.Is(err, repository.ErrGenreNotFound) {
			return fmt.Errorf("%s: %w", op, ErrGenreNotFound)
		}
		if errors.Is(err, repository.ErrConflict) {
			return fmt.Errorf("%s: %w", op, ErrMovieGenreAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (s *MovieService) DeleteMovieGenre(ctx context.Context, movieID int64, genreID int64) error {
	const op = "service.DeleteMovieGenre"

	err := s.repository.DeleteMovieGenre(ctx, movieID, genreID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrMovieGenreNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

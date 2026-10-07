package service

import (
	"errors"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/drobyshevv/movie-reservation/internal/repository"
)

type GenreService struct {
	repo GenreRepository
}

func NewGenreService(repo GenreRepository) *GenreService {
	return &GenreService{
		repo: repo,
	}
}

type GenreRepository interface {
	GetGenres() ([]models.Genre, error)
	GetGenre(id int64) (*models.Genre, error)
	CreateGenre(typeGenre string) (*models.Genre, error)
	UpdateGenre(id int64, typeGenre string) (*models.Genre, error)
	DeleteGenre(id int64) error
	GetGenreMovies(id int64) ([]models.Movie, error)
}

func (s *GenreService) GetGenres() ([]models.Genre, error) {
	const op = "service.GetGenres"

	genres, err := s.repo.GetGenres()
	if err != nil {
		return nil, err
	}

	return genres, nil
}

func (s *GenreService) GetGenre(id int64) (*models.Genre, error) {
	const op = "service.GetGenre"

	genre, err := s.repo.GetGenre(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrGenreNotFound)
		}
		return nil, err
	}

	return genre, nil
}

func (s *GenreService) CreateGenre(typeGenre string) (*models.Genre, error) {
	const op = "service.CreateGenre"

	genre, err := s.repo.CreateGenre(typeGenre)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, fmt.Errorf("%s: %w", op, ErrGenreAlreadyExists)
		}
		return nil, err
	}

	return genre, nil
}

func (s *GenreService) UpdateGenre(id int64, typeGenre string) (*models.Genre, error) {
	const op = "service.UpdateGenre"

	genre, err := s.repo.UpdateGenre(id, typeGenre)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrGenreNotFound)
		}
		return nil, err
	}

	return genre, nil
}

func (s *GenreService) DeleteGenre(id int64) error {
	const op = "service.DeleteGenre"

	err := s.repo.DeleteGenre(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrGenreNotFound)
		}
		return err
	}

	return nil
}

func (s *GenreService) GetGenreMovies(id int64) ([]models.Movie, error) {
	const op = "service.GetGenreMovies"

	movies, err := s.repo.GetGenreMovies(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrGenreNotFound)
		}
		return nil, err
	}

	return movies, nil
}

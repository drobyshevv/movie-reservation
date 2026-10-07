package service

import (
	"errors"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/drobyshevv/movie-reservation/internal/repository"
)

type HallService struct {
	repo HallRepository
}

func NewHallService(repo HallRepository) *HallService {
	return &HallService{
		repo: repo,
	}
}

type HallRepository interface {
	GetHalls() ([]models.Hall, error)
	GetHall(id int64) (*models.Hall, error)
	CreateHall(params models.CreateHallParams) (*models.Hall, error)
	UpdateHall(id int64, params models.UpdateHallParams) (*models.Hall, error)
	DeleteHall(id int64) error
}

func (s *HallService) GetHalls() ([]models.Hall, error) {
	const op = "service.GetHalls"

	halls, err := s.repo.GetHalls()
	if err != nil {
		return nil, err
	}

	return halls, nil
}

func (s *HallService) GetHall(id int64) (*models.Hall, error) {
	const op = "service.GetHall"

	hall, err := s.repo.GetHall(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrHallNotFound)
		}
		return nil, err
	}

	return hall, nil
}

func (s *HallService) CreateHall(params models.CreateHallParams) (*models.Hall, error) {
	const op = "service.CreateHall"

	hall, err := s.repo.CreateHall(params)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, fmt.Errorf("%s: %w", op, ErrHallAlreadyExists)
		}
		return nil, err
	}

	return hall, nil
}

func (s *HallService) UpdateHall(id int64, params models.UpdateHallParams) (*models.Hall, error) {
	const op = "service.UpdateHall"

	hall, err := s.repo.UpdateHall(id, params)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrHallNotFound)
		}
		return nil, err
	}

	return hall, nil
}

func (s *HallService) DeleteHall(id int64) error {
	const op = "service.DeleteHall"

	err := s.repo.DeleteHall(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrHallNotFound)
		}
		return err
	}

	return nil
}

package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/drobyshevv/movie-reservation/internal/repository"
)

type HallService struct {
	repository HallRepository
}

func NewHallService(repository HallRepository) *HallService {
	return &HallService{
		repository: repository,
	}
}

type HallRepository interface {
	GetHalls(ctx context.Context) ([]models.Hall, error)
	GetHall(ctx context.Context, id int64) (*models.Hall, error)
	CreateHall(ctx context.Context, params models.CreateHallParams) (*models.Hall, error)
	UpdateHall(ctx context.Context, id int64, params models.UpdateHallParams) (*models.Hall, error)
	DeleteHall(ctx context.Context, id int64) error
	GetHallSeats(ctx context.Context, id int64) ([]models.Seat, error)
	CreateHallSeats(ctx context.Context, id int64, params []models.CreateHallSeatsParams) ([]models.Seat, error)
	DeleteHallSeats(ctx context.Context, id int64) error
}

func (s *HallService) GetHalls(ctx context.Context) ([]models.Hall, error) {
	const op = "service.GetHalls"

	halls, err := s.repository.GetHalls(ctx)
	if err != nil {
		return nil, err
	}

	return halls, nil
}

func (s *HallService) GetHall(ctx context.Context, id int64) (*models.Hall, error) {
	const op = "service.GetHall"

	hall, err := s.repository.GetHall(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrHallNotFound)
		}
		return nil, err
	}

	return hall, nil
}

func (s *HallService) CreateHall(ctx context.Context, params models.CreateHallParams) (*models.Hall, error) {
	const op = "service.CreateHall"

	hall, err := s.repository.CreateHall(ctx, params)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, fmt.Errorf("%s: %w", op, ErrHallAlreadyExists)
		}
		return nil, err
	}

	return hall, nil
}

func (s *HallService) UpdateHall(ctx context.Context, id int64, params models.UpdateHallParams) (*models.Hall, error) {
	const op = "service.UpdateHall"

	hall, err := s.repository.UpdateHall(ctx, id, params)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrHallNotFound)
		}
		return nil, err
	}

	return hall, nil
}

func (s *HallService) DeleteHall(ctx context.Context, id int64) error {
	const op = "service.DeleteHall"

	err := s.repository.DeleteHall(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrHallNotFound)
		}
		return err
	}

	return nil
}

func (s *HallService) GetHallSeats(ctx context.Context, id int64) ([]models.Seat, error) {
	const op = "service.GetHallSeats"

	seats, err := s.repository.GetHallSeats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return seats, nil
}

func (s *HallService) CreateHallSeats(ctx context.Context, id int64, params []models.CreateHallSeatsParams) ([]models.Seat, error) {
	const op = "service.CreateHallSeats"

	seats, err := s.repository.CreateHallSeats(ctx, id, params)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrHallNotFound)
		}
		if errors.Is(err, repository.ErrConflict) {
			return nil, fmt.Errorf("%s: %w", op, ErrSeatAlreadyExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return seats, nil
}

func (s *HallService) DeleteHallSeats(ctx context.Context, id int64) error {
	const op = "servie.DeleteHallSeats"

	err := s.repository.DeleteHallSeats(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrHallNotFound) {
			return fmt.Errorf("%s: %w", op, ErrHallNotFound)
		}
		if errors.Is(err, repository.ErrSeatNotFound) {
			return fmt.Errorf("%s: %w", op, ErrSeatNotFound)
		}
		if errors.Is(err, repository.ErrSessionExists) {
			return fmt.Errorf("%s: %w", op, ErrHallHasSessions)
		}
	}
	return nil
}

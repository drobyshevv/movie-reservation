package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/repository"
)

type SeatService struct {
	repository SeatRepository
}

func NewSeatService(repository SeatRepository) *SeatService {
	return &SeatService{
		repository: repository,
	}
}

type SeatRepository interface {
	DeleteSeat(ctx context.Context, id int64) error
}

func (s *SeatService) DeleteSeat(ctx context.Context, id int64) error {
	op := "service.DeleteSeat"
	err := s.repository.DeleteSeat(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, ErrSeatNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

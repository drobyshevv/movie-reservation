package bookings

import (
	"context"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	"github.com/google/uuid"
)

type BookingUseCase struct {
	repository BookingRepository
}

func NewBookingUseCase(repository BookingRepository) *BookingUseCase {
	return &BookingUseCase{
		repository: repository,
	}
}

type BookingRepository interface {
	Get(ctx context.Context, id int64) (*domain.Booking, error)
	Create(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	ListBySession(ctx context.Context, sessionID int64) ([]SeatStatusOutput, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]UserBookingOutput, error)
	UpdateStatus(ctx context.Context, id int64, status domain.BookingStatus, updateFn func(booking *domain.Booking, status domain.BookingStatus) error) error
}

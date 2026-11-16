package bookings

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *BookingUseCase) Create(ctx context.Context, input CreateBookingInput) (*domain.Booking, error) {
	const op = "usecase.bookings.Create"

	booking, err := domain.NewBooking(input.SessionID, input.UserID, input.SeatID)
	if err != nil {
		return nil, fmt.Errorf("%s init domain model: %w", op, err)
	}

	created, err := u.repository.Create(ctx, booking)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return created, nil
}

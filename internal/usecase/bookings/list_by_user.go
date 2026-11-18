package bookings

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (u *BookingUseCase) ListByUser(ctx context.Context, userID uuid.UUID) ([]UserBookingOutput, error) {
	const op = "usecase.bookings.ListByUser"
	bookings, err := u.repository.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return bookings, nil
}

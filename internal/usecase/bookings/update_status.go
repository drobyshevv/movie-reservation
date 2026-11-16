package bookings

import (
	"context"
	"fmt"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

func (u *BookingUseCase) UpdateStatus(ctx context.Context, id int64, status domain.BookingStatus) error {
	const op = "usecase.bookings.UpdateStatus"

	err := u.repository.UpdateStatus(ctx, func(id int64, booking *domain.Booking, status domain.BookingStatus) error {
		if !status.IsValid() {
			return fmt.Errorf("%s: validate status: %w", op, domain.ErrInvalidBookingStatus)
		}

		if err := booking.ChangeStatus(status); err != nil {
			return fmt.Errorf("%s: change status: %w", op, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

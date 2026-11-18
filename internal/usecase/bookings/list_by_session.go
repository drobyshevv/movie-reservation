package bookings

import (
	"context"
	"fmt"
)

func (u *BookingUseCase) ListBySession(ctx context.Context, sessionID int64) ([]SeatStatusOutput, error) {
	const op = "usecase.bookings.ListBySession"
	seats, err := u.repository.ListBySession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return seats, nil
}

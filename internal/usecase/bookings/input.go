package bookings

import "github.com/google/uuid"

type CreateBookingInput struct {
	SessionID int64
	UserID    uuid.UUID
	SeatID    int64
}

package bookings

import (
	"time"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

type SeatStatusOutput struct {
	SeatID int64
	Status domain.BookingStatus
}

type UserBookingOutput struct {
	ID            int64
	MovieTitle    string
	HallName      string
	StartAt       time.Time
	SeatNumber    string
	BookingStatus domain.BookingStatus
	SessionStatus domain.SessionStatus
}

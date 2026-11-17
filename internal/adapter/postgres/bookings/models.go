package bookings

import (
	"time"

	"github.com/google/uuid"
)

type bookingRow struct {
	id        int64
	sessionID int64
	userID    uuid.UUID
	seatID    int64
	status    string
	createdAt time.Time
}

type seatStatusRow struct {
	seatID int64
	status string
}

type userBookingRow struct {
	id         int64
	movieTitle string
	hallName   string
	startAt    time.Time
	seatNumber string
	status     string
}

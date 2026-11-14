package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidSessionID        = errors.New("session ID must be positive")
	ErrInvalidUserID           = errors.New("user ID must be positive")
	ErrInvalidBookingIDs       = errors.New("invalid booking IDs")
	ErrBookingNotPending       = errors.New("only pending bookings can be confirmed")
	ErrBookingNotConfirmed     = errors.New("only confirmed bookings can be completed")
	ErrBookingAlreadyCancelled = errors.New("booking is already cancelled")
	ErrBookingAlreadyCompleted = errors.New("booking is already completed")
	ErrInvalidBookingStatus    = errors.New("invalid booking status")
)

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCancelled BookingStatus = "cancelled"
	BookingCompleted BookingStatus = "completed"
)

type Booking struct {
	id        int64
	sessionID int64
	userID    uuid.UUID
	seatID    int64
	status    BookingStatus
	createdAt time.Time
}

func (s BookingStatus) IsValid() bool {
	switch s {
	case BookingPending, BookingConfirmed, BookingCancelled, BookingCompleted:
		return true
	default:
		return false
	}
}

func NewBooking(sessionID int64, userID uuid.UUID, seatID int64) (*Booking, error) {
	if sessionID <= 0 {
		return nil, ErrInvalidSessionID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if seatID <= 0 {
		return nil, ErrInvalidSeatID
	}
	return &Booking{
		sessionID: sessionID,
		userID:    userID,
		seatID:    seatID,
		status:    BookingPending,
	}, nil
}

func RestoreBooking(
	id, sessionID, seatID int64,
	userID uuid.UUID,
	status BookingStatus,
	createdAt time.Time,
) (*Booking, error) {
	if id <= 0 || sessionID <= 0 || seatID <= 0 {
		return nil, ErrInvalidSessionID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	return &Booking{
		id:        id,
		sessionID: sessionID,
		userID:    userID,
		seatID:    seatID,
		status:    status,
		createdAt: createdAt,
	}, nil
}

func (b *Booking) ChangeStatus(status BookingStatus) error {
	switch status {
	case BookingConfirmed:
		if b.status != BookingPending {
			return ErrBookingNotPending
		}
	case BookingCompleted:
		if b.status != BookingConfirmed {
			return ErrBookingNotConfirmed
		}
	case BookingCancelled:
		if b.status == BookingCancelled {
			return ErrBookingAlreadyCancelled
		}
		if b.status == BookingCompleted {
			return ErrBookingAlreadyCompleted
		}
	case BookingPending:
		return ErrInvalidBookingStatus
	default:
		return ErrInvalidBookingStatus
	}

	b.status = status
	return nil
}

func (b *Booking) ID() int64             { return b.id }
func (b *Booking) SessionID() int64      { return b.sessionID }
func (b *Booking) UserID() uuid.UUID     { return b.userID }
func (b *Booking) SeatID() int64         { return b.seatID }
func (b *Booking) Status() BookingStatus { return b.status }
func (b *Booking) CreatedAt() time.Time  { return b.createdAt }

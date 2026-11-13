package dto

import (
	"time"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	"github.com/google/uuid"

	bookinguc "github.com/drobyshevv/movie-reservation/internal/usecase/bookings"
)

type CreateBookingRequest struct {
	SessionID int64  `json:"session_id" validate:"required,gt=0"`
	UserID    string `json:"user_id" validate:"required,uuid"`
	SeatID    int64  `json:"seat_id" validate:"required,gt=0"`
}

type BookingResponse struct {
	ID        int64     `json:"id"`
	SessionID int64     `json:"session_id"`
	UserID    string    `json:"user_id"`
	SeatID    int64     `json:"seat_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type SeatStatusResponse struct {
	SeatID int64  `json:"seat_id"`
	Status string `json:"status"`
}

type UserBookingResponse struct {
	ID         int64     `json:"id"`
	MovieTitle string    `json:"movie_title"`
	HallName   string    `json:"hall_name"`
	StartAt    time.Time `json:"start_at"`
	SeatNumber string    `json:"seat_number"`
	Status     string    `json:"status"`
}

func (b *CreateBookingRequest) ToCreateInput() (*bookinguc.CreateBookingInput, error) {
	userID, err := uuid.Parse(b.UserID)
	if err != nil {
		return nil, err
	}

	return &bookinguc.CreateBookingInput{
		SessionID: b.SessionID,
		UserID:    userID,
		SeatID:    b.SeatID,
	}, nil
}

func ToBookingResponse(b *domain.Booking) BookingResponse {
	if b == nil {
		return BookingResponse{}
	}
	return BookingResponse{
		ID:        b.ID(),
		SessionID: b.SessionID(),
		UserID:    b.UserID().String(),
		SeatID:    b.SeatID(),
		Status:    string(b.Status()),
		CreatedAt: b.CreatedAt(),
	}
}

func ToSeatStatusesResponse(st []bookinguc.SeatStatusOutput) []SeatStatusResponse {
	result := make([]SeatStatusResponse, len(st))
	for i, s := range st {
		result[i] = SeatStatusResponse{
			SeatID: s.SeatID,
			Status: string(s.Status),
		}
	}
	return result
}

func ToUserBooking(id int64, movieTitle, hallName, seatNumber string, startAt time.Time, status domain.BookingStatus) UserBookingResponse {
	return UserBookingResponse{
		ID:         id,
		MovieTitle: movieTitle,
		HallName:   hallName,
		StartAt:    startAt,
		SeatNumber: seatNumber,
		Status:     string(status),
	}
}

package dto

import (
	"errors"
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

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required"`
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
	ID            int64     `json:"id"`
	MovieTitle    string    `json:"movie_title"`
	HallName      string    `json:"hall_name"`
	StartAt       time.Time `json:"start_at"`
	SeatNumber    string    `json:"seat_number"`
	BookingStatus string    `json:"booking_status"`
	SessionStatus string    `json:"session_status"`
}

func (b *CreateBookingRequest) ToInput() (*bookinguc.CreateBookingInput, error) {
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

func (b *UpdateStatusRequest) ToInput() (domain.BookingStatus, error) {
	status := domain.BookingStatus(b.Status)
	if !status.IsValid() {
		return "", errors.New("invalid booking status")
	}
	return status, nil
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

func ToSeatStatusesResponse(items []bookinguc.SeatStatusOutput) []SeatStatusResponse {
	result := make([]SeatStatusResponse, len(items))
	for i, item := range items {
		result[i] = SeatStatusResponse{
			SeatID: item.SeatID,
			Status: string(item.Status),
		}
	}
	return result
}

func ToUserBookings(items []bookinguc.UserBookingOutput) []UserBookingResponse {
	result := make([]UserBookingResponse, len(items))
	for i, item := range items {
		result[i] = UserBookingResponse{
			ID:            item.ID,
			MovieTitle:    item.MovieTitle,
			HallName:      item.HallName,
			StartAt:       item.StartAt,
			SeatNumber:    item.SeatNumber,
			BookingStatus: string(item.BookingStatus),
			SessionStatus: string(item.SessionStatus),
		}
	}
	return result
}

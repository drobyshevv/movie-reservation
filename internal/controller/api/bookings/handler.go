package bookings

import (
	"context"
	"log/slog"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	bookinguc "github.com/drobyshevv/movie-reservation/internal/usecase/bookings"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BookingHandler struct {
	log       *slog.Logger
	usecase   BookingUseCase
	validator *validator.Validate
}

func NewBookingHandler(log *slog.Logger, usecase BookingUseCase, validator *validator.Validate) *BookingHandler {
	return &BookingHandler{
		log:       log,
		usecase:   usecase,
		validator: validator,
	}
}

type BookingUseCase interface {
	Get(ctx context.Context, id int64) (*domain.Booking, error)
	Create(ctx context.Context, input *bookinguc.CreateBookingInput) (*domain.Booking, error)
	ListBySession(ctx context.Context, sessionID int64) ([]bookinguc.SeatStatusOutput, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]bookinguc.UserBookingOutput, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	Delete(ctx context.Context, id int64) error
}

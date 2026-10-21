package halls

import (
	"context"
	"log/slog"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	"github.com/go-playground/validator/v10"

	halluc "github.com/drobyshevv/movie-reservation/internal/usecase/halls"
)

type HallHandler struct {
	log       *slog.Logger
	usecase   HallUseCase
	validator *validator.Validate
}

func NewHallHandler(log *slog.Logger, usecase HallUseCase, validator *validator.Validate) *HallHandler {
	return &HallHandler{
		log:       log,
		usecase:   usecase,
		validator: validator,
	}
}

type HallUseCase interface {
	List(ctx context.Context) ([]domain.Hall, error)
	Get(ctx context.Context, id int64) (*domain.Hall, error)
	GetWithSeats(ctx context.Context, id int64) (halluc.HallWithSeats, error)
	GetSummary(ctx context.Context, id int64) (halluc.HallSummary, error)
	CreateWithSeats(ctx context.Context, input halluc.CreateHallWithSeatsInput) (halluc.HallWithSeats, error)
	Update(ctx context.Context, id int64, name string) (*domain.Hall, error)
	Delete(ctx context.Context, id int64) error
	ReplaceSeats(ctx context.Context, id int64, input []halluc.HallSeatsInput) ([]domain.Seat, error)
	DeleteSeats(ctx context.Context, id int64) error
}

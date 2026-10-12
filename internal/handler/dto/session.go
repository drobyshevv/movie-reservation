package dto

import (
	"time"

	"github.com/drobyshevv/movie-reservation/internal/models"
)

type CreateSessionRequest struct {
	HallID  int64      `json:"hall_id" validate:"required,gt=0"`
	MovieID int64      `json:"movie_id" validate:"required,gt=0"`
	StartAt time.Time  `json:"start_at" validate:"required"`
	EndAt   *time.Time `json:"end_at" validate:"omitempty"`
}

type UpdateSessionRequest struct {
	HallID  int64      `json:"hall_id" validate:"required,gt=0"`
	MovieID int64      `json:"movie_id" validate:"required,gt=0"`
	StartAt time.Time  `json:"start_at" validate:"required"`
	EndAt   *time.Time `json:"end_at" validate:"omitempty"`
}

type UpdateSessionStatusRequest struct {
	Status models.StatusType `json:"status" validate:"required,oneof=scheduled active cancelled"`
}

type SessionResponse struct {
	ID        int64             `json:"id"`
	HallID    int64             `json:"hall_id"`
	MovieID   int64             `json:"movie_id"`
	Status    models.StatusType `json:"status"`
	StartAt   time.Time         `json:"start_at"`
	EndAt     time.Time         `json:"end_at"`
	CreatedAt time.Time         `json:"created_at"`
}

type SessionWithDetailsResponse struct {
	ID        int64             `json:"id"`
	Status    models.StatusType `json:"status"`
	StartAt   time.Time         `json:"start_at"`
	EndAt     time.Time         `json:"end_at"`
	CreatedAt time.Time         `json:"created_at"`
	Hall      HallShort         `json:"hall"`
	Movie     MovieShort        `json:"movie"`
}

type HallShort struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type MovieShort struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

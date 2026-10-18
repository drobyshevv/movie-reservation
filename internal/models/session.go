package models

import (
	"time"

	"github.com/drobyshevv/movie-reservation/internal/domain"
)

type CreateSessionInput struct {
	HallID  int64
	MovieID int64
	StartAt time.Time
	EndAt   *time.Time
}

type UpdateSessionInput struct {
	Status  *domain.StatusType
	StartAt *time.Time
	EndAt   *time.Time
}

type SessionWithDetails struct {
	ID        int64
	Status    domain.StatusType
	StartAt   time.Time
	EndAt     time.Time
	CreatedAt time.Time
	Hall      HallShort
	Movie     MovieShort
}

type HallShort struct {
	ID   int64
	Name string
}

type MovieShort struct {
	ID    int64
	Title string
}

package sessions

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

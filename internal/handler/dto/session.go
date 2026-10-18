package dto

import (
	"time"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	"github.com/drobyshevv/movie-reservation/internal/models"
)

type CreateSessionRequest struct {
	HallID  int64      `json:"hall_id" validate:"required,gt=0"`
	MovieID int64      `json:"movie_id" validate:"required,gt=0"`
	StartAt time.Time  `json:"start_at" validate:"required"`
	EndAt   *time.Time `json:"end_at" validate:"omitempty"`
}

type UpdateSessionRequest struct {
	Status  *domain.StatusType `json:"status" validate:"omitempty,oneof=scheduled active cancelled"`
	StartAt *time.Time         `json:"start_at"`
	EndAt   *time.Time         `json:"end_at"`
}

type SessionResponse struct {
	ID        int64             `json:"id"`
	HallID    int64             `json:"hall_id"`
	MovieID   int64             `json:"movie_id"`
	Status    domain.StatusType `json:"status"`
	StartAt   time.Time         `json:"start_at"`
	EndAt     time.Time         `json:"end_at"`
	CreatedAt time.Time         `json:"created_at"`
}

type SessionWithDetailsResponse struct {
	ID        int64             `json:"id"`
	Status    domain.StatusType `json:"status"`
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

func (s *SessionWithDetailsResponse) FromModel(session *models.SessionWithDetails) {
	s.ID = session.ID
	s.Status = session.Status
	s.StartAt = session.StartAt
	s.EndAt = session.EndAt
	s.CreatedAt = session.CreatedAt
	s.Hall = HallShort(session.Hall)
	s.Movie = MovieShort(session.Movie)
}

func SessionWithDetailsListFromModels(
	sessions []models.SessionWithDetails,
) []SessionWithDetailsResponse {
	resp := make([]SessionWithDetailsResponse, 0, len(sessions))

	for i := range sessions {
		var item SessionWithDetailsResponse
		item.FromModel(&sessions[i])

		resp = append(resp, item)
	}

	return resp
}

func (s *CreateSessionRequest) ToModel() models.CreateSessionInput {
	return models.CreateSessionInput{
		HallID:  s.HallID,
		MovieID: s.MovieID,
		StartAt: s.StartAt,
		EndAt:   s.EndAt,
	}
}

func (s *SessionResponse) FromModel(session *domain.Session) {
	s.ID = session.ID()
	s.HallID = session.HallID()
	s.MovieID = session.MovieID()
	s.Status = session.Status()
	s.StartAt = session.Period().Start()
	s.EndAt = session.Period().End()
	s.CreatedAt = session.CreatedAt()
}

func (s *UpdateSessionRequest) ToModel() models.UpdateSessionInput {
	return models.UpdateSessionInput{
		Status:  s.Status,
		StartAt: s.StartAt,
		EndAt:   s.EndAt,
	}
}

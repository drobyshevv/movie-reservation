package sessions

import (
	"context"
	"log/slog"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	"github.com/drobyshevv/movie-reservation/internal/models"
	"github.com/go-playground/validator/v10"
)

type SessionHandler struct {
	log       *slog.Logger
	service   SessionService
	validator *validator.Validate
}

func NewSessionHandler(log *slog.Logger, service SessionService, validator *validator.Validate) *SessionHandler {
	return &SessionHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

type SessionService interface {
	GetSessions(ctx context.Context) ([]models.SessionWithDetails, error)
	GetSession(ctx context.Context, id int64) (*models.SessionWithDetails, error)
	CreateSession(ctx context.Context, params models.CreateSessionInput) (*domain.Session, error)
	UpdateSession(ctx context.Context, id int64, params models.UpdateSessionInput) (*domain.Session, error)
	DeleteSession(ctx context.Context, id int64) error
}

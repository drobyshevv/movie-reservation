package handler

import (
	"context"
	"log/slog"

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
	GetSessions(ctx context.Context) ([]models.Session, error)
	GetSession(ctx context.Context, id int64) (models.Session, error)
	CreateSession(ctx context.Context, hallID int64, movieID int64, params models.CreateSessionParams) ([]models.Session, error)
	UpdateSession(ctx context.Context, id int64, params models.UpdateSessionParams) (models.Session, error)
	DeleteSession(ctx context.Context, id int64) error
}

//TODO Если меняется StartAt, а EndAt не передан — пересчитывать EndAt автоматически или оставить старый?

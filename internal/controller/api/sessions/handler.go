package sessions

import (
	"context"
	"log/slog"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	sessionuc "github.com/drobyshevv/movie-reservation/internal/usecase/sessions"
	"github.com/go-playground/validator/v10"
)

type SessionHandler struct {
	log       *slog.Logger
	usecase   SessionUseCase
	validator *validator.Validate
}

func NewSessionHandler(log *slog.Logger, usecase SessionUseCase, validator *validator.Validate) *SessionHandler {
	return &SessionHandler{
		log:       log,
		usecase:   usecase,
		validator: validator,
	}
}

type SessionUseCase interface {
	List(ctx context.Context) ([]sessionuc.SessionWithDetails, error)
	Get(ctx context.Context, id int64) (*sessionuc.SessionWithDetails, error)
	Create(ctx context.Context, input sessionuc.CreateSessionInput) (*domain.Session, error)
	Update(ctx context.Context, id int64, input sessionuc.UpdateSessionInput) (*domain.Session, error)
	Delete(ctx context.Context, id int64) error
}

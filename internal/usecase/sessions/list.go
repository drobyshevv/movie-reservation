package sessions

import (
	"context"
	"fmt"
)

func (u *SessionUseCase) List(ctx context.Context) ([]SessionWithDetails, error) {
	const op = "usecase.sessions.List"

	sessions, err := u.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sessions, nil
}

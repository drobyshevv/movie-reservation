package sessions

import (
	"context"
	"fmt"
)

func (u *UseCase) Get(ctx context.Context, id int64) (*SessionWithDetails, error) {
	const op = "usecase.sessions.Get"

	session, err := u.repository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

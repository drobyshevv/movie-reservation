package halls

import (
	"context"
	"fmt"
)

func (u *HallUseCase) GetSummary(ctx context.Context, id int64) (*HallSummary, error) {
	const op = "usecase.halls.GetSummary"

	hall, err := u.repository.GetSummary(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return hall, nil
}

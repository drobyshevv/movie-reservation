package halls

import "github.com/drobyshevv/movie-reservation/internal/domain"

type HallWithSeats struct {
	Hall  domain.Hall
	Seats []domain.Seat
}

type HallSummary struct {
	ID         int64
	Name       string
	TotalSeats int
}

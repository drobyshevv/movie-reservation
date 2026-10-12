package models

import "time"

type CreateSessionParams struct {
	HallID  int64
	MovieID int64
	StartAt time.Time
	EndAt   *time.Time
}

type UpdateSessionParams struct {
	HallID  *int64
	MovieID *int64
	StartAt *time.Time
	EndAt   *time.Time
}

type SessionWithDetails struct {
	ID        int64
	Status    StatusType
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

package models

import "time"

type StatusType string

const (
	StatusScheduled StatusType = "scheduled"
	StatusActive    StatusType = "active"
	StatusCancelled StatusType = "cancelled"
)

type Session struct {
	ID        int64
	HallID    int64
	MovieID   int64
	Status    StatusType
	StartAt   time.Time
	EndAt     time.Time
	CreatedAt time.Time
}

type CreateSessionParams struct {
	HallID  int64
	MovieID int64
	StartAt time.Time
	EndAt   *time.Time
}

type UpdateSessionParams struct {
	HallID  *int64
	MovieID *int64
	Status  *StatusType
	StartAt *time.Time
	EndAt   *time.Time
}

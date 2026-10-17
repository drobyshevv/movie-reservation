package domain

import (
	"errors"
	"time"
)

type StatusType string

const (
	StatusScheduled StatusType = "scheduled"
	StatusActive    StatusType = "active"
	StatusCancelled StatusType = "cancelled"
)

type Session struct {
	id        int64
	hallID    int64
	movieID   int64
	status    StatusType
	period    TimeRange
	createdAt time.Time
}

func NewSession(
	hallID, movieID int64,
	period TimeRange,
) (Session, error) {
	if hallID <= 0 || movieID <= 0 {
		return Session{}, errors.New("invalid hall or movie ID")
	}

	return Session{
		hallID:  hallID,
		movieID: movieID,
		period:  period,
		status:  StatusScheduled,
	}, nil
}

func RestoreSession(
	id, hallID, movieID int64,
	status StatusType,
	period TimeRange,
	createdAt time.Time,
) (Session, error) {
	if id <= 0 || hallID <= 0 || movieID <= 0 {
		return Session{}, errors.New("invalid session IDs")
	}

	return Session{
		id:        id,
		hallID:    hallID,
		movieID:   movieID,
		status:    status,
		period:    period,
		createdAt: createdAt,
	}, nil
}

func (s *Session) ChangePeriod(period TimeRange) error {
	if s.status != StatusScheduled {
		return errors.New("only scheduled sessions can be changed")
	}

	s.period = period
	return nil
}

func (s *Session) ChangeStatus(status StatusType) error {
	switch status {
	case StatusScheduled:
		if s.status != StatusCancelled {
			return errors.New("only cancelled sessions can be scheduled")
		}
	case StatusActive:
		if s.status != StatusScheduled {
			return errors.New("only scheduled sessions can be activated")
		}
	case StatusCancelled:
		if s.status == StatusCancelled {
			return errors.New("session is already cancelled")
		}
	default:
		return errors.New("invalid session status")
	}

	s.status = status
	return nil
}

func (s *Session) ID() int64            { return s.id }
func (s *Session) HallID() int64        { return s.hallID }
func (s *Session) MovieID() int64       { return s.movieID }
func (s *Session) Status() StatusType   { return s.status }
func (s *Session) Period() TimeRange    { return s.period }
func (s *Session) CreatedAt() time.Time { return s.createdAt }

package domain

import "errors"

var (
	ErrEmptySeatRow      = errors.New("seat row cannot be empty")
	ErrInvalidSeatNumber = errors.New("seat number must be positive")
	ErrInvalidSeatID     = errors.New("seat ID must be positive")
	ErrInvalidSeatHallID = errors.New("seat hall ID must be positive")
)

type Seat struct {
	id     int64
	hallID int64
	row    string
	number int16
}

func NewSeat(hallID int64, row string, number int16) (*Seat, error) {
	if hallID <= 0 {
		return nil, ErrInvalidSeatHallID
	}
	if row == "" {
		return nil, ErrEmptySeatRow
	}
	if number <= 0 {
		return nil, ErrInvalidSeatNumber
	}

	return &Seat{
		hallID: hallID,
		row:    row,
		number: number,
	}, nil
}

func RestoreSeat(id int64, hallID int64, row string, number int16) (*Seat, error) {
	if id <= 0 {
		return nil, ErrInvalidSeatID
	}

	seat, err := NewSeat(hallID, row, number)
	if err != nil {
		return nil, err
	}

	seat.id = id
	return seat, nil
}

func (s *Seat) ID() int64     { return s.id }
func (s *Seat) HallID() int64 { return s.hallID }
func (s *Seat) Row() string   { return s.row }
func (s *Seat) Number() int16 { return s.number }

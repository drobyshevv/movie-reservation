package models

type Seat struct {
	ID     int64
	HallID int64
	Row    string
	Number int16
}

type CreateHallSeatsParams struct {
	Row    string
	Number int16
}

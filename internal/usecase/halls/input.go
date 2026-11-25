package halls

type CreateHallWithSeatsInput struct {
	Name  string
	Seats []HallSeatsInput
}

type HallSeatsInput struct {
	Row    string
	Number int16
}

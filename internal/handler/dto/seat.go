package dto

type HallSeatResponse struct {
	ID     int64  `json:"id"`
	Row    string `json:"row"`
	Number int16  `json:"number"`
}

type CreateHallSeatsRequest struct {
	Row    string `json:"row" validate:"required"`
	Number int16  `json:"number" validate:"required,gt=0"`
}

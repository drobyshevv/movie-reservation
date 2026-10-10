package dto

import "github.com/drobyshevv/movie-reservation/internal/models"

type HallSeatResponse struct {
	ID     int64  `json:"id"`
	Row    string `json:"row"`
	Number int16  `json:"number"`
}

type CreateHallSeatsRequest struct {
	Row    string `json:"row" validate:"required"`
	Number int16  `json:"number" validate:"required,gt=0"`
}

func MapCreateHallSeatsRequest(req []CreateHallSeatsRequest) []models.CreateHallSeatsParams {
	var seats []models.CreateHallSeatsParams

	for _, r := range req {
		for n := int16(1); n <= r.Number; n++ {
			seats = append(seats, models.CreateHallSeatsParams{
				Row:    r.Row,
				Number: r.Number,
			})
		}
	}

	return seats
}

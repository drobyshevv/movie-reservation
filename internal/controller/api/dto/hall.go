package dto

import (
	"strconv"

	"github.com/drobyshevv/movie-reservation/internal/domain"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
	halluc "github.com/drobyshevv/movie-reservation/internal/usecase/halls"
)

type CreateHallRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type CreateHallWithSeatsRequest struct {
	Name  string                   `json:"name" validate:"required,max=100"`
	Seats []CreateHallSeatsRequest `json:"seats" validate:"required,dive"`
}

type CreateHallSeatsRequest struct {
	Row    string `json:"row" validate:"required,max=10"`
	Number int16  `json:"number" validate:"required,gt=0"`
}

type UpdateHallRequest struct {
	Name string `json:"name" validate:"omitempty,max=100"`
}

func (req *CreateHallWithSeatsRequest) ToInput() (halluc.CreateHallWithSeatsInput, error) {
	input := halluc.CreateHallWithSeatsInput{
		Name:  req.Name,
		Seats: make([]halluc.HallSeatsInput, 0, len(req.Seats)),
	}

	seen := make(map[string]struct{}, len(req.Seats))

	for _, s := range req.Seats {
		key := s.Row + "_" + strconv.Itoa(int(s.Number))
		if _, exists := seen[key]; exists {
			return halluc.CreateHallWithSeatsInput{}, errs.ErrDuplicateSeatsInRequest
		}
		seen[key] = struct{}{}

		input.Seats = append(input.Seats, halluc.HallSeatsInput{
			Row:    s.Row,
			Number: s.Number,
		})
	}

	return input, nil
}

func ToSeatsInput(req []CreateHallSeatsRequest) ([]halluc.HallSeatsInput, error) {
	input := make([]halluc.HallSeatsInput, 0, len(req))
	seen := make(map[string]struct{}, len(req))

	for _, s := range req {
		key := s.Row + "_" + strconv.Itoa(int(s.Number))
		if _, exists := seen[key]; exists {
			return nil, errs.ErrDuplicateSeatsInRequest
		}
		seen[key] = struct{}{}

		input = append(input, halluc.HallSeatsInput{
			Row:    s.Row,
			Number: s.Number,
		})
	}

	return input, nil
}

type HallResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type HallWithSeatsResponse struct {
	ID    int64              `json:"id"`
	Name  string             `json:"name"`
	Seats []HallSeatResponse `json:"seats"`
}

type HallSeatResponse struct {
	ID     int64  `json:"id"`
	Row    string `json:"row"`
	Number int16  `json:"number"`
}

type HallSummaryResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	TotalSeats int    `json:"total_seats"`
}

func (h *HallResponse) FromModel(m domain.Hall) {
	h.ID = m.ID()
	h.Name = m.Name()
}

func ToHallListResponse(halls []domain.Hall) []HallResponse {
	resp := make([]HallResponse, len(halls))
	for i, hall := range halls {
		resp[i] = HallResponse{
			ID:   hall.ID(),
			Name: hall.Name(),
		}
	}
	return resp
}

func (h *HallWithSeatsResponse) FromModel(m halluc.HallWithSeats) {
	h.ID = m.Hall.ID()
	h.Name = m.Hall.Name()

	h.Seats = make([]HallSeatResponse, len(m.Seats))
	for i, seat := range m.Seats {
		h.Seats[i] = HallSeatResponse{
			ID:     seat.ID(),
			Row:    seat.Row(),
			Number: seat.Number(),
		}
	}
}

func (h *HallSummaryResponse) FromModel(m halluc.HallSummary) {
	h.ID = m.ID
	h.Name = m.Name
	h.TotalSeats = m.TotalSeats
}

func ToSeatListResponse(seats []domain.Seat) []HallSeatResponse {
	resp := make([]HallSeatResponse, len(seats))
	for i, seat := range seats {
		resp[i] = HallSeatResponse{
			ID:     seat.ID(),
			Row:    seat.Row(),
			Number: seat.Number(),
		}
	}
	return resp
}

package dto

import "github.com/drobyshevv/movie-reservation/internal/models"

type CreateHallRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Capacity int    `json:"capacity" validate:"required,gt=0,lte=1000"`
}

type UpdateHallRequest struct {
	Name     *string `json:"name" validate:"omitempty,max=100"`
	Capacity *int    `json:"capacity" validate:"omitempty,gt=0,lte=1000"`
}

type HallResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

func (h *HallResponse) FromModel(m *models.Hall) {
	h.ID = m.ID
	h.Name = m.Name
	h.Capacity = m.Capacity
}

func (h *CreateHallRequest) ToModel() *models.CreateHallParams {
	return &models.CreateHallParams{
		Name:     h.Name,
		Capacity: h.Capacity,
	}
}

func (h *UpdateHallRequest) ToModel() *models.UpdateHallParams {
	return &models.UpdateHallParams{
		Name:     h.Name,
		Capacity: h.Capacity,
	}
}

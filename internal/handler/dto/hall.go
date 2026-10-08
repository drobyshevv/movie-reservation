package dto

import "github.com/drobyshevv/movie-reservation/internal/models"

type CreateHallRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type UpdateHallRequest struct {
	Name *string `json:"name" validate:"omitempty,max=100"`
}

type HallResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (h *HallResponse) FromModel(m *models.Hall) {
	h.ID = m.ID
	h.Name = m.Name
}

func (h *CreateHallRequest) ToModel() *models.CreateHallParams {
	return &models.CreateHallParams{
		Name: h.Name,
	}
}

func (h *UpdateHallRequest) ToModel() *models.UpdateHallParams {
	return &models.UpdateHallParams{
		Name: h.Name,
	}
}

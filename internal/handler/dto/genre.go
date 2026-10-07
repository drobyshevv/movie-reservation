package dto

import "github.com/drobyshevv/movie-reservation/internal/models"

type CreateGenreRequest struct {
	TypeGenre string `json:"type_genre" validate:"required,max=100"`
}

type UpdateGenreRequest struct {
	TypeGenre string `json:"type_genre" validate:"required,max=100"`
}

type GenreResponse struct {
	ID        int64  `json:"id"`
	TypeGenre string `json:"type_genre"`
}

func (g *GenreResponse) FromModel(m *models.Genre) {
	g.ID = m.ID
	g.TypeGenre = m.TypeGenre
}

func (g *CreateGenreRequest) ToModel() *models.CreateGenreParams {
	return &models.CreateGenreParams{
		TypeGenre: g.TypeGenre,
	}
}

func (g *UpdateGenreRequest) ToModel() *models.UpdateGenreParams {
	return &models.UpdateGenreParams{
		TypeGenre: g.TypeGenre,
	}
}

package dto

import (
	"github.com/drobyshevv/movie-reservation/internal/domain"
	genreuc "github.com/drobyshevv/movie-reservation/internal/usecase/genres"
)

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

type GenreMovieResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

func (g *GenreResponse) FromModel(m *domain.Genre) {
	g.ID = m.ID()
	g.TypeGenre = m.Name()
}

func (g *CreateGenreRequest) ToModel() string {
	return g.TypeGenre
}

func (g *UpdateGenreRequest) ToModel() string {
	return g.TypeGenre
}

func (g *GenreMovieResponse) FromModel(gm *genreuc.GenreMovie) {
	g.ID = gm.ID
	g.Title = gm.Title
}

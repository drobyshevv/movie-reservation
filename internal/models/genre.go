package models

type Genre struct {
	ID        int64
	TypeGenre string
}

type CreateGenreParams struct {
	TypeGenre string
}

type UpdateGenreParams struct {
	TypeGenre string
}

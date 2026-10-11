package models

type Hall struct {
	ID   int64
	Name string
}

type CreateHallParams struct {
	Name string
}

type UpdateHallParams struct {
	Name *string
}

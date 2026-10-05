package models

type Hall struct {
	ID       int64
	Name     string
	Capacity int
}

type CreateHallParams struct {
	Name     string
	Capacity int
}

type UpdateHallParams struct {
	Name     *string
	Capacity *int
}

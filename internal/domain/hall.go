package domain

import "errors"

var (
	ErrHallNameIsEmpty = errors.New("hall name is empty")
)

type Hall struct {
	id   int64
	name string
}

func NewHall(name string) (*Hall, error) {
	if len(name) == 0 {
		return nil, ErrHallNameIsEmpty
	}
	return &Hall{
		name: name,
	}, nil
}

func RestoreHall(id int64, name string) (Hall, error) {
	if name == "" {
		return Hall{}, ErrHallNameIsEmpty
	}

	return Hall{
		id:   id,
		name: name,
	}, nil
}

func (h *Hall) ChangeName(name string) error {
	if len(name) == 0 {
		return ErrHallNameIsEmpty
	}
	return nil
}

func (h *Hall) ID() int64    { return h.id }
func (h *Hall) Name() string { return h.name }

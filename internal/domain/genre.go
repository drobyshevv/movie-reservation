package domain

import (
	"errors"
	"strings"
)

var (
	ErrEmptyTypeGenre   = errors.New("genre type cannot be empty")
	ErrTypeGenreTooLong = errors.New("genre type is too long")
)

type Genre struct {
	id   int64
	name string
}

func NewGenre(typeGenre string) (Genre, error) {
	typeGenre = strings.TrimSpace(typeGenre)

	if typeGenre == "" {
		return Genre{}, ErrEmptyTypeGenre
	}

	if len([]rune(typeGenre)) > 50 {
		return Genre{}, ErrTypeGenreTooLong
	}

	return Genre{name: typeGenre}, nil
}

func RestoreGenre(id int64, typeGenre string) (Genre, error) {
	if id <= 0 {
		return Genre{}, errors.New("invalid genre ID")
	}

	genre, err := NewGenre(typeGenre)
	if err != nil {
		return Genre{}, err
	}

	genre.id = id
	return genre, nil
}

func (g *Genre) ChangeType(typeGenre string) error {
	typeGenre = strings.TrimSpace(typeGenre)

	if typeGenre == "" {
		return ErrEmptyTypeGenre
	}

	if len([]rune(typeGenre)) > 50 {
		return ErrTypeGenreTooLong
	}

	g.name = typeGenre
	return nil
}

func (g Genre) ID() int64    { return g.id }
func (g Genre) Name() string { return g.name }

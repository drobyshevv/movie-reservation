package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidMovieDuration = errors.New("movie duration must be positive")
	ErrInvalidMovieID       = errors.New("movie ID must be positive")
	ErrEmptyMovieTitle      = errors.New("movie title cannot be empty")
)

type Movie struct {
	id          int64
	title       string
	description *string
	duration    time.Duration
	createdAt   time.Time
}

func NewMovie(title string, description *string, duration time.Duration) (Movie, error) {
	if title == "" {
		return Movie{}, ErrEmptyMovieTitle
	}

	if duration <= 0 {
		return Movie{}, ErrInvalidMovieDuration
	}

	return Movie{
		title:       title,
		description: description,
		duration:    duration,
	}, nil
}

func RestoreMovie(
	id int64,
	title string,
	description *string,
	duration time.Duration,
	createdAt time.Time,
) (Movie, error) {
	if id <= 0 {
		return Movie{}, errors.New("movie ID must be positive")
	}

	movie, err := NewMovie(title, description, duration)
	if err != nil {
		return Movie{}, err
	}

	movie.id = id
	movie.createdAt = createdAt

	return movie, nil
}

func (m *Movie) ChangeTitle(title string) error {
	if title == "" {
		return ErrEmptyMovieTitle
	}
	m.title = title
	return nil
}

func (m *Movie) ChangeDescription(description *string) {
	m.description = description
}

func (m *Movie) ChangeDuration(duration time.Duration) error {
	if duration <= 0 {
		return ErrInvalidMovieDuration
	}

	m.duration = duration
	return nil
}

func (m Movie) ID() int64               { return m.id }
func (m Movie) Title() string           { return m.title }
func (m Movie) Description() *string    { return m.description }
func (m Movie) Duration() time.Duration { return m.duration }
func (m Movie) CreatedAt() time.Time    { return m.createdAt }

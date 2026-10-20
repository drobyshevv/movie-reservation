package movie

import "github.com/drobyshevv/movie-reservation/internal/domain"

type MovieWithGenres struct {
	Movie  domain.Movie
	Genres []domain.Genre
}

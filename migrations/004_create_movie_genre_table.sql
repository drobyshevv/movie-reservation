-- +goose Up
CREATE TABLE movie_genre (
    movie_id integer references movies (id) on delete cascade,
    genre_id integer references genres (id) on delete cascade,
    primary key(movie_id, genre_id)
);

-- +goose Down
DROP TABLE IF EXISTS movie_genre;
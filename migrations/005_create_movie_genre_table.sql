-- +goose Up
CREATE TABLE movie_genre (
    movie_id BIGINT REFERENCES movies (id) ON DELETE CASCADE,
    genre_id BIGINT REFERENCES genres (id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, genre_id)
);

-- +goose Down
DROP TABLE IF EXISTS movie_genre;
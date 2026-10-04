-- +goose Up

-- movie_genre table indexes
CREATE INDEX idx_movie_genre_genre_id ON movie_genre (genre_id);

-- sessions table indexes
CREATE INDEX idx_sessions_hall_id ON sessions (hall_id);
CREATE INDEX idx_sessions_movie_id ON sessions (movie_id);

-- bookings table indexes
CREATE INDEX idx_bookings_session_id ON bookings (session_id);
CREATE INDEX idx_bookings_user_id ON bookings (user_id);
CREATE INDEX idx_bookings_user_created ON bookings (user_id, created_at DESC);

CREATE UNIQUE INDEX idx_bookings_active_seat
ON bookings (session_id, seat_number)
WHERE status IN ('pending', 'confirmed');

-- +goose Down

DROP INDEX IF EXISTS idx_movie_genre_genre_id;
DROP INDEX IF EXISTS idx_sessions_hall_id;
DROP INDEX IF EXISTS idx_sessions_movie_id;
DROP INDEX IF EXISTS idx_bookings_session_id;
DROP INDEX IF EXISTS idx_bookings_user_id;
DROP INDEX IF EXISTS idx_bookings_user_created;
DROP INDEX IF EXISTS idx_bookings_active_seat;
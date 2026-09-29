-- +goose Up

-- movie_genre table indexes
create index idx_movie_genre_genre_id on movie_genre (genre_id);

-- sessions table indexes
create index idx_sessions_hall_id on sessions (hall_id);
create index idx_sessions_movie_id on sessions (movie_id);

-- bookings table indexes
create index idx_bookings_session_id on bookings (session_id);
create index idx_bookings_user_id on bookings (user_id);
create index idx_bookings_user_created on bookings (user_id, created_at DESC);

create unique index idx_bookings_active_seat 
on bookings (session_id, seat_number)
where status in ('pending', 'confirmed');

-- +goose Down
drop index if EXISTS idx_movie_genre_genre_id;
drop index if EXISTS idx_sessions_hall_id;
drop index if EXISTS idx_sessions_movie_id;
drop index if EXISTS idx_bookings_session_id;
drop index if EXISTS idx_bookings_user_id;
drop index if EXISTS idx_bookings_user_created;
drop index if EXISTS idx_bookings_status;
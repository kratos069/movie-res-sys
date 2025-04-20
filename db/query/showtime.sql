-- name: CreateShowtime :one
INSERT INTO showtimes (movie_id, start_time, price)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListShowtimesByDate :many
SELECT s.*, m.title, m.poster_url
FROM showtimes s
JOIN movies m ON m.movie_id = s.movie_id
WHERE DATE(s.start_time) = $1
ORDER BY s.start_time;

-- name: GetShowtime :one
SELECT * FROM showtimes
WHERE showtime_id = $1;

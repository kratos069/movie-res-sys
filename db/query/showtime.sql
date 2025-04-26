-- name: CreateShowtime :one
INSERT INTO showtimes (movie_id, start_time, price)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListShowtimesByDate :many
SELECT s.showtime_id, s.movie_id, s.start_time, s.price, s.created_at, m.title, m.poster_url
FROM showtimes s
JOIN movies m ON m.movie_id = s.movie_id
WHERE s.start_time >= $1
ORDER BY s.start_time;

-- name: ListShowtimesBetween :many
SELECT s.showtime_id, s.movie_id, s.start_time, s.price, s.created_at, m.title, m.poster_url
FROM showtimes s
JOIN movies m ON m.movie_id = s.movie_id
WHERE s.start_time >= $1 AND s.start_time < $2
ORDER BY s.start_time;

-- name: GetShowtime :one
SELECT * FROM showtimes
WHERE showtime_id = $1;

-- name: DeleteShowtime :exec
DELETE FROM showtimes
WHERE showtime_id = $1;
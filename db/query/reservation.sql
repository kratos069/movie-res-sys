-- name: ReserveSeat :one
INSERT INTO reservations (user_id, showtime_id, seat_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CancelReservation :exec
DELETE FROM reservations
WHERE reservation_id = $1 AND user_id = $2;

-- name: ListReservationsByUser :many
SELECT r.*, s.start_time, m.title, se.row, se.number
FROM reservations r
JOIN showtimes s ON s.showtime_id = r.showtime_id
JOIN movies m ON m.movie_id = s.movie_id
JOIN seats se ON se.seat_id = r.seat_id
WHERE r.user_id = $1
ORDER BY s.start_time;

-- name: ListAvailableSeatsForShowtime :many
SELECT *
FROM seats
WHERE seat_id NOT IN (
  SELECT seat_id FROM reservations WHERE showtime_id = $1
)
ORDER BY row, number;

-- name: ListReservationsByShowtime :many
SELECT r.*, u.name, se.row, se.number
FROM reservations r
JOIN users u ON u.user_id = r.user_id
JOIN seats se ON se.seat_id = r.seat_id
WHERE r.showtime_id = $1;

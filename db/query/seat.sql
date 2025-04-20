-- name: ListAllSeats :many
SELECT * FROM seats
ORDER BY row, number;

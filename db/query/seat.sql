-- name: ListAllSeats :many
SELECT * FROM seats
ORDER BY row, number;

-- name: ListSeatsForShowtime :many
SELECT 
    s.seat_id,
    s.row,
    s.number,
    CASE WHEN r.seat_id IS NOT NULL THEN true ELSE false END AS is_booked
FROM seats s
LEFT JOIN reservations r 
    ON s.seat_id = r.seat_id AND r.showtime_id = $1
ORDER BY s.row, s.number;
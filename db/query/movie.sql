-- name: CreateMovie :one
INSERT INTO movies (title, description, poster_url, genre_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListMovies :many
SELECT * FROM movies
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: GetMovie :one
SELECT * FROM movies
WHERE movie_id = $1;

-- name: UpdateMovie :one
UPDATE movies
SET title = $2,
    description = $3,
    poster_url = $4,
    genre_id = $5
WHERE movie_id = $1
RETURNING *;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE movie_id = $1;

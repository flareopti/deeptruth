-- name: CreateAuthor :one
INSERT INTO authors (
    name,
    rating,
    description,
    created_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetAuthor :one
SELECT * FROM authors 
WHERE id = $1 LIMIT 1;

-- name: ListAuthor :many
SELECT * FROM authors
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAuthor :one
UPDATE authors
SET rating = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $1;
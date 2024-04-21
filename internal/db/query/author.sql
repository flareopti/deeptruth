-- name: CreateAuthor :one
INSERT INTO authors (
    name,
    rating,
    description
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAuthor :one
SELECT * FROM authors 
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAuthorRating :one
UPDATE authors
SET rating = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $1;
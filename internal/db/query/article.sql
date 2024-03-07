-- name: CreateArticle :one
INSERT INTO articles (
    author_id,
    title,
    content,
    verdict,
    rating
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetArticle :one
SELECT * FROM articles
WHERE id = $1 LIMIT 1;

-- name: ListArticle :many
SELECT * FROM articles
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateArticle :one
UPDATE articles
SET rating = $2, verdict = $3
WHERE id = $1
RETURNING *;

-- name: DeleteArticle :exec
DELETE FROM articles WHERE id = $1;
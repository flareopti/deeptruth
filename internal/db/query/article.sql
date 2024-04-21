-- name: CreateArticle :one
INSERT INTO articles (
    author_id,
    title,
    content,
    rating
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetArticle :one
SELECT * FROM articles
WHERE id = $1 LIMIT 1;

-- name: GetArticlesCount :one
SELECT COUNT(*) FROM articles;

-- name: ListArticles :many
SELECT * FROM articles
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateArticleRating :one
UPDATE articles
SET rating = $2
WHERE id = $1
RETURNING *;

-- name: DeleteArticle :exec
DELETE FROM articles WHERE id = $1;
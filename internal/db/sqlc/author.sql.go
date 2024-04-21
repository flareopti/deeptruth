// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: author.sql

package db

import (
	"context"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (
    name,
    rating,
    description
) VALUES (
    $1, $2, $3
) RETURNING id, name, description, rating, created_at
`

type CreateAuthorParams struct {
	Name        string `json:"name"`
	Rating      int32  `json:"rating"`
	Description string `json:"description"`
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	row := q.db.QueryRow(ctx, createAuthor, arg.Name, arg.Rating, arg.Description)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Rating,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $1
`

func (q *Queries) DeleteAuthor(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteAuthor, id)
	return err
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, name, description, rating, created_at FROM authors 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAuthor(ctx context.Context, id int64) (Author, error) {
	row := q.db.QueryRow(ctx, getAuthor, id)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Rating,
		&i.CreatedAt,
	)
	return i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT id, name, description, rating, created_at FROM authors
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListAuthorsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAuthors(ctx context.Context, arg ListAuthorsParams) ([]Author, error) {
	rows, err := q.db.Query(ctx, listAuthors, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Author{}
	for rows.Next() {
		var i Author
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Rating,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAuthorRating = `-- name: UpdateAuthorRating :one
UPDATE authors
SET rating = $2
WHERE id = $1
RETURNING id, name, description, rating, created_at
`

type UpdateAuthorRatingParams struct {
	ID     int64 `json:"id"`
	Rating int32 `json:"rating"`
}

func (q *Queries) UpdateAuthorRating(ctx context.Context, arg UpdateAuthorRatingParams) (Author, error) {
	row := q.db.QueryRow(ctx, updateAuthorRating, arg.ID, arg.Rating)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Rating,
		&i.CreatedAt,
	)
	return i, err
}

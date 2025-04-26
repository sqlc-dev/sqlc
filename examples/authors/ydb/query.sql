-- name: ListAuthors :many 
SELECT * FROM authors;

-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $p0;

-- name: ListAuthorsWithIdModulo :many
SELECT * FROM authors
WHERE id % 2 = 0;

-- name: GetAuthorsByName :many
SELECT * FROM authors
WHERE name = $p0;

-- name: ListAuthorsWithNullBio :many
SELECT * FROM authors
WHERE bio IS NULL;

-- name: CreateOrUpdateAuthor :execresult 
UPSERT INTO authors (id, name, bio) VALUES ($p0, $p1, $p2);

-- name: CreateOrUpdateAuthorRetunringBio :one
UPSERT INTO authors (id, name, bio) VALUES ($p0, $p1, $p2) RETURNING bio;

-- name: DeleteAuthor :exec 
DELETE FROM authors
WHERE id = $p0;


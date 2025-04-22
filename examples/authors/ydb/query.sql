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


-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = @p1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :exec
INSERT INTO authors (name, bio) VALUES (@p1, @p2);

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = @p1;

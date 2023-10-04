-- name: GetAuthor :one
SELECT * FROM Authors
WHERE ID = ? LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM Authors
ORDER BY Name;

-- name: CreateAuthor :execresult
INSERT INTO Authors (
  Name, Bio
) VALUES (
  ?, ?
);

-- name: DeleteAuthor :exec
DELETE FROM Authors
WHERE ID = ?;

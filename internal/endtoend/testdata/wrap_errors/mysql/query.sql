-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = ? LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :execlastid
INSERT INTO authors (
          name, bio
) VALUES (
  ?, ?
);

-- name: DeleteAuthorExec :exec
DELETE FROM authors
WHERE id = ?;

-- name: DeleteAuthorExecRows :execrows
DELETE FROM authors
WHERE id = ?;

-- name: DeleteAuthorExecLastID :execlastid
DELETE FROM authors
WHERE id = ?;

-- name: DeleteAuthorExecResult :execresult
DELETE FROM authors
WHERE id = ?;

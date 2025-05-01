-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :execlastid
INSERT INTO authors (
          name, bio
) VALUES (
  $1, $2
);

-- name: DeleteAuthorExec :exec
DELETE FROM authors
WHERE id = $1;

-- name: DeleteAuthorExecRows :execrows
DELETE FROM authors
WHERE id = $1;

-- name: DeleteAuthorExecLastID :execlastid
DELETE FROM authors
WHERE id = $1;

-- name: DeleteAuthorExecResult :execresult
DELETE FROM authors
WHERE id = $1;

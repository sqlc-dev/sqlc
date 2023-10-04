-- name: DeleteAuthor :exec
UPDATE
  authors,
  books
SET
  authors.deleted_at = now(),
  books.deleted_at = now()
WHERE
  books.is_amazing = 1
  AND authors.name = sqlc.arg(name);
-- name: ListBooksByAuthor :many
SELECT b.*
FROM books b
WHERE b.author_id = sqlc.arg(author_id)
  AND sqlc.macro.NotDeleted(b.deleted_at)
ORDER BY b.title;

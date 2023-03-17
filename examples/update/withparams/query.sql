/* name: CreateAuthor :execresult */
INSERT INTO
  authors (name, deleted_at, updated_at)
VALUES
  (?, ?, ?);

/* name: CreateBook :execresult */
INSERT INTO
  books (is_amazing)
VALUES
  (?);

/* name: GetAuthor :one */
SELECT
  *
FROM
  authors
WHERE
  id = ?
LIMIT
  1;

/* name: DeleteAuthor :exec */
UPDATE
  authors,
  books
SET
  authors.deleted_at = now(),
  authors.updated_at = now()
WHERE
  books.is_amazing = 1
  AND authors.name = sqlc.arg(name);
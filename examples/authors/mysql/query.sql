/* name: GetAuthor :one */
SELECT * FROM authors
WHERE id = ? LIMIT 1;

/* name: ListAuthors :many */
SELECT * FROM authors
ORDER BY name;

/* name: CreateAuthor :execresult */
INSERT INTO authors (
  name, bio
) VALUES (
  ?, ?
);

/* name: DeleteAuthor :exec */
DELETE FROM authors
WHERE id = ?;

/* name: UpdateAuthor :exec */
UPDATE authors
SET
 name = coalesce(sqlc.narg('name'), name),
 bio = coalesce(sqlc.narg('bio'), bio)
WHERE id = sqlc.arg('id');

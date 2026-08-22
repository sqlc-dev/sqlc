-- name: GetAuthor :one
select  id,name , bio
from   authors
where id =  $1 limit 1;

-- a helpful comment
-- name: ListAuthors :many
SELECT id, name, bio FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING *;

-- name: SearchAuthors :many
SELECT id, name, bio, created_at FROM authors WHERE name LIKE $1 AND bio IS NOT NULL AND id > $2 AND name <> $3 ORDER BY name, id LIMIT $4;

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = @id

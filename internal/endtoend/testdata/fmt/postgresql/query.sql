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

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = @id

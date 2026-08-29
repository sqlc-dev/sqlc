-- name: GetAuthor :one
select  id,name , bio
from   authors
where id =  $1 limit 1; -- the primary lookup

/* This listing powers the admin page.
   Keep it ordered by name so the UI stays stable. */
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

-- name: PickyQuery :many
SELECT id, -- the primary key
       name
FROM authors
WHERE id > $1;

-- name: SearchAuthors :many
SELECT id, name, bio, created_at FROM authors WHERE name LIKE $1 AND bio IS NOT NULL AND id > $2 AND name <> $3 ORDER BY name, id LIMIT $4;

-- name: AtParamsGlued :many
SELECT name FROM authors WHERE name = @slug AND @filter::bool;

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = @id
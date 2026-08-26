-- name: GetAuthor :one
SELECT id, name, bio FROM authors WHERE id = $1;

-- name: FilterAuthors :many
SELECT id, name FROM authors WHERE name = $name AND royalties > $royalties;

-- name: SearchAuthors :many
SELECT id FROM authors WHERE name LIKE ? AND bio IS NOT NULL;

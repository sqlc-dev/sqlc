-- name: GetAuthor :one
SELECT id, name, bio FROM authors WHERE id = @id;

-- name: FilterAuthors :many
SELECT id, name FROM authors WHERE name = @name AND royalties > @royalties;

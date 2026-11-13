-- name: AuthorBooks :many
SELECT sqlc.embed(a), sqlc.embed(b) FROM authors a LEFT JOIN books b ON b.author_id = a.id;

-- name: AuthorHasBio :many
SELECT a.name, (a.bio is not null) as has_bio FROM authors a;

-- name: ListAuthors :many
SELECT id, name, bio, royalties, created FROM authors;

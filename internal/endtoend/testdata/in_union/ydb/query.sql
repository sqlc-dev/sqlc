-- name: GetAuthors :many
SELECT * FROM authors
WHERE id IN (SELECT author_id FROM book1 UNION SELECT author_id FROM book2);



-- name: GetAuthors :many
SELECT * FROM authors
WHERE author_id IN (SELECT author_id FROM book1 UNION SELECT author_id FROM book2);

-- name: ListAuthorsUnion :many
SELECT name AS foo FROM authors
UNION
SELECT first_name AS foo FROM people
ORDER BY foo;


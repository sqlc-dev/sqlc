-- name: ListAuthorsUnion :many
SELECT name as foo FROM authors
UNION
SELECT first_name as foo FROM people
ORDER BY foo;

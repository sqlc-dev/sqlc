-- name: ListAuthors :iter
-- Lists every author by name.
SELECT * FROM authors
ORDER BY name;

-- name: ListAuthorNamesAfter :iter
SELECT name FROM authors
WHERE id > $1
ORDER BY id;

-- name: DeleteAuthorsByName :iter
DELETE FROM authors
WHERE name = $1
RETURNING id, bio;

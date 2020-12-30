/* name: Get :one */
SELECT * FROM authors
WHERE author_id = ?;

/* name: Create :execresult */
INSERT INTO authors (name) VALUES (?);

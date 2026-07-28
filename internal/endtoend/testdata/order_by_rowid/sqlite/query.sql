-- name: ListAuthorsByRowid :many
SELECT id, name FROM authors ORDER BY rowid;

-- name: ListAuthorsByQualifiedRowid :many
SELECT id, name FROM authors ORDER BY authors._rowid_ DESC;

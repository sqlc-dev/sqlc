-- name: ByIDs :many
SELECT id, title FROM documents
WHERE id IN (SELECT CAST(value AS BIGINT) FROM STRING_SPLIT(@ids, ','));

-- name: Tokens :many
SELECT value FROM STRING_SPLIT(@list, ',');

-- name: TestSubqueryUnion :many
SELECT tmp.* FROM (
  SELECT * FROM authors
  UNION
  SELECT * FROM authors
) tmp LIMIT 5;

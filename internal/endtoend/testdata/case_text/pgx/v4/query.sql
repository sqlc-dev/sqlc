-- name: CaseStatementText :many
SELECT CASE 
  WHEN id = $1 THEN 'foo'
  ELSE 'bar'
END is_one
FROM foo;

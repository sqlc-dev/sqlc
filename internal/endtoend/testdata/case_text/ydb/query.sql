-- name: CaseStatementText :many
SELECT CASE 
  WHEN id = $id THEN 'foo'
  ELSE 'bar'
END AS is_one
FROM foo;





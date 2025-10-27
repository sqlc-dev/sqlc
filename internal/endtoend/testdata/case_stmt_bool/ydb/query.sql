-- name: CaseStatementBoolean :many
SELECT CASE 
  WHEN id = $id THEN true
  ELSE false
END AS is_one
FROM foo;





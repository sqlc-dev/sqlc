CREATE TABLE foo (id text not null);

-- name: CaseStatementBoolean :many
SELECT CASE 
  WHEN id = $1 THEN true
  ELSE false
END is_one
FROM foo;

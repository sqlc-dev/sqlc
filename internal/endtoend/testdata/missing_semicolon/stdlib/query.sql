-- name: FirstQuery :many
SELECT * FROM foo;

-- name: SecondQuery :many
SELECT * FROM foo WHERE email = $1

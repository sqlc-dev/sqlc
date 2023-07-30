-- name: TestIN :many
SELECT * FROM foo WHERE retyped IN (sqlc.slice(paramName));

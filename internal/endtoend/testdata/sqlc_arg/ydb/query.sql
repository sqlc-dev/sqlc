-- name: FuncParamIdent :many
SELECT name FROM foo WHERE name = $slug;

-- name: FuncParamString :many
SELECT name FROM foo WHERE name = $slug;


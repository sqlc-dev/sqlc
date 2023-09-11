CREATE TABLE foo (id int not null);

-- name: FuncParamIdent :many
SELECT id FROM foo WHERE id = ?;

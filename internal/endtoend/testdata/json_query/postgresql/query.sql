CREATE TABLE foo (bar jsonb not null);

-- name: SelectFieldText :one
SELECT * FROM foo WHERE bar ->> 'field_name' = $1;

-- name: SelectFieldJson :one
SELECT * FROM foo WHERE bar -> 'field_name' = $1;

-- name: SelectTypeCast :one
SELECT * FROM foo WHERE (bar ->> 'field_name')::text = $1;
